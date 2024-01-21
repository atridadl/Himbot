package command

import (
	"context"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

func Ping(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	// Command Logic
	return &api.InteractionResponseData{
		Content: option.NewNullableString("Pong!"),
	}
}
