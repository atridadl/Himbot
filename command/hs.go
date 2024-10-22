package command

import (
	"fmt"
	"himbot/lib"
	"time"

	"github.com/bwmarrin/discordgo"
)

func HsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
    if !lib.CheckAndApplyCooldown(s, i, "hs", 10*time.Second) {
        return
    }

    options := i.ApplicationCommandData().Options
    if len(options) == 0 || options[0].Type != discordgo.ApplicationCommandOptionString {
        lib.RespondWithError(s, i, "Please provide a nickname.")
        return
    }
    nickname := options[0].StringValue()

    user, err := lib.GetUser(i)
    if err != nil {
        lib.RespondWithError(s, i, "Error processing command: "+err.Error())
        return
    }

    response := fmt.Sprintf("%s was %s's nickname in high school!", nickname, user.Username)

    err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Content: response,
        },
    })

    if err != nil {
        fmt.Println("Error responding to interaction:", err)
        lib.RespondWithError(s, i, "An error occurred while processing the command")
    }
}
