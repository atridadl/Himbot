package command

import (
	"github.com/bwmarrin/discordgo"
)

func PingCommand(s *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	// Customize the response based on whether it's a guild or DM
	responseContent := "Pong!"

	return responseContent, nil
}
