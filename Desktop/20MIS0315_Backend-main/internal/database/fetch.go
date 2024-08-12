package database

// Video represents a video entry
type Video struct {
	VideoID      string
	Title        string
	Description  string
	ChannelTitle string
	PublishedAt  string
}

// GetVideos retrieves videos from the database
func GetVideos() ([]Video, error) {
	rows, err := DB.Query("SELECT video_id, title, description, channel_title, published_at FROM videos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []Video
	for rows.Next() {
		var video Video
		if err := rows.Scan(&video.VideoID, &video.Title, &video.Description, &video.ChannelTitle, &video.PublishedAt); err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return videos, nil
}
