package youtube

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"20MIS0315_Backend/internal/database"
)

const apiKey = "AIzaSyDFsovXnSzsML2WbXdItm3NzCYzo8dBP8"

// Video represents a video entry
type Video struct {
	VideoID      string
	Title        string
	Description  string
	ChannelTitle string
	PublishDate  time.Time
	ThumbnailURL string
}

// FetchAndStoreVideos fetches videos from YouTube and stores them in the database
func FetchAndStoreVideos(query string, publishedAfter string) {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&order=date&type=video&q=%s&publishedAfter=%s&key=%s", query, publishedAfter, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching data from YouTube API: %v", err)
		return
	}
	defer resp.Body.Close()

	var apiResponse struct {
		Items []struct {
			ID struct {
				VideoID string `json:"videoId"`
			} `json:"id"`
			Snippet struct {
				PublishedAt  string `json:"publishedAt"`
				Title        string `json:"title"`
				Description  string `json:"description"`
				ChannelTitle string `json:"channelTitle"`
				Thumbnails   struct {
					Default struct {
						URL string `json:"url"`
					} `json:"default"`
				} `json:"thumbnails"`
			} `json:"snippet"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		log.Printf("Error decoding YouTube API response: %v", err)
		return
	}

	for _, item := range apiResponse.Items {
		video := Video{
			VideoID:      item.ID.VideoID,
			Title:        item.Snippet.Title,
			Description:  item.Snippet.Description,
			ChannelTitle: item.Snippet.ChannelTitle,
			PublishDate:  parseTime(item.Snippet.PublishedAt),
			ThumbnailURL: item.Snippet.Thumbnails.Default.URL,
		}
		if err := storeVideo(video); err != nil {
			log.Printf("Error storing video: %v", err)
		}
	}
}

// parseTime parses the YouTube publish date into a time.Time object
func parseTime(dateStr string) time.Time {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		log.Printf("Error parsing time: %v", err)
		return time.Time{}
	}
	return t
}

// storeVideo inserts a video into the PostgreSQL database
func storeVideo(video Video) error {
	query := `
        INSERT INTO videos (video_id, title, description, channel_title, published_at, thumbnail_url)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (video_id) DO NOTHING;
    `
	_, err := database.DB.Exec(query, video.VideoID, video.Title, video.Description, video.ChannelTitle, video.PublishDate, video.ThumbnailURL)
	return err
}
