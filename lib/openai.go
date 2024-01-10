package lib

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

var client *openai.Client

func init() {
	godotenv.Load(".env")
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY environment variable not set")
		os.Exit(1)
	}
	client = openai.NewClient(apiKey)
}

func OpenAITextGeneration(prompt string) (string, error) {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4TurboPreview,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func OpenAIImageGeneration(prompt string) (imageFile *bytes.Buffer, err error) {
	resp, err := client.CreateImage(context.Background(), openai.ImageRequest{
		Prompt: prompt,
		Model:  "dall-e-3",
		Size:   "1024x1024",
	})

	if err != nil {
		return nil, errors.New("there was an error generating the image based on this prompt... this usually happens when the generated image violates safety requirements")
	}

	imageRes, err := http.Get(resp.Data[0].URL)

	if err != nil {
		return nil, err
	}

	defer imageRes.Body.Close()

	imageBytes, err := io.ReadAll(imageRes.Body)

	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(imageBytes), nil
}
