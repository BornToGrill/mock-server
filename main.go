package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	port        int
	status      int
	delay       int
	body        string
	contentType string
)

func main() {

	flag.IntVar(&port, "port", 8080, "Port of the server")
	flag.IntVar(&status, "status", 200, "Status code returned")
	flag.IntVar(&delay, "delay", 0, "Amount of delay in response time")
	flag.StringVar(&body, "body", "", "Body returned")
	flag.StringVar(&contentType, "content-type", "application/json", "Content type of the response")

	flag.Parse()

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		respStatus := getStatus(r)
		respBody := getBody(r)
		respContentType := getContentType(r)
		delay := getDelay(r)

		// Artificial response time
		time.Sleep(time.Duration(delay) * time.Millisecond)

		rw.Header().Add("content-type", respContentType)
		rw.WriteHeader(respStatus)
		rw.Write([]byte(respBody))
	})

	listenAddr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server listening on %s\n", listenAddr)

	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func getStatus(r *http.Request) int {
	override := r.URL.Query().Get("status")
	if override == "" {
		return status
	}
	status, _ := strconv.Atoi(override)
	return status
}

func getDelay(r *http.Request) int {
	override := r.URL.Query().Get("delay")
	if override == "" {
		return delay
	}
	d, _ := strconv.Atoi(override)
	return d
}

func getContentType(r *http.Request) string {
	override := r.URL.Query().Get("contentType")
	if override == "" {
		return contentType
	}
	return override
}

func getBody(r *http.Request) string {
	override := r.URL.Query().Get("body")
	if override == "" {
		return body
	}
	return override
}
