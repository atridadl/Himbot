package lib

import (
	"net"
	"os"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

func ErrorResponse(err error) *api.InteractionResponseData {
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
