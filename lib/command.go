package lib

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type CommandFunc func(s *discordgo.Session, i *discordgo.InteractionCreate) (string, error)

func HandleCommand(commandName string, cooldownDuration time.Duration, handler CommandFunc) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if !CheckAndApplyCooldown(s, i, commandName, cooldownDuration) {
			return
		}

		// Acknowledge the interaction immediately
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
		if err != nil {
			ThrowWithError(commandName, "Error deferring response: "+err.Error())
			return
		}

		// Execute the command handler
		response, err := handler(s, i)

		if err != nil {
			RespondWithError(s, i, "Error processing command: "+err.Error())
			return
		}

		// Send the follow-up message with the response
		_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: response,
		})

		if err != nil {
			ThrowWithError(commandName, "Error sending follow-up message: "+err.Error())
		}
	}
}
