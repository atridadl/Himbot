# HimBot

A discord bot written in Typescript and running on the Deno runtime.

## It's dangerous to go alone! Take this!

- Install Deno [here](https://deno.com/manual@v1.34.1/getting_started/installation) (required)
- The [Deno VSCode Extension](https://marketplace.visualstudio.com/items?itemName=denoland.vscode-deno) (recommended)

## Structure

Commands and Events are all stored in named files within the src/commands and src/events directories respectively.
Usage and example ts files can be found in the examples folder.

To generate a new Command or Event run `deno task new:command` or `deno task new:event` respectively.

## Running Locally

- Copy .env.example and rename to .env
- Create a Discord Bot with all gateway permissions enabled
- Generate a token for this discord bot and paste it in the .env for BOT_TOKEN
- Run `deno run --allow-all mod.ts` to run locally

## Adding the bot to a server

Use the following link (replacing BOT_TOKEN with your own Token of course...) to add your bot:
https://discord.com/oauth2/authorize?client_id=BOT_TOKEN&scope=bot&permissions=8
