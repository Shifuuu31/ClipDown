package main

// import (
// 	"fmt"
// 	"io"
// 	"log"
// 	"os"

// 	"github.com/kkdai/youtube/v2"
// )

// func downloadYouTubeVideo(videoURL string, format string, quality string) error {
// 	// Create a new YouTube client
// 	client := youtube.Client{}

// 	// Get the video information
// 	video, err := client.GetVideo(videoURL)
// 	if err != nil {
// 		return fmt.Errorf("failed to get video info: %w", err)
// 	}

// 	for i := 0; i < len(video.CaptionTracks); i++ {
// 		fmt.Println(video.Description)
// 	}
// 	os.Exit(0)
// 	var selectedIndex int
// 	// Debug: Print available formats
// 	fmt.Println("Available Formats:")
// 	for i, f := range video.Formats {
// 		fmt.Printf("[%d]	Format: %s	|	vQuality: %s	|	aQuality: %s\n", i+1, f.MimeType, f.QualityLabel, f.AudioQuality)
// 	}
// 	fmt.Print("Choose the quality that you prefer (enter the number): ")
// 	_, err = fmt.Scanf("%d", &selectedIndex) // Read an integer input
// 	if err != nil {
// 		log.Fatalf("Invalid input: %v", err)
// 	}
// 	fmt.Printf("\nselected: %d\n\n", selectedIndex) // Read an integer input

// 	// Validate the selected index
// 	if selectedIndex < 1 || selectedIndex > len(video.Formats) {
// 		log.Fatalf("Selected index [%d] is out of range. Please select a number between 1 and %d.", selectedIndex, len(video.Formats))
// 	}

// 	// Find the desired format
// 	var selectedFormat *youtube.Format = &video.Formats[selectedIndex-1]
// 	// for _, f := range video.Formats {
// 	// 	// Check if the base mime type matches and quality matches
// 	// 	if f.MimeType[:len(format)+5] == "video/"+format && f.QualityLabel == quality { // 5 is the length of `video/`
// 	// 		selectedFormat = &f
// 	// 		break
// 	// 	}
// 	// }

// 	if selectedFormat == nil {
// 		return fmt.Errorf("desired format and quality not found")
// 	}

// 	// Create the output file
// 	outputFileName := fmt.Sprintf("%s.%s", video.Title, selectedFormat.MimeType[6:10])
// 	file, err := os.Create(outputFileName)
// 	if err != nil {
// 		return fmt.Errorf("failed to create output file: %w", err)
// 	}
// 	defer file.Close()

// 	// Download the video
// 	stream, size, err := client.GetStream(video, selectedFormat)
// 	if err != nil {
// 		return fmt.Errorf("failed to get video stream: %w", err)
// 	}

// 	// Read from the stream and write to file
// 	_, err = io.Copy(file, stream)
// 	if err != nil {
// 		return fmt.Errorf("failed to download video: %w", err)
// 	}

// 	fmt.Printf("Downloaded: %s --> Size: %v\n", outputFileName, size)
// 	return nil
// }

// func main() {
// 	if len(os.Args) < 4 {
// 		log.Fatalf("Usage: %s <videoURL> <format> <quality>", os.Args[0])
// 	}
// 	videoURL := os.Args[1]
// 	format := os.Args[2]
// 	quality := os.Args[3]

// 	err := downloadYouTubeVideo(videoURL, format, quality)
// 	if err != nil {
// 		log.Fatalf("Error downloading video: %v", err)
// 	}
// }
