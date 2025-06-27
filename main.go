package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

var (
	App *AppConfig
	Rdb *redis.Client
)

func main() {
	// Configuration
	token := "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"
	owner := "Ashrockzzz2003"
	repo := "webhook.soc"
	issueNumber := 1
	commentBody := "This comment is made by DevPool with my Personal Access Token"

	// API Endpoint
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d/comments",
		owner,
		repo,
		issueNumber,
	)

	payload := map[string]string{
		"body": commentBody,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating the request: %v\n", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Comment posted successfully!")
	} else {
		fmt.Printf("Failed to post comment. Status: %s\n", resp.Status)
	}
}
