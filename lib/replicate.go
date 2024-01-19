package lib

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/replicate/replicate-go"
)

func ReplicateTextGeneration(prompt string) (string, error) {
	client, clientError := replicate.NewClient(replicate.WithTokenFromEnv())
	if clientError != nil {
		return "", clientError
	}

	input := replicate.PredictionInput{
		"prompt":         prompt,
		"max_new_tokens": 4096,
	}
	webhook := replicate.Webhook{
		URL:    "https://example.com/webhook",
		Events: []replicate.WebhookEventType{"start", "completed"},
	}

	prediction, predictionError := client.Run(context.Background(), "mistralai/mistral-7b-instruct-v0.2:79052a3adbba8116ebc6697dcba67ad0d58feff23e7aeb2f103fc9aa545f9269", input, &webhook)

	if predictionError != nil {
		return "", predictionError
	}

	test, ok := prediction.([]interface{})

	if !ok {
		return "", errors.New("there was an error generating the image based on this prompt... this usually happens when the generated image violates safety requirements")
	}

	strs := make([]string, len(test))
	for i, v := range test {
		str, ok := v.(string)
		if !ok {
			return "", errors.New("element is not a string")
		}
		strs[i] = str
	}

	result := strings.Join(strs, "")

	return result, nil
}

func ReplicateImageGeneration(prompt string, filename string) (*bytes.Buffer, error) {
	client, clientError := replicate.NewClient(replicate.WithTokenFromEnv())
	if clientError != nil {
		return nil, clientError
	}

	input := replicate.PredictionInput{
		"prompt":                 prompt,
		"refiner":                "expert_ensemble_refiner",
		"num_inference_steps":    69,
		"disable_safety_checker": true,
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

	// Save image to a temporary file
	var tmpfile *os.File
	var err error

	if filename != "" {
		tmpfile, err = os.CreateTemp("", filename)
	} else {
		tmpfile, err = os.CreateTemp("", "image.*.jpg")
	}

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
