package command

import (
	"context"
	"errors"
	"himbot/lib"
	"strconv"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
)

func Pic(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	// Cooldown Logic
	allowed, cooldownString := lib.CooldownHandler(*data.Event, "pic", time.Minute*5)

	if !allowed {
		return lib.ErrorResponse(errors.New(cooldownString))
	}

	// Command Logic
	var options struct {
		Prompt string `discord:"prompt"`
	}

	if err := data.Options.Unmarshal(&options); err != nil {
		lib.CancelCooldown(data.Event.Member.User.ID.String(), "pic")
		return lib.ErrorResponse(err)
	}

	// Get current epoch timestamp
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// Concatenate clean username and timestamp to form filename
	filename := data.Event.Sender().Username + "_" + timestamp + ".jpg"

	imageFile, err := lib.OpenAIImageGeneration(options.Prompt, filename)

	if err != nil {
		lib.CancelCooldown(data.Event.Member.User.ID.String(), "pic")
		return lib.ErrorResponse(err)
	}

	file := sendpart.File{
		Name:   filename,
		Reader: imageFile,
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString("Prompt: " + options.Prompt),
		Files:   []sendpart.File{file},
	}
}
