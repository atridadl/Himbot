package command

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"himbot/lib"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
)

func Ask(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	// Cooldown Logic
	allowed, cooldownString := lib.CooldownHandler(*data.Event, "ask", time.Minute)

	if !allowed {
		return lib.ErrorResponse(errors.New(cooldownString))
	}

	// Command Logic
	var options struct {
		Prompt string `discord:"prompt"`
	}

	if err := data.Options.Unmarshal(&options); err != nil {
		lib.CancelCooldown(data.Event.User.ID.String(), "ask")
		return lib.ErrorResponse(err)
	}

	respString, err := lib.OpenAITextGeneration(options.Prompt)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		lib.CancelCooldown(data.Event.User.ID.String(), "ask")
		return &api.InteractionResponseData{
			Content:         option.NewNullableString("ChatCompletion Error!"),
			AllowedMentions: &api.AllowedMentions{},
		}
	}

	if len(respString) > 1800 {
		textFile := bytes.NewBuffer([]byte(respString))

		file := sendpart.File{
			Name:   "himbot_response.md",
			Reader: textFile,
		}

		return &api.InteractionResponseData{
			Content:         option.NewNullableString("Prompt: " + options.Prompt + "\n"),
			AllowedMentions: &api.AllowedMentions{},
			Files:           []sendpart.File{file},
		}
	}
	return &api.InteractionResponseData{
		Content:         option.NewNullableString("Prompt: " + options.Prompt + "\n--------------------\n" + respString),
		AllowedMentions: &api.AllowedMentions{},
	}
}
