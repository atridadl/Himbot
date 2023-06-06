import { ApplyOptions } from '@sapphire/decorators';
import { Command } from '@sapphire/framework';
import { Message } from 'discord.js';

@ApplyOptions<Command.Options>({
	description: 'Quack!'
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
		return this.sendQuack(message);
	}

	// Chat Input (slash) command
	public async chatInputRun(interaction: Command.ChatInputCommandInteraction) {
		return this.sendQuack(interaction);
	}

	private async sendQuack(interactionOrMessage: Message | Command.ChatInputCommandInteraction | Command.ContextMenuCommandInteraction) {
		const duckResponse = await fetch('https://random-d.uk/api/v2/quack');
		const duckData = await duckResponse.json();

		interactionOrMessage instanceof Message
			? await interactionOrMessage.channel.send({
					content: duckData.url
			  })
			: await interactionOrMessage.reply({
					content: duckData.url,
					fetchReply: true
			  });
	}
}
