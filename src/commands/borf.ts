import { ApplyOptions } from '@sapphire/decorators';
import { Command } from '@sapphire/framework';
import { Message } from 'discord.js';

// @ts-ignore
@ApplyOptions<Command.Options>({
	description: 'Borf! Borf!'
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
		return this.sendBorf(message);
	}

	// Chat Input (slash) command
	public async chatInputRun(interaction: Command.ChatInputCommandInteraction) {
		return this.sendBorf(interaction);
	}

	private async sendBorf(interactionOrMessage: Message | Command.ChatInputCommandInteraction | Command.ContextMenuCommandInteraction) {
		const dogResponse = await fetch('https://dog.ceo/api/breeds/image/random');
		const dogData = (await dogResponse.json()) as { message: string; status: string };

		interactionOrMessage instanceof Message
			? await interactionOrMessage.channel.send({
					content: dogData?.status === 'success' ? dogData.message : 'Error: I had troubles fetching perfect puppies for you... :('
			  })
			: await interactionOrMessage.reply({
					content: dogData.status === 'success' ? dogData.message : 'Error: I had troubles fetching perfect puppies for you... :(',
					fetchReply: true
			  });
	}
}
