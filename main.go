package main

import (
	"context"
	"himbot/command"
	"log"
	"os"
	"os/signal"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/joho/godotenv"
)

var commands = []api.CreateCommandData{
	{
		Name:        "ping",
		Description: "ping pong!",
	},
	{
		Name:        "hs",
		Description: "This command was your nickname in highschool!",
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  "nickname",
				Description: "Your nickname in highschool.",
				Required:    true,
			},
		},
	},
}

func main() {

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		godotenv.Load(".env")

		if token == "" {
			log.Fatalln("No $DISCORD_TOKEN given.")
		}
	}

	h := newHandler(state.New("Bot " + token))
	h.s.AddInteractionHandler(h)
	h.s.AddIntents(gateway.IntentGuilds)
	h.s.AddHandler(func(*gateway.ReadyEvent) {
		me, _ := h.s.Me()
		log.Println("connected to the gateway as", me.Tag())
	})

	if err := cmdroute.OverwriteCommands(h.s, commands); err != nil {
		log.Fatalln("cannot update commands:", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := h.s.Connect(ctx); err != nil {
		log.Fatalln("cannot connect:", err)
	}
}

type handler struct {
	*cmdroute.Router
	s *state.State
}

func newHandler(s *state.State) *handler {
	h := &handler{s: s}

	h.Router = cmdroute.NewRouter()
	// Automatically defer handles if they're slow.
	h.Use(cmdroute.Deferrable(s, cmdroute.DeferOpts{}))
	h.AddFunc("ping", command.Ping)
	h.AddFunc("ask", command.Ask)
	h.AddFunc("pic", command.Pic)
	h.AddFunc("hs", command.HS)

	return h
}
