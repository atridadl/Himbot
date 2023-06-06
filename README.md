# HimBot

A discord bot written in Typescript.

## It's dangerous to go alone! Take this!

-   Install Nodejs 18 or later [here](https://nodejs.org/en/download) (required)
-   The pnpm package manager `npm i -g pnpm` (recommended)

## Structure

Commands and Listeners are all stored in named files within the src/commands and src/listeners directories respectively.

## Running Locally

-   Copy .env.example and rename to .env
-   Create a Discord Bot with all gateway permissions enabled
-   Generate a token for this discord bot and paste it in the .env for DISCORD_TOKEN
-   Run `pnpm dev` to run locally

## Adding the bot to a server

Use the following link (replacing DISCORD_CLIENT_ID with your own bot's client ID of course...) to add your bot:
https://discord.com/oauth2/authorize?client_id=DISCORD_CLIENT_ID&scope=bot&permissions=8
