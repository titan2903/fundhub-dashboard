package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Campaign struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	ImageURL         string `json:"image_url"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
	Slug             string `json:"slug"`
}

type CampaignPageData struct {
	PageTitle string
	Campaigns []Campaign
}

func main() {
	_ = godotenv.Load()
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "/"
	}

	baseApiUrl := os.Getenv("BASE_API_URL")
	if baseApiUrl == "" {
		baseApiUrl = "http://fundhub.api.local"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html"))

	http.HandleFunc(baseURL, func(w http.ResponseWriter, r *http.Request) {
		// Make a GET request to the API endpoint
		apiURL := fmt.Sprintf("%s/%s", baseApiUrl, "api/v1/campaigns")
		response, err := http.Get(apiURL)
		if err != nil {
			log.Error("Error making API request:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer response.Body.Close()

		log.Infof("Response Status Code: %d", response.StatusCode)
		log.Infof("Response Headers: %+v", response.Header)

		// Check the Content-Type
		contentType := response.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "application/json") {
			log.Errorf("Unexpected content type: %s", contentType)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// Check for an empty response body
		if response.ContentLength == 0 {
			log.Warn("Empty response body")
			http.Error(w, "Empty response body", http.StatusInternalServerError)
			return
		}

		// Read and log the response body
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			log.Error("Error reading response body:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		log.Infof("Response Body: %s", responseBody)

		// Create a new reader for the response body
		response.Body = ioutil.NopCloser(bytes.NewBuffer(responseBody))

		// Decode the JSON response into the CampaignPageData struct
		var campaignData struct {
			Meta struct {
				Message string `json:"message"`
				Code    int    `json:"code"`
				Status  string `json:"status"`
			} `json:"meta"`
			Data []Campaign `json:"data"`
		}

		err = json.NewDecoder(response.Body).Decode(&campaignData)
		if err != nil {
			log.Error("Error decoding JSON:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Check if the API response status is success
		if campaignData.Meta.Status != "success" {
			log.Error("API request failed:", campaignData.Meta.Message)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Create the CampaignPageData
		data := CampaignPageData{
			PageTitle: "List of Campaigns",
			Campaigns: campaignData.Data,
		}

		tmpl.Execute(w, data)
	})

	log.Infof("Server started on :%s, base URL: %s, base API URL: %s", port, baseURL, baseApiUrl)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Error("Error starting the server:", err)
	}
}
