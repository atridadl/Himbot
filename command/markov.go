package command

import (
	"log"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// MarkovCommand generates a random message using Markov chains
func MarkovCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Println("MarkovCommand called")

	// Get the channel ID from the interaction
	channelID := i.ChannelID
	log.Printf("Channel ID: %s", channelID)

	// Fetch the last 100 messages from the channel
	messages, err := s.ChannelMessages(channelID, 100, "", "", "")
	if err != nil {
		log.Printf("Error fetching messages: %v", err)
		respondWithError(s, i, "Failed to fetch messages")
		return
	}
	log.Printf("Fetched %d messages", len(messages))

	// Build the Markov chain from the fetched messages
	chain := buildMarkovChain(messages)
	log.Printf("Built Markov chain with %d entries", len(chain))

	// Generate a new message using the Markov chain
	newMessage := generateMessage(chain)
	log.Printf("Generated message: %s", newMessage)

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
		return
	}
	log.Println("Successfully responded to interaction")
}

// buildMarkovChain creates a Markov chain from a list of messages
func buildMarkovChain(messages []*discordgo.Message) map[string][]string {
	chain := make(map[string][]string)
	for _, msg := range messages {
		words := strings.Fields(msg.Content)
		log.Printf("Processing message: %s", msg.Content)
		// Build the chain by associating each word with the word that follows it
		for i := 0; i < len(words)-1; i++ {
			chain[words[i]] = append(chain[words[i]], words[i+1])
		}
	}
	log.Printf("Built chain with %d entries", len(chain))
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

// respondWithError sends an error message as a response to the interaction
func respondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	log.Printf("Responding with error: %s", message)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
	if err != nil {
		log.Printf("Error sending error response: %v", err)
	}
}
