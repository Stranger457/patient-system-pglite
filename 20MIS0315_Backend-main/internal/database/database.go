package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error

	connStr := "user=vallipichowdappa dbname=postgres sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database connection:", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}
	log.Println("Database connection established")
}

func InsertVideo(videoID, title, description, channelTitle, publishedAt string) error {
	query := `
        INSERT INTO videos (video_id, title, description, channel_title, published_at)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (video_id) DO NOTHING;
    `
	_, err := DB.Exec(query, videoID, title, description, channelTitle, publishedAt)
	return err
}
