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

var ReplicatePromptPrefix = "Your name is Himbot. You are a helpful but sarcastic and witty discord bot. Please respond with a natural response to the following prompt with that personality in mind:"

func ReplicateTextGeneration(prompt string) (string, error) {
	client, clientError := replicate.NewClient(replicate.WithTokenFromEnv())
	if clientError != nil {
		return "", clientError
	}

	input := replicate.PredictionInput{
		"prompt":         ReplicatePromptPrefix + prompt,
		"max_new_tokens": 4096,
	}

	webhook := replicate.Webhook{
		URL:    "https://example.com/webhook",
		Events: []replicate.WebhookEventType{"start", "completed"},
	}

	prediction, predictionError := client.Run(context.Background(), "meta/llama-2-70b-chat:44878be272d384b00155fe821d8534275d5471b88c25dc6e78c97ce3644b1e99", input, &webhook)

	if predictionError != nil {
		return "", predictionError
	}

	if prediction == nil {
		return "", errors.New("there was an error generating a response based on this prompt... please reach out to @himbothyswaggins to fix this issue")
	}

	test, ok := prediction.([]interface{})

	if !ok {
		return "", errors.New("there was an error generating a response based on this prompt... please reach out to @himbothyswaggins to fix this issue")
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
		"width":                  1024,
		"height":                 1024,
		"prompt":                 prompt,
		"scheduler":              "K_EULER",
		"refine":                 "expert_ensemble_refiner",
		"num_outputs":            1,
		"negative_prompt":        "ugly, deformed, noisy, blurry, low contrast, text, BadDream, lowres, low resolution, mutated body parts, extra limbs, mutated body parts, inaccurate hands, too many hands, deformed fingers, too many fingers, deformed eyes, deformed faces, unrealistic faces",
		"num_inference_steps":    35,
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

	if prediction == nil {
		return nil, errors.New("there was an error generating the image based on this prompt... please reach out to @himbothyswaggins to fix this issue")
	}

	test, ok := prediction.([]interface{})

	if !ok {
		return nil, errors.New("there was an error generating the image based on this prompt... please reach out to @himbothyswaggins to fix this issue")
	}

	imgUrl, ok := test[0].(string)

	if !ok {
		return nil, errors.New("there was an error generating the image based on this prompt... please reach out to @himbothyswaggins to fix this issue")
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
