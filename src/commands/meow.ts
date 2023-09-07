import { ApplyOptions } from '@sapphire/decorators';
import { Command } from '@sapphire/framework';
import { Message } from 'discord.js';

// @ts-ignore
@ApplyOptions<Command.Options>({
	description: 'Meow!'
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
		return this.sendMeow(message);
	}

	// Chat Input (slash) command
	public async chatInputRun(interaction: Command.ChatInputCommandInteraction) {
		return this.sendMeow(interaction);
	}

	private async sendMeow(interactionOrMessage: Message | Command.ChatInputCommandInteraction | Command.ContextMenuCommandInteraction) {
		const catResponse = await fetch('https://cataas.com/cat?json=true');
		const catData = await catResponse.json();
		interactionOrMessage instanceof Message
			? await interactionOrMessage.channel.send({
					content: `https://cataas.com/${catData.url}`
			  })
			: await interactionOrMessage.reply({
					content: `https://cataas.com/${catData.url}`,
					fetchReply: true
			  });
	}
}
