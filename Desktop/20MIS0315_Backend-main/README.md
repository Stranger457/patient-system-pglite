# 20MIS0315_Backend

#Project Goal

To make an API to fetch latest videos sorted in reverse chronological order of their publishing date-time from YouTube for a given tag/search query in a paginated response.


**#Basic Requirements:**

●  Server should call the YouTube API continuously in background with some interval (say 10 seconds but fetch more data points per call as there is a limit to number of calls you can make per day) for fetching the latest videos for a predefined search query and should store the data of videos (specifically these fields - Video title, description, publishing datetime, thumbnails URLs and any other fields you require) in a database with proper indexes , Make sure you have at least 500k data points stored , try to have more than this it would be great .


●  A GET API which returns the stored video data in a paginated response sorted in descending order of published datetime.


●  A basic search API to search the stored videos using their title and description and other filters which seem appropriate to you.

●  Dockerize the project.


●  It should be optimized.

**#YOUTUBE-API-FETCHER**

Description:

This project provides a comprehensive solution for fetching and managing YouTube video data. It consists of a backend service implemented in Go that uses the YouTube Data API to retrieve and store video details. Key features include:

• YouTube Video Fetching: Efficiently fetches the latest YouTube videos based on search queries and stores them in a PostgreSQL database.


• Pagination and Search: Implements pagination for fetching large sets of video data and allows searching through video titles and descriptions using partial matches.


• Dashboard: A user-friendly web interface for viewing and interacting with stored video data. The dashboard supports search functionality, sorting, and filtering options. The project is designed to be scalable and optimized for performance, with features including support for multiple API keys and data storage optimization. It is dockerized for easy deployment and includes a basic search API to enhance usability.

**#INSTALLATION:**

Clone the repository: 

git clone https://github.com/vallipichowdappa/youtube-api-fetcher


add keys : 

export youtube_apikeys = “your first api,second api,third api”

**#SET ENVIRONMENT**

install go(golang)

download docker,

postman/curl,

install postgresql

api key,

**#DEVELOPMENT**

To run the application docker-compose up,

To again build the image docker-compose build ,

To stop the application docker-compose down

start postgresql database • run main file “ go run main.go”

Navigate to http://localhost:8080/

Run the api-end points in postman/curl

Run all the tests. After running the tests :

API-Test 1 : http://localhost:8080/api/videos/search?query=tea.


Image ==>  https://drive.google.com/file/d/1AJakY5iz2zvJ6HovOQvKJpKtk80ZHR1r/view?usp=sharing


API-Test 2 : on cricket with youtube image : http://localhost:8080/api/videos/search?query=cricket


Image ==>  https://drive.google.com/file/d/1ZOtXKdzBnMehLikwlTukiGLN9vhZ6mhg/view?usp=sharing
