package lib

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// respondWithError sends an error message as a response to the interaction
func RespondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	log.Printf("Responding with error: %s", message)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Error sending error response: %v", err)
	}
}

func ThrowWithError(command, message string) error {
	return fmt.Errorf("error in command '%s': %s", command, message)
}
