# HimBot

A discord bot written in Go.

## It's dangerous to go alone! Take this!

-   Install Go 1.21.5 or higher (required)

## Running Locally

-   Copy .env.example and rename to .env
-   Create a Discord Bot with all gateway permissions enabled
-   Generate a token for this discord bot and paste it in the .env for DISCORD_TOKEN
-   Generate and provide an Replicate token and paste it in the .env for REPLICATE_API_TOKEN
-   Run `go run main.go` to run locally

## Adding the bot to a server

Use the following link (replacing DISCORD_CLIENT_ID with your own bot's client ID of course...) to add your bot:
https://discord.com/oauth2/authorize?client_id=DISCORD_CLIENT_ID&scope=bot&permissions=8
