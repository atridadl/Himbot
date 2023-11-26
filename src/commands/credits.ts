import { ApplyOptions } from '@sapphire/decorators';
import { Command } from '@sapphire/framework';
import { Message } from 'discord.js';

// @ts-ignore
@ApplyOptions<Command.Options>({
	description: 'Check how many credits I have left!'
})
export class UserCommand extends Command {
	// Register Chat Input and Context Menu command
	public override registerApplicationCommands(registry: Command.Registry) {
		registry.registerChatInputCommand((builder) => builder.setName(this.name).setDescription(this.description));
	}

	// Message command
	public async messageRun(message: Message) {
		return this.credits(message);
	}

	// Chat Input (slash) command
	public async chatInputRun(interaction: Command.ChatInputCommandInteraction) {
		return this.credits(interaction);
	}

	private async credits(interactionOrMessage: Message | Command.ChatInputCommandInteraction | Command.ContextMenuCommandInteraction) {
		const askMessage =
			interactionOrMessage instanceof Message
				? await interactionOrMessage.channel.send({ content: '🤔 Thinking... 🤔' })
				: await interactionOrMessage.reply({ content: '🤔 Thinking... 🤔', fetchReply: true });

		const creditCountResponse = await fetch(`https://api.stability.ai/v1/user/balance`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${process.env.STABILITY_API_KEY}`
			}
		});

		const balance = ((await creditCountResponse.json()) as { credits: number }).credits || 0;

		const content = `I have ${balance} credits remaining for image generation!`;

		if (interactionOrMessage instanceof Message) {
			return askMessage.edit({ content });
		}

		return interactionOrMessage.editReply({
			content: content
		});
	}
}
