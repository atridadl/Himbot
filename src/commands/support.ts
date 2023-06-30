import { ApplyOptions } from '@sapphire/decorators';
import { Command } from '@sapphire/framework';
import { Message } from 'discord.js';

@ApplyOptions<Command.Options>({
	description: "Help support HimBot's AI features."
})
export class UserCommand extends Command {
	// Register Chat Input and Context Menu command
	public override registerApplicationCommands(registry: Command.Registry) {
		registry.registerChatInputCommand((builder) => builder.setName(this.name).setDescription(this.description));
	}

	// Message command
	public async messageRun(message: Message) {
		return this.sendSupport(message);
	}

	// Chat Input (slash) command
	public async chatInputRun(interaction: Command.ChatInputCommandInteraction) {
		return this.sendSupport(interaction);
	}

	private async sendSupport(interactionOrMessage: Message | Command.ChatInputCommandInteraction | Command.ContextMenuCommandInteraction) {
		interactionOrMessage instanceof Message
			? await interactionOrMessage.channel.send({
					content: 'Thanks! The link to donate is here: https://buy.stripe.com/eVa6p95ho0Kn33a8wx'
			  })
			: await interactionOrMessage.reply({
					content: 'Thanks! The link to donate is here: https://buy.stripe.com/eVa6p95ho0Kn33a8wx',
					fetchReply: true
			  });
	}
}
