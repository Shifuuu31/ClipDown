package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/kkdai/youtube/v2"
)

type YtVid struct {
	Title          string
	Thumbnail      string
	Duration       string
	Description    string
	Mp3, Mp4, Webm []*youtube.Format
}

var (
	tmpl   = template.Must(template.ParseGlob("static/*.html"))
	Data   YtVid
	Client youtube.Client
	Video  *youtube.Video
	mu     sync.Mutex // Mutex for handling concurrency
)

// Root handler for the main page
func rootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "home.html", nil)
}

func detailHandler(w http.ResponseWriter, r *http.Request) {
	videoUrl := r.URL.Query().Get("videoUrl")
	if videoUrl == "" {
		http.Error(w, "No video URL provided", http.StatusBadRequest)
		return
	}

	video, err := Client.GetVideo(videoUrl)
	if err != nil {
		http.Error(w, "Failed to retrieve video", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Set global Video and Data variables
	Video = video
	Data.Title = video.Title
	Data.Duration = video.Duration.String()
	Data.Description = video.Description
	Data.Thumbnail = video.Thumbnails[0].URL

	// Find the thumbnail with "mq"
	for i := 0; i < len(video.Thumbnails); i++ {
		if strings.Contains(video.Thumbnails[i].URL, "mq") {
			Data.Thumbnail = video.Thumbnails[i].URL
		}
	}

	// Categorize the formats
	for i := 0; i < len(video.Formats); i++ {
		// video.Formats[i].ContentLength /= (1024 * 1024) // Convert size to MB
		if strings.Contains(video.Formats[i].MimeType, "audio") {
			if video.Formats[i].AudioQuality == "AUDIO_QUALITY_LOW" {
				video.Formats[i].AudioQuality = "Low"
			}
			if video.Formats[i].AudioQuality == "AUDIO_QUALITY_MEDIUM" {
				video.Formats[i].AudioQuality = "Medium"
			}
			if video.Formats[i].AudioQuality == "AUDIO_QUALITY_HIGH" {
				video.Formats[i].AudioQuality = "High"
			}
			Data.Mp3 = append(Data.Mp3, &video.Formats[i])
		} else if strings.Contains(video.Formats[i].MimeType, "video/mp4") {
			Data.Mp4 = append(Data.Mp4, &video.Formats[i])
		} else if strings.Contains(video.Formats[i].MimeType, "video/webm") {
			Data.Webm = append(Data.Webm, &video.Formats[i])
		}
	}

	// Render download page
	tmpl.ExecuteTemplate(w, "download.html", Data)
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/details", detailHandler)
	http.HandleFunc("/downloadAudio", downloadMp3)
	http.HandleFunc("/downloadMp4", downloadMp4)
	http.HandleFunc("/downloadWebm", downloadWebm)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// DownloadMp3 handles MP3 downloads
func downloadMp3(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/downloadAudio" {
		http.Error(w, "404 Page Not Found", http.StatusNotFound)
		return
	}

	in := r.URL.Query().Get("quality-size")
	if in == "" {
		http.Error(w, "Download failed, try again", http.StatusBadRequest)
		return
	}

	// Extract quality and size from the query string
	quality := strings.Split(in, "-")[0]
	size, err := strconv.Atoi(strings.Split(in, "-")[1])
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid size: %v", err), http.StatusBadRequest)
		return
	}

	// Lock the mutex to handle concurrency
	mu.Lock()
	defer mu.Unlock()

	// Iterate through available MP3 formats
	for _, format := range Data.Mp3 {
		if format.AudioQuality == quality && format.ContentLength == int64(size) {
			err := Download(w, format, "mp3")
			if err != nil {
				http.Error(w, fmt.Sprintf("Download failed: %v", err), http.StatusInternalServerError)
			}
			return
		}
	}

	http.Error(w, "No matching MP3 found", http.StatusBadRequest)
}

// DownloadMp4 handles MP4 downloads
func downloadMp4(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/downloadMp4" {
		http.Error(w, "404 Page Not Found", http.StatusNotFound)
		return
	}

	in := r.URL.Query().Get("quality-size")
	if in == "" {
		http.Error(w, "Download failed, try again", http.StatusBadRequest)
		return
	}

	// Extract quality and size from the query string
	quality := strings.Split(in, "-")[0]
	size, err := strconv.Atoi(strings.Split(in, "-")[1])
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid size: %v", err), http.StatusBadRequest)
		return
	}

	// Lock the mutex to handle concurrency
	mu.Lock()
	defer mu.Unlock()

	// Iterate through available MP4 formats
	for _, format := range Data.Mp4 {
		if format.QualityLabel == quality && format.ContentLength == int64(size) {
			err := Download(w, format, "mp4")
			if err != nil {
				http.Error(w, fmt.Sprintf("Download failed: %v", err), http.StatusInternalServerError)
			}
			return
		}
	}

	http.Error(w, "No matching MP4 found", http.StatusBadRequest)
}

// DownloadWebm handles WEBM downloads
func downloadWebm(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/downloadWebm" {
		http.Error(w, "404 Page Not Found", http.StatusNotFound)
		return
	}

	in := r.URL.Query().Get("quality-size")
	if in == "" {
		http.Error(w, "Download failed, try again", http.StatusBadRequest)
		return
	}

	// Extract quality and size from the query string
	quality := strings.Split(in, "-")[0]
	size, err := strconv.Atoi(strings.Split(in, "-")[1])
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid size: %v", err), http.StatusBadRequest)
		return
	}

	// Lock the mutex to handle concurrency
	mu.Lock()
	defer mu.Unlock()

	// Iterate through available WEBM formats
	for _, format := range Data.Webm {
		if format.QualityLabel == quality && format.ContentLength == int64(size) {
			err := Download(w, format, "webm")
			if err != nil {
				http.Error(w, fmt.Sprintf("Download failed: %v", err), http.StatusInternalServerError)
			}
			return
		}
	}

	http.Error(w, "No matching WEBM found", http.StatusBadRequest)
}

// Download helper function for downloading a video file
func Download(w http.ResponseWriter, format *youtube.Format, ext string) error {
	// Ensure that Video is initialized
	if Video == nil {
		return fmt.Errorf("Video is not initialized")
	}

	// Set headers for the file download
	outputFileName := fmt.Sprintf("%s.%s", sanitizeFileName(Video.Title), ext)
	w.Header().Set("Content-Type", format.MimeType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", outputFileName))
	w.Header().Set("Content-Length", strconv.FormatInt(format.ContentLength, 10))

	// Get the stream for the video format
	stream, size, err := Client.GetStream(Video, format)
	if err != nil {
		return fmt.Errorf("failed to get video stream: %w", err)
	}

	// Ensure stream and size are valid
	if stream == nil || size == 0 {
		return fmt.Errorf("invalid video stream or size")
	}

	// Read the stream into the response writer
	buf := make([]byte, 1024*16) // Buffer size
	totalBytesRead := int64(0)

	for {
		n, err := stream.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read video stream: %w", err)
		}
		if n == 0 {
			break
		}

		_, err = w.Write(buf[:n])
		if err != nil {
			return fmt.Errorf("failed to write data to response: %w", err)
		}

		// Update the number of bytes written
		totalBytesRead += int64(n)
	}

	return nil
}

// Helper function to sanitize filenames
func sanitizeFileName(name string) string {
	// Remove unwanted characters from the filename
	return strings.ReplaceAll(name, " ", "_")
}
