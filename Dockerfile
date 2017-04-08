# Version: 0.1.0 
FROM golang:1.8-alpine
MAINTAINER Dinesh Weerapurage "xydinesh@gmail.com"
ENV REFRESHED_AT 04-08-2017
RUN apk update
RUN apk add git
EXPOSE 8080
