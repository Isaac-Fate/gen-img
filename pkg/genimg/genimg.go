package genimg

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
)

const (
	OpenaiImagesApiEndpoint = "https://api.openai.com/v1/images/generations"
	ImageSizeSquare         = "1024x1024"
	ImageSizeLandscape      = "1792x1024"
)

type imageGenerationResponseData struct {
	Data []struct {
		Url string `json:"url"`
	}
}

func GenerateImage(endpoint string, apiKey string, prompt string, imageSize string) (string, error) {

	// Create the payload
	payload := map[string]any{
		"model":           "dall-e-3",
		"prompt":          prompt,
		"size":            imageSize,
		"response_format": "url",
	}

	// Convert to JSON
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Create the request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payloadJson))
	if err != nil {
		return "", err
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// Read the response body

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}

	// Decode the JSON response into a struct
	var responseData imageGenerationResponseData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return "", err
	}

	// Get the image URL

	if len(responseData.Data) == 0 {
		return "", errors.New("data is empty")
	}

	imageUrl := responseData.Data[0].Url

	return imageUrl, nil
}

func GenerateImages(ctx context.Context, endpoint string, apiKey string, prompt string, imageSize string, numImages int, maxWorkers int) []string {

	// A buffered channel of all image URLs to return
	imageUrlChan := make(chan string, numImages)

	// A buffered channel of all requests to send
	jobChan := make(chan struct{}, numImages)

	// Create a wait goup
	var wg sync.WaitGroup

	for range maxWorkers {

		// Increment the wait group
		wg.Add(1)

		// Start a worker
		go worker(ctx, &wg, endpoint, apiKey, prompt, imageSize, jobChan, imageUrlChan)

	}

	// Send the jobs
	for range numImages {
		jobChan <- struct{}{}
	}

	// Done sending jobs
	close(jobChan)

	go func() {
		// Wait for the wait group to finish
		wg.Wait()

		// Close the channel
		close(imageUrlChan)
	}()

	// Slice of all URLs of generated images
	var imageUrls []string
	for imageUrl := range imageUrlChan {
		imageUrls = append(imageUrls, imageUrl)
	}

	return imageUrls
}

func worker(ctx context.Context, wg *sync.WaitGroup, endpoint string, apiKey string, prompt string, imageSize string, jobChan <-chan struct{}, imageUrlChan chan<- string) {

	// Decrement the wait group counter
	defer wg.Done()

	for range jobChan {
		select {

		case <-ctx.Done():
			return

		default:
			// Generate a single image
			imageUrl, err := GenerateImage(endpoint, apiKey, prompt, imageSize)

			// Collect the image URL
			if err == nil {
				imageUrlChan <- imageUrl
			}

		}
	}

}
