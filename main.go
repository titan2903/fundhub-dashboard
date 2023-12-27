package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"

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

	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://fundhub.api.local/api/v1/campaigns"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html"))

	http.HandleFunc(baseURL, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Make a GET request to the API endpoint
		response, err := http.Get(apiURL)
		if err != nil {
			log.Errorf("Error making API request: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer response.Body.Close()

		log.Infof("Response Status Code: %d", response.StatusCode)
		log.Infof("Response Headers: %+v", response.Header)

		// Decode the JSON response into the CampaignPageData struct
		var campaignData struct {
			Meta struct {
				Message string `json:"message"`
				Code    int    `json:"code"`
				Status  string `json:"status"`
			} `json:"meta"`
			Data []Campaign `json:"data"`
		}

		log.Infof("body response: %v", response.Body)

		_ = json.NewDecoder(response.Body).Decode(&campaignData)
		// if err != nil {
		// 	log.Errorf("Error decoding JSON: %v", err)
		// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		// 	return
		// }

		// Create the CampaignPageData
		data := CampaignPageData{
			PageTitle: "List of Campaigns",
			Campaigns: campaignData.Data,
		}

		tmpl.Execute(w, data)
	})

	log.Infof("Server started on :%s, base URL: %s, base API URL: %s", port, baseURL, apiURL)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Errorf("Error starting the server: %v", err)
	}
}
