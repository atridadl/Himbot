package main

import (
	"himbot/command"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "ping pong!",
		},
		{
			Name:        "hs",
			Description: "This command was your nickname in highschool!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "nickname",
					Description: "Your nickname in highschool.",
					Required:    true,
				},
			},
		},
		{
			Name:        "markov",
			Description: "Why did the Markov chain break up with its partner? Because it couldn't handle the past!",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping":   command.PingCommand,
		"hs":     command.HsCommand,
		"markov": command.MarkovCommand,
	}
)

func main() {
	godotenv.Load(".env")

	token := os.Getenv("DISCORD_TOKEN")

	if token == "" {
		log.Fatalln("No $DISCORD_TOKEN given.")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	dg.AddHandler(ready)
	dg.AddHandler(interactionCreate)

	dg.Identify.Intents = discordgo.IntentsGuilds

	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}

	log.Println("Bot is now running. Press CTRL-C to exit.")
	registerCommands(dg)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
}

func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
		h(s, i)
	}
}

func registerCommands(s *discordgo.Session) {
	log.Println("Registering commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
}
