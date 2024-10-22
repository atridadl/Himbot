package command

import (
	"fmt"
	"himbot/lib"
	"time"

	"github.com/bwmarrin/discordgo"
)

func PingCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !lib.CheckAndApplyCooldown(s, i, "ping", 5*time.Second) {
		return
	}

	// Customize the response based on whether it's a guild or DM
	responseContent := "Pong!"

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: responseContent,
		},
	})

	if err != nil {
		fmt.Println("Error responding to interaction:", err)
		// Optionally, you could try to send an error message to the user
		lib.RespondWithError(s, i, "An error occurred while processing the command")
	}
}
