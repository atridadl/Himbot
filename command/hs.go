package command

import (
	"context"
	"himbot/lib"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

func HS(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	var options struct {
		Arg string `discord:"nickname"`
	}

	if err := data.Options.Unmarshal(&options); err != nil {
		return lib.ErrorResponse(err)
	}

	user := lib.GetUserObject(*data.Event)

	return &api.InteractionResponseData{
		Content: option.NewNullableString(options.Arg + " was " + user.DisplayName() + "'s nickname in highschool!"),
	}
}
