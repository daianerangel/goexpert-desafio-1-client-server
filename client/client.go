package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func fetchQuote() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func main() {
	quote, err := fetchQuote()
	if err != nil {
		log.Fatalf("Failed to fetch quote: %v", err)
	}

	err = ioutil.WriteFile("cotacao.txt", []byte(quote), 0644)
	if err != nil {
		log.Fatalf("Failed to write quote to file: %v", err)
	}
}
