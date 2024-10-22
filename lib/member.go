package lib

import (
	"github.com/bwmarrin/discordgo"
)

// InteractionUser represents a user from an interaction, abstracting away the differences
// between guild members and DM users.
type InteractionUser struct {
	ID       string
	Username string
	Bot      bool
}

// GetUser extracts user information from an interaction, handling both guild and DM cases.
func GetUser(i *discordgo.InteractionCreate) (*InteractionUser, error) {
	if i.Member != nil && i.Member.User != nil {
		// Guild interaction
		return &InteractionUser{
			ID:       i.Member.User.ID,
			Username: i.Member.User.Username,
			Bot:      i.Member.User.Bot,
		}, nil
	} else if i.User != nil {
		// DM interaction
		return &InteractionUser{
			ID:       i.User.ID,
			Username: i.User.Username,
			Bot:      i.User.Bot,
		}, nil
	}

	return nil, ThrowWithError("GetUser", "Unable to extract user information from interaction")
}

// IsInGuild checks if the interaction occurred in a guild.
func IsInGuild(i *discordgo.InteractionCreate) bool {
	return i.Member != nil
}

// GetGuildID safely retrieves the guild ID if the interaction is from a guild.
func GetGuildID(i *discordgo.InteractionCreate) string {
	if i.GuildID != "" {
		return i.GuildID
	}
	return ""
}
