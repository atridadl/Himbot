package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"himbot/lib"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/joho/godotenv"
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
	allowed := lib.CooldownHandler(*data.Event, "ask", time.Minute)

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

	respString, err := lib.ReplicateTextGeneration(options.Prompt)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return &api.InteractionResponseData{
			Content:         option.NewNullableString("ChatCompletion Error!"),
			AllowedMentions: &api.AllowedMentions{},
		}
	}

	if len(respString) > 1800 {
		textFile := bytes.NewBuffer([]byte(respString))

		file := sendpart.File{
			Name:   "himbot_response.txt",
			Reader: textFile,
		}

		return &api.InteractionResponseData{
			Content:         option.NewNullableString("Prompt: " + options.Prompt + "\n" + "Response:\n"),
			AllowedMentions: &api.AllowedMentions{},
			Files:           []sendpart.File{file},
		}
	}
	return &api.InteractionResponseData{
		Content:         option.NewNullableString("Prompt: " + options.Prompt + "\n" + "Response: " + respString),
		AllowedMentions: &api.AllowedMentions{},
	}
}

func (h *handler) cmdPic(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	// Cooldown Logic
	allowed := lib.CooldownHandler(*data.Event, "pic", time.Minute)

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

	imageFile, err := lib.ReplicateImageGeneration(options.Prompt)

	if err != nil {
		return errorResponse(err)
	}

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
	allowed := lib.CooldownHandler(*data.Event, "hdPic", time.Minute*10)

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

	imageFile, err := lib.OpenAIImageGeneration(options.Prompt)

	if err != nil {
		return errorResponse(err)
	}

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
	var content string
	switch e := err.(type) {
	case *net.OpError:
		content = "**Network Error:** " + e.Error()
	case *os.PathError:
		content = "**File Error:** " + e.Error()
	default:
		content = "**Error:** " + err.Error()
	}

	return &api.InteractionResponseData{
		Content:         option.NewNullableString(content),
		Flags:           discord.EphemeralMessage,
		AllowedMentions: &api.AllowedMentions{},
	}
}
