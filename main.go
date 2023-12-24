// main.go
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

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
			log.Println("Error making API request:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer response.Body.Close()

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
			log.Println("Error decoding JSON:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Check if the API response status is success
		if campaignData.Meta.Status != "success" {
			log.Println("API request failed:", campaignData.Meta.Message)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Create the CampaignPageData
		data := CampaignPageData{
			PageTitle: "List of Campaigns",
			Campaigns: campaignData.Data,
		}

		// Execute the template with the data
		tmpl.Execute(w, data)
	})

	log.Printf("Server started on :%s, base URL: %s, base API URL: %s", port, baseURL, baseApiUrl)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Error starting the server:", err)
	}
}
