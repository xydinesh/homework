package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func GetCommits(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://api.github.com/repos/stedolan/jq/commits?per_page=5")
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}
	defer resp.Body.Close()

	var mc multiCommits
	var nc newCommits

	if err := json.NewDecoder(resp.Body).Decode(&mc); err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}

	for _, commit := range mc {
		nc = append(nc, NewCommit{
			Committer: commit.Commit.Committer.Name,
			Message:   commit.Commit.Message,
			Sha:       commit.Sha,
		})
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(nc)
}
