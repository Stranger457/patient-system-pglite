package server

import (
	"20MIS0315_Backend/internal/database"
	"encoding/json"
	"net/http"
)

// Video represents a video entry
type Video struct {
	VideoID      string `json:"video_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ChannelTitle string `json:"channel_title"`
	PublishedAt  string `json:"published_at"`
}

// InsertVideoHandler handles POST requests to insert video data
func InsertVideoHandler(w http.ResponseWriter, r *http.Request) {
	var video Video
	// Decode the JSON request body into a Video struct
	err := json.NewDecoder(r.Body).Decode(&video)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Insert the video data into the database
	err = database.InsertVideo(video.VideoID, video.Title, video.Description, video.ChannelTitle, video.PublishedAt)
	if err != nil {
		http.Error(w, "Error inserting video into database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Video inserted successfully"))
}
