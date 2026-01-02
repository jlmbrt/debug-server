package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	dateTimeFormat = "02 Jan 2006 15:04:05 MST"
)

type Config struct {
	Address string
	Port    string
}

func (c *Config) Bind() string {
	return fmt.Sprintf("%s:%s", c.Address, c.Port)
}

var config Config

func main() {
	config = loadConfig()

	http.HandleFunc("/", handler)
	http.HandleFunc("/ready", readyHandler)

	fmt.Printf("Listening on %s\n", config.Bind())

	err := http.ListenAndServe(config.Bind(), nil)

	if err != nil {
		log.Fatalln("There's an error with the server", err)
	}
}

func loadConfig() Config {
	return Config{
		Address: getEnv("ADDR", "0.0.0.0"),
		Port:    getEnv("PORT", "8080"),
	}
}

func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); len(value) > 0 {
		return value
	} else {
		return defaultValue
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	sleep, _ := strconv.ParseInt(r.FormValue("sleep"), 10, 8)
	code, _ := strconv.ParseInt(r.FormValue("response-code"), 10, 16)
	body := r.FormValue("response-body")
	contentType := r.FormValue("response-content-type")

	if sleep > 0 {
		time.Sleep(time.Duration(sleep) * time.Second)
	}

	if code < 100 || code > 599 {
		code = 200
	}

	if len(contentType) == 0 {
		contentType = "application/json"
	}

	if len(body) == 0 {
		body = "{}"
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(int(code))
	w.Write([]byte(body))
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func logRequest(r *http.Request) {
	dateTime := time.Now()

	payload := make(map[string]any)

	payload["datetime"] = dateTime.Format(dateTimeFormat)
	payload["timestamp"] = dateTime.UnixMilli()
	payload["from"] = r.RemoteAddr
	payload["method"] = r.Method
	payload["path"] = r.URL.Path
	payload["query"] = r.URL.Query()
	payload["url"] = r.URL.RequestURI()
	payload["headers"] = r.Header

	r.ParseForm()

	contentType := r.Header.Get("Content-Type")

	switch contentType {
	case "application/json":
		rawBody, _ := io.ReadAll(r.Body)
		var body any
		json.Unmarshal(rawBody, &body)
		payload["body"] = body
	case "application/x-www-form-urlencoded", "multipart/form-data":
		payload["body"] = r.Form
	default:
		rawBody, _ := io.ReadAll(r.Body)
		payload["body"] = string(rawBody)
	}

	prettyBody, _ := json.MarshalIndent(payload, "", "  ")

	fmt.Println(string(prettyBody))
}
