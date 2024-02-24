package lib

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

var PromptPrefix = "Your name is Himbot. You are a helpful but sarcastic and witty discord bot. Please respond with a natural response to the following prompt with that personality in mind:"

func OpenAITextGeneration(prompt string) (string, error) {
	godotenv.Load(".env")
	apiKey := os.Getenv("OPENAI_API_KEY")

	client := openai.NewClient(apiKey)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4Turbo1106,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: PromptPrefix + prompt,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("Ask command error: %v\n", err)
		return "", errors.New("https://fly.storage.tigris.dev/atridad/himbot/no.gif")
	}

	return resp.Choices[0].Message.Content, nil
}

func OpenAIImageGeneration(prompt string, filename string) (*bytes.Buffer, error) {
	godotenv.Load(".env")
	apiKey := os.Getenv("OPENAI_API_KEY")

	client := openai.NewClient(apiKey)

	imageResponse, err := client.CreateImage(
		context.Background(),
		openai.ImageRequest{
			Model:          openai.CreateImageModelDallE3,
			Prompt:         prompt,
			Size:           openai.CreateImageSize1024x1024,
			Quality:        openai.CreateImageQualityStandard,
			ResponseFormat: openai.CreateImageResponseFormatURL,
			N:              1,
		},
	)

	if err != nil {
		fmt.Printf("Pic command error: %v\n", err)
		return nil, errors.New("https://fly.storage.tigris.dev/atridad/himbot/hornypolice.gif")
	}

	imgUrl := imageResponse.Data[0].URL

	imageRes, imageGetErr := http.Get(imgUrl)
	if imageGetErr != nil {
		return nil, imageGetErr
	}

	defer imageRes.Body.Close()

	imageBytes, imgReadErr := io.ReadAll(imageRes.Body)
	if imgReadErr != nil {
		return nil, imgReadErr
	}

	// Save image to a temporary file
	tmpfile, err := os.Create(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(imageBytes); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	// Upload the image to S3
	_, uploadErr := UploadToS3(tmpfile.Name())
	if uploadErr != nil {
		log.Printf("Failed to upload image to S3: %v", uploadErr)
	}

	imageFile := bytes.NewBuffer(imageBytes)
	return imageFile, nil
}
