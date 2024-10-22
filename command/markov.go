package command

import (
	"himbot/lib"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func MarkovCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !lib.CheckAndApplyCooldown(s, i, "markov", 30*time.Second) {
		return
	}

	// Get the channel ID from the interaction
	channelID := i.ChannelID

	// Get the number of messages to fetch from the option
	numMessages := 100 // Default value
	if len(i.ApplicationCommandData().Options) > 0 {
		if i.ApplicationCommandData().Options[0].Name == "messages" {
			numMessages = int(i.ApplicationCommandData().Options[0].IntValue())
			if numMessages <= 0 {
				numMessages = 100
			} else if numMessages > 1000 {
				numMessages = 1000 // Limit to 1000 messages max
			}
		}
	}

	// Fetch messages
	allMessages, err := fetchMessages(s, channelID, numMessages)
	if err != nil {
		lib.RespondWithError(s, i, "Failed to fetch messages: "+err.Error())
		return
	}

	// Build the Markov chain from the fetched messages
	chain := buildMarkovChain(allMessages)

	// Generate a new message using the Markov chain
	newMessage := generateMessage(chain)

	// Check if the generated message is empty and provide a fallback message
	if newMessage == "" {
		newMessage = "I couldn't generate a message. The channel might be empty or contain no usable text."
	}

	// Respond to the interaction with the generated message
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: newMessage,
		},
	})

	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
		lib.RespondWithError(s, i, "An error occurred while processing the command")
	}
}

func fetchMessages(s *discordgo.Session, channelID string, numMessages int) ([]*discordgo.Message, error) {
	var allMessages []*discordgo.Message
	var lastMessageID string

	for len(allMessages) < numMessages {
		batchSize := 100
		if numMessages-len(allMessages) < 100 {
			batchSize = numMessages - len(allMessages)
		}

		batch, err := s.ChannelMessages(channelID, batchSize, lastMessageID, "", "")
		if err != nil {
			return nil, err
		}

		if len(batch) == 0 {
			break // No more messages to fetch
		}

		allMessages = append(allMessages, batch...)
		lastMessageID = batch[len(batch)-1].ID

		if len(batch) < 100 {
			break // Less than 100 messages returned, we've reached the end
		}
	}

	return allMessages, nil
}

// buildMarkovChain creates a Markov chain from a list of messages
func buildMarkovChain(messages []*discordgo.Message) map[string][]string {
	chain := make(map[string][]string)
	for _, msg := range messages {
		words := strings.Fields(msg.Content)
		// Build the chain by associating each word with the word that follows it
		for i := 0; i < len(words)-1; i++ {
			chain[words[i]] = append(chain[words[i]], words[i+1])
		}
	}
	return chain
}

// generateMessage creates a new message using the Markov chain
func generateMessage(chain map[string][]string) string {
	if len(chain) == 0 {
		return ""
	}

	words := []string{}
	var currentWord string

	// Start with a random word from the chain
	for word := range chain {
		currentWord = word
		break
	}

	// Generate up to 20 words
	for i := 0; i < 20; i++ {
		words = append(words, currentWord)
		if nextWords, ok := chain[currentWord]; ok && len(nextWords) > 0 {
			// Randomly select the next word from the possible follow-ups
			currentWord = nextWords[rand.Intn(len(nextWords))]
		} else {
			break
		}
	}

	return strings.Join(words, " ")
}
