import { ApplyOptions } from '@sapphire/decorators';
import { Command } from '@sapphire/framework';
import { Message } from 'discord.js';

// @ts-ignore
@ApplyOptions<Command.Options>({
	description: 'Dad joke for daddies only!'
})
export class UserCommand extends Command {
	// Register Chat Input and Context Menu command
	public override registerApplicationCommands(registry: Command.Registry) {
		// Register Chat Input command
		registry.registerChatInputCommand({
			name: this.name,
			description: this.description
		});
	}

	// Message command
	public async messageRun(message: Message) {
		return this.sendJoke(message);
	}

	// Chat Input (slash) command
	public async chatInputRun(interaction: Command.ChatInputCommandInteraction) {
		return this.sendJoke(interaction);
	}

	private async sendJoke(interactionOrMessage: Message | Command.ChatInputCommandInteraction | Command.ContextMenuCommandInteraction) {
		const jokeResponse = await fetch('https://icanhazdadjoke.com/', {
			headers: {
				Accept: 'application/json'
			}
		});
		const jokeData = (await jokeResponse.json()) as { status: number; joke: string };

		interactionOrMessage instanceof Message
			? await interactionOrMessage.channel.send({
					content: jokeData.status === 200 ? jokeData.joke : '404 Joke Not Found'
			  })
			: await interactionOrMessage.reply({
					content: jokeData.status === 200 ? jokeData.joke : '404 Joke Not Found',
					fetchReply: true
			  });
	}
}
