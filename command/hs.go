package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func HsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	nickname := options[0].StringValue()

	var username string
	if i.Member != nil {
		username = i.Member.User.Username
	} else if i.User != nil {
		username = i.User.Username
	} else {
		username = "User"
	}

	response := fmt.Sprintf("%s was %s's nickname in highschool!", nickname, username)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}
