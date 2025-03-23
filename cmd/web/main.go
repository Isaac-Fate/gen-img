package main

import (
	"context"
	"gen-img/internal/templs/components"
	"gen-img/internal/templs/pages"
	"gen-img/pkg/genimg"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const numWorkers = 10

func main() {

	// Load the dotenv
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Get the API key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		panic("OPENAI_API_KEY is not set")
	}

	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		pages.Home().Render(ctx.Request.Context(), ctx.Writer)
	})

	router.POST("/generate", func(ctx *gin.Context) {

		endpoint := ctx.PostForm("endpoint")
		prompt := ctx.PostForm("prompt")
		imageSize := ctx.PostForm("imageSize")

		numImagesAsString := ctx.PostForm("numImages")
		numImages, err := strconv.Atoi(numImagesAsString)
		if err != nil {
			panic(err)
		}

		c, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		images := genimg.GenerateImages(c, endpoint, apiKey, prompt, imageSize, numImages, numWorkers)

		components.ImageList(images).Render(ctx.Request.Context(), ctx.Writer)
	})

	router.Run()
}
