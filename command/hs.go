package command

import (
	"fmt"
	"himbot/lib"

	"github.com/bwmarrin/discordgo"
)

func HsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 || options[0].Type != discordgo.ApplicationCommandOptionString {
		return "", fmt.Errorf("please provide a nickname")
	}
	nickname := options[0].StringValue()

	user, err := lib.GetUser(i)
	if err != nil {
		return "", fmt.Errorf("error processing command: %w", err)
	}

	response := fmt.Sprintf("%s was %s's nickname in high school!", nickname, user.Username)

	return response, nil
}
