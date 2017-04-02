package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type weatherProvider interface {
	temperature(city string) (float64, error)
}

type weatherUndergroundProvider struct {
	apiKey string // 8d41ab297b1a83a0
}

func (w weatherUndergroundProvider) temperature(city string) (float64, error) {
	resp, err := http.Get("http://api.wunderground.com/api/" + w.apiKey + "/conditions/q/" + city + ".json")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var d struct {
		Observation struct {
			Temp float64 `json:"temp_c"`
		} `json:"current_observation"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return 0, err
	}

	kelvin := d.Observation.Temp + 273.15
	log.Printf("weatherUnderground: %s: %.2f", city, kelvin)
	return kelvin, nil
}

type openWeatherProvider struct{}

func (w openWeatherProvider) temperature(city string) (float64, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=9e8784c71f8117a9ffa2eafded53d5d9&q=" + city)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	var d struct {
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return 0, err
	}

	log.Printf("openWeatherMap: %s: %.2f", city, d.Main.Temp)
	return d.Main.Temp, nil
}

type multiWeatherProvider []weatherProvider

func (w multiWeatherProvider) temperature(city string) (float64, error) {
	temps := make(chan float64, len(w))
	errs := make(chan error, len(w))

	for _, provider := range w {
		go func(p weatherProvider) {
			k, err := p.temperature(city)
			if err != nil {
				errs <- err
				return
			}
			temps <- k
		}(provider)
	}

	sum := 0.0

	for i := 0; i < len(w); i++ {
		select {
		case temp := <-temps:
			sum += temp
		case err := <-errs:
			return 0, err
		}
	}

	return sum / float64(len(w)), nil
}

func main() {
	mw := multiWeatherProvider{
		openWeatherProvider{},
		weatherUnderground{apiKey: "8d41ab297b1a83a0"},
	}

	http.HandleFunc("/", hello)
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		data, err := mw.temperature(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"city": city,
			"temp": data,
			"took": time.Since(begin).String(),
		})
	})
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello GO!\n"))
}
