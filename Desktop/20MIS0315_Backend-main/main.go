package main

import (
	"20MIS0315_Backend/internal/database"
	"20MIS0315_Backend/internal/server"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type YouTubeResponse struct {
	Kind          string   `json:"kind"`
	Etag          string   `json:"etag"`
	NextPageToken string   `json:"nextPageToken"`
	RegionCode    string   `json:"regionCode"`
	PageInfo      PageInfo `json:"pageInfo"`
	Items         []Item   `json:"items"`
}

type PageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}

type Item struct {
	Kind    string  `json:"kind"`
	Etag    string  `json:"etag"`
	ID      ID      `json:"id"`
	Snippet Snippet `json:"snippet"`
}

type ID struct {
	Kind    string `json:"kind"`
	VideoID string `json:"videoId"`
}

type Snippet struct {
	PublishedAt          string     `json:"publishedAt"`
	ChannelID            string     `json:"channelId"`
	Title                string     `json:"title"`
	Description          string     `json:"description"`
	Thumbnails           Thumbnails `json:"thumbnails"`
	ChannelTitle         string     `json:"channelTitle"`
	LiveBroadcastContent string     `json:"liveBroadcastContent"`
	PublishTime          string     `json:"publishTime"`
}

type Thumbnails struct {
	Default Thumbnail `json:"default"`
	Medium  Thumbnail `json:"medium"`
	High    Thumbnail `json:"high"`
}

type Thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

const developerKey = "AIzaSyDFsovXnSzsML2WbXdItm3NzCYzo8dBP8"

func main() {
	flag.Parse()

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the database
	database.InitDB()

	// Load JSON data from a file
	data, err := ioutil.ReadFile("data.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	var youTubeResponse YouTubeResponse

	// Parse JSON data
	err = json.Unmarshal(data, &youTubeResponse)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	// Insert data into PostgreSQL
	for _, item := range youTubeResponse.Items {
		video := item.Snippet
		err = database.InsertVideo(item.ID.VideoID, video.Title, video.Description, video.ChannelTitle, video.PublishedAt)
		if err != nil {
			log.Printf("Error inserting video: %v", err)
		}
	}

	log.Println("Data inserted successfully!")

	// Initialize the YouTube client
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	// Initialize the router
	r := mux.NewRouter()

	// Define API routes
	r.HandleFunc("/api/videos", server.GetPaginatedVideosHandler).Methods("GET")
	r.HandleFunc("/api/videos/search", func(w http.ResponseWriter, r *http.Request) {
		server.SearchVideosHandler(w, r, service)
	}).Methods("GET")

	// Serve static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public")))

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
