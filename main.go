package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"himbot/lib"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/joho/godotenv"
	"github.com/replicate/replicate-go"
	openai "github.com/sashabaranov/go-openai"
)

var commands = []api.CreateCommandData{
	{
		Name:        "ping",
		Description: "ping pong!",
	},
	{
		Name:        "ask",
		Description: "Ask Himbot! Cooldown: 1 Minute.",
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  "prompt",
				Description: "The prompt to send to Himbot.",
				Required:    true,
			},
		},
	},
	{
		Name:        "pic",
		Description: "Generate an image using Stable Diffusion! Cooldown: 1 Minute.",
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  "prompt",
				Description: "The prompt for the image generation.",
				Required:    true,
			},
		},
	},
	{
		Name:        "hdpic",
		Description: "Generate an image using DALL·E 3! Cooldown: 10 Minutes.",
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  "prompt",
				Description: "The prompt for the image generation.",
				Required:    true,
			},
		},
	},
	{
		Name:        "hs",
		Description: "This command was your nickname in highschool!",
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  "nickname",
				Description: "Your nickname in highschool.",
				Required:    true,
			},
		},
	},
}

func main() {
	godotenv.Load(".env")

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatalln("No $DISCORD_TOKEN given.")
	}

	h := newHandler(state.New("Bot " + token))
	h.s.AddInteractionHandler(h)
	h.s.AddIntents(gateway.IntentGuilds)
	h.s.AddHandler(func(*gateway.ReadyEvent) {
		me, _ := h.s.Me()
		log.Println("connected to the gateway as", me.Tag())
	})

	if err := cmdroute.OverwriteCommands(h.s, commands); err != nil {
		log.Fatalln("cannot update commands:", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := h.s.Connect(ctx); err != nil {
		log.Fatalln("cannot connect:", err)
	}
}

type handler struct {
	*cmdroute.Router
	s *state.State
}

func newHandler(s *state.State) *handler {
	h := &handler{s: s}

	h.Router = cmdroute.NewRouter()
	// Automatically defer handles if they're slow.
	h.Use(cmdroute.Deferrable(s, cmdroute.DeferOpts{}))
	h.AddFunc("ping", h.cmdPing)
	h.AddFunc("ask", h.cmdAsk)
	h.AddFunc("pic", h.cmdPic)
	h.AddFunc("hdpic", h.cmdHDPic)
	h.AddFunc("hs", h.cmdHS)

	return h
}

func (h *handler) cmdPing(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	// Command Logic
	return &api.InteractionResponseData{
		Content: option.NewNullableString("Pong!"),
	}
}

func (h *handler) cmdAsk(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	// Cooldown Logic
	allowed := lib.CooldownHandler(*data.Event)

	if !allowed {
		return errorResponse(errors.New("please wait for the cooldown"))
	}

	// Command Logic
	var options struct {
		Arg string `discord:"prompt"`
	}

	if err := data.Options.Unmarshal(&options); err != nil {
		return errorResponse(err)
	}

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4TurboPreview,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: options.Arg,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return &api.InteractionResponseData{
			Content:         option.NewNullableString("ChatCompletion Error!"),
			AllowedMentions: &api.AllowedMentions{}, // don't mention anyone
		}
	}

	return &api.InteractionResponseData{
		Content:         option.NewNullableString(resp.Choices[0].Message.Content),
		AllowedMentions: &api.AllowedMentions{}, // don't mention anyone
	}
}

func (h *handler) cmdPic(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	// Cooldown Logic
	allowed := lib.CooldownHandler(*data.Event)

	if !allowed {
		return errorResponse(errors.New("please wait for the cooldown"))
	}

	// Command Logic
	var options struct {
		Prompt string `discord:"prompt"`
	}

	if err := data.Options.Unmarshal(&options); err != nil {
		return errorResponse(err)
	}

	client, clientError := replicate.NewClient(replicate.WithTokenFromEnv())
	if clientError != nil {
		return errorResponse(clientError)
	}
	if err := data.Options.Unmarshal(&options); err != nil {
		return errorResponse(err)
	}

	input := replicate.PredictionInput{
		"prompt": options.Prompt,
	}
	webhook := replicate.Webhook{
		URL:    "https://example.com/webhook",
		Events: []replicate.WebhookEventType{"start", "completed"},
	}

	prediction, predictionError := client.Run(context.Background(), "stability-ai/sdxl:39ed52f2a78e934b3ba6e2a89f5b1c712de7dfea535525255b1aa35c5565e08b", input, &webhook)

	if predictionError != nil {
		return errorResponse(predictionError)
	}

	test, ok := prediction.([]interface{})

	if !ok {
		fmt.Println("prediction is not []interface{}")
	}

	imgUrl, ok := test[0].(string)

	if !ok {
		fmt.Println("prediction.Output[0] is not a string")
	}

	imageRes, imageGetErr := http.Get(imgUrl)
	if imageGetErr != nil {
		return errorResponse(imageGetErr)
	}

	defer imageRes.Body.Close()

	imageBytes, imgReadErr := io.ReadAll(imageRes.Body)
	if imgReadErr != nil {
		return errorResponse(imgReadErr)
	}

	imageFile := bytes.NewBuffer(imageBytes)

	file := sendpart.File{
		Name:   "himbot_response.png",
		Reader: imageFile,
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString("Prompt: " + options.Prompt),
		Files:   []sendpart.File{file},
	}
}

func (h *handler) cmdHDPic(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	// Cooldown Logic
	allowed := lib.CooldownHandler(*data.Event)

	if !allowed {
		return errorResponse(errors.New("please wait for the cooldown"))
	}

	// Command Logic
	var options struct {
		Prompt string `discord:"prompt"`
	}

	if err := data.Options.Unmarshal(&options); err != nil {
		return errorResponse(err)
	}

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	// Send the generation request to DALL·E 3
	resp, err := client.CreateImage(context.Background(), openai.ImageRequest{
		Prompt: options.Prompt,
		Model:  "dall-e-3",
		Size:   "1024x1024",
	})
	if err != nil {
		log.Printf("Image creation error: %v\n", err)
		return errorResponse(fmt.Errorf("failed to generate image"))
	}

	imageRes, err := http.Get(resp.Data[0].URL)

	if err != nil {
		return errorResponse(err)
	}

	defer imageRes.Body.Close()

	imageBytes, err := io.ReadAll(imageRes.Body)

	if err != nil {
		return errorResponse(err)
	}

	imageFile := bytes.NewBuffer(imageBytes)

	file := sendpart.File{
		Name:   "himbot_response.png",
		Reader: imageFile,
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString("Prompt: " + options.Prompt),
		Files:   []sendpart.File{file},
	}
}

func (h *handler) cmdHS(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	var options struct {
		Arg string `discord:"nickname"`
	}

	if err := data.Options.Unmarshal(&options); err != nil {
		return errorResponse(err)
	}

	user := lib.GetUserObject(*data.Event)

	return &api.InteractionResponseData{
		Content: option.NewNullableString(options.Arg + " was " + user.DisplayName() + "'s nickname in highschool!"),
	}
}

func errorResponse(err error) *api.InteractionResponseData {
	return &api.InteractionResponseData{
		Content:         option.NewNullableString("**Error:** " + err.Error()),
		Flags:           discord.EphemeralMessage,
		AllowedMentions: &api.AllowedMentions{ /* none */ },
	}
}
