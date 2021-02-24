package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	port        int
	status      int
	contentType string
	delay       int
	body        string

	infoLog  = log.New(os.Stdout, "[INFO] ", log.LstdFlags)
	errorLog = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "mock-server is a mock http-server utility with configurable response parameters.\n"+
			"Usage of mock-server:\n")

		flag.PrintDefaults()
	}

	flag.IntVar(&port, "port", 8080, "Port of the server")
	flag.IntVar(&status, "status", 200, "Status code returned")
	flag.StringVar(&contentType, "content-type", "application/json", "Content type of the response")
	flag.IntVar(&delay, "delay", 0, "Amount of delay in response time")
	flag.StringVar(&body, "body", "", "Body returned")

	flag.Parse()

	if stdin, ok := readStdinWhenAvailable(); ok {
		body = stdin
	}

	http.HandleFunc("/", mockServerEndpoint)

	listenAddr := fmt.Sprintf(":%d", port)
	infoLog.Print("Server listening on ", listenAddr)

	errorLog.Fatal(http.ListenAndServe(listenAddr, nil))
}

func mockServerEndpoint(rw http.ResponseWriter, r *http.Request) {
	respStatus := getIntOrDefault(r, "status", status)
	respContentType := getStringOrDefault(r, "contentType", contentType)
	respDelay := getIntOrDefault(r, "delay", delay)
	respBody := getStringOrDefault(r, "body", body)

	infoLog.Printf(`%s
	Status: %d
	ContentType: %s
	Delay: %d
	Body:
%s`,
		r.URL.RequestURI(),
		respStatus,
		respContentType,
		respDelay,
		respBody,
	)

	// Artificial response time
	time.Sleep(time.Duration(delay) * time.Millisecond)

	rw.Header().Add("content-type", respContentType)
	rw.WriteHeader(respStatus)
	rw.Write([]byte(respBody))
}

func getStringOrDefault(r *http.Request, key, fallback string) string {
	override := r.URL.Query().Get(key)
	if override == "" {
		return fallback
	}
	return override
}

func getIntOrDefault(r *http.Request, key string, fallback int) int {
	override := r.URL.Query().Get(key)
	if override == "" {
		return fallback
	}
	d, err := strconv.Atoi(override)
	if err != nil {
		errorLog.Printf("Failed to convert URL parameter `%s` returning default. Error: %v", key, err)
		return fallback
	}
	return d
}

func readStdinWhenAvailable() (string, bool) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Stdin available
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		return string(b), true
	}
	// Stdin empty
	return "", false
}
