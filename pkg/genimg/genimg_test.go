package genimg_test

import (
	"context"
	"errors"
	"fmt"
	"gen-img/pkg/genimg"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func loadApiKey() (string, error) {
	err := godotenv.Load("../../.env")

	if err != nil {
		return "", err
	}

	apiKey := os.Getenv("OPENAI_API_KEY")

	if apiKey == "" {
		return "", errors.New("OPENAI_API_KEY is not set")
	} else {
		return apiKey, nil
	}
}

func TestLoadOpenaiAPiKey(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Errorf("error loading .env file: %v", err)
	}
}

func TestGenerateImage(t *testing.T) {

	apiKey, err := loadApiKey()

	if err != nil {
		t.Errorf("error loading api key: %v", err)
	}

	imageUrl, err := genimg.GenerateImage(genimg.OpenaiImagesApiEndpoint, apiKey, "a robot painting on a canvas", genimg.ImageSizeLandscape)

	if err != nil {
		t.Errorf("error generating image: %v", err)
	}

	fmt.Println(imageUrl)
}

func TestGenerateImages(t *testing.T) {
	apiKey, err := loadApiKey()

	if err != nil {
		t.Errorf("error loading api key: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	imageUrls := genimg.GenerateImages(ctx, genimg.OpenaiImagesApiEndpoint, apiKey, "a robot painting on a canvas", genimg.ImageSizeLandscape, 10, 10)

	for _, imageUrl := range imageUrls {
		fmt.Println(imageUrl)
	}
}
