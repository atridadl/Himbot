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
-   Generate and provide an OpenAPI token and paste it in the .env for OPENAI_API_KEY
-   Run `pnpm dev` to run locally

## Adding the bot to a server

Use the following link (replacing DISCORD_CLIENT_ID with your own bot's client ID of course...) to add your bot:
https://discord.com/oauth2/authorize?client_id=DISCORD_CLIENT_ID&scope=bot&permissions=8

## Commands

### **ask**

##### Description

A command that returns the answer to your prompt for OpenAI's GPT 4.0 turbo model.

#### Usage

`/ask prompt:prompt_text`

### **borf**

#### Description

A command that returns a random picture of a dog.

#### Usage

`/borf`

### **dad**

#### Description

A command that returns a random dad joke.

#### Usage

`/dad`

### **fancypic**

#### Description

A command that return an image using DALL.E with "hd" quality.

#### Usage

`/fancypic prompt`

### **pic**

#### Description

A command that return an image using DALL.E with "standard" quality.

#### Usage

`/pic prompt`

### **ping**

#### Description

A command that returns the latency for the bot, and the latency for the Discord API. Useful for bot debugging.

#### Usage

`/ping`

### **quack**

#### Description

A command that returns a random picture of a duck.

#### Usage

`/quack`

### **support**

#### Description

A command that returns a Stripe link that can be used to fund credits for the AI commands.

#### Usage

`/support`

### **hs**

#### Description

A command that returns "my nickname in highschool" response given your prompt

#### Usage

`/hs nickname:butts`
