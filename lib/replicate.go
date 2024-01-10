package lib

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/replicate/replicate-go"
)

func ReplicateImageGeneration(prompt string) (*bytes.Buffer, error) {
	client, clientError := replicate.NewClient(replicate.WithTokenFromEnv())
	if clientError != nil {
		return nil, clientError
	}

	input := replicate.PredictionInput{
		"prompt": prompt,
	}
	webhook := replicate.Webhook{
		URL:    "https://example.com/webhook",
		Events: []replicate.WebhookEventType{"start", "completed"},
	}

	prediction, predictionError := client.Run(context.Background(), "stability-ai/sdxl:39ed52f2a78e934b3ba6e2a89f5b1c712de7dfea535525255b1aa35c5565e08b", input, &webhook)

	if predictionError != nil {
		return nil, predictionError
	}

	test, ok := prediction.([]interface{})

	if !ok {
		return nil, errors.New("there was an error generating the image based on this prompt... this usually happens when the generated image violates safety requirements")
	}

	imgUrl, ok := test[0].(string)

	if !ok {
		return nil, errors.New("there was an error generating the image based on this prompt... this usually happens when the generated image violates safety requirements")
	}

	imageRes, imageGetErr := http.Get(imgUrl)
	if imageGetErr != nil {
		return nil, imageGetErr
	}

	defer imageRes.Body.Close()

	imageBytes, imgReadErr := io.ReadAll(imageRes.Body)
	if imgReadErr != nil {
		return nil, imgReadErr
	}

	imageFile := bytes.NewBuffer(imageBytes)
	return imageFile, nil
}
