package server

import (
	"20MIS0315_Backend/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"google.golang.org/api/youtube/v3"
)

// Video represents the data structure for a video
type Video struct {
	VideoID      string `json:"video_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ChannelTitle string `json:"channel_title"`
	PublishedAt  string `json:"published_at"`
}

const apiKey = "AIzaSyDFsovXnSzsML2WbXdItm3NzCYzo8dBP8s"

// GetPaginatedVideosHandler handles the API request for paginated videos
func GetPaginatedVideosHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("pageSize")

	pageInt := 1
	pageSizeInt := 10
	var err error

	if page != "" {
		pageInt, err = strconv.Atoi(page)
		if err != nil || pageInt < 1 {
			http.Error(w, "Invalid page parameter", http.StatusBadRequest)
			return
		}
	}

	if pageSize != "" {
		pageSizeInt, err = strconv.Atoi(pageSize)
		if err != nil || pageSizeInt < 1 {
			http.Error(w, "Invalid pageSize parameter", http.StatusBadRequest)
			return
		}
	}

	offset := (pageInt - 1) * pageSizeInt

	log.Printf("Fetching videos with page: %d, pageSize: %d, offset: %d", pageInt, pageSizeInt, offset)

	rows, err := database.DB.Query(`
        SELECT id, title, description, publish_date, thumbnail_url 
        FROM videos 
        ORDER BY publish_date DESC 
        LIMIT $1 OFFSET $2`, pageSizeInt, offset)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, "Error fetching videos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var videos []map[string]interface{}
	for rows.Next() {
		var id, title, description, thumbnailURL string
		var publishDate string
		err := rows.Scan(&id, &title, &description, &publishDate, &thumbnailURL)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Error scanning video", http.StatusInternalServerError)
			return
		}

		log.Printf("Fetched video: %s, %s", title, publishDate)

		video := map[string]interface{}{
			"id":           id,
			"title":        title,
			"description":  description,
			"publishDate":  publishDate,
			"thumbnailUrl": thumbnailURL,
		}
		videos = append(videos, video)
	}

	if len(videos) == 0 {
		log.Println("No videos found")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

// SearchVideosHandler handles the search API requests
func SearchVideosHandler(w http.ResponseWriter, r *http.Request, service *youtube.Service) {
	query := r.URL.Query().Get("query")
	maxResults := int64(25)

	call := service.Search.List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(maxResults)
	response, err := call.Do()
	if err != nil {
		http.Error(w, "Error fetching videos", http.StatusInternalServerError)
		return
	}

	// Process the result and write it to the response
	videos := make([]map[string]interface{}, 0)
	for _, item := range response.Items {
		if item.Id.Kind == "youtube#video" {
			video := map[string]interface{}{
				"id":           item.Id.VideoId,
				"title":        item.Snippet.Title,
				"description":  item.Snippet.Description,
				"publishDate":  item.Snippet.PublishedAt,
				"thumbnailUrl": item.Snippet.Thumbnails.Default.Url,
			}
			videos = append(videos, video)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}
