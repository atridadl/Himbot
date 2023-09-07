import { ApplyOptions } from '@sapphire/decorators';
import { Command } from '@sapphire/framework';
import { Message } from 'discord.js';

// @ts-ignore
@ApplyOptions<Command.Options>({
	description: 'Disclosures for privacy.'
})
export class UserCommand extends Command {
	// Register Chat Input and Context Menu command
	public override registerApplicationCommands(registry: Command.Registry) {
		registry.registerChatInputCommand((builder) => builder.setName(this.name).setDescription(this.description));
	}

	// Message command
	public async messageRun(message: Message) {
		return this.sendDisclosure(message);
	}

	// Chat Input (slash) command
	public async chatInputRun(interaction: Command.ChatInputCommandInteraction) {
		return this.sendDisclosure(interaction);
	}

	private async sendDisclosure(interactionOrMessage: Message | Command.ChatInputCommandInteraction | Command.ContextMenuCommandInteraction) {
		interactionOrMessage instanceof Message
			? await interactionOrMessage.channel.send({
					content: `Commands such as /ask and /pic do not send your discord data to the AI providers. They will, however, send anything you include in the prompts. /ask uses OpenAI's gpt3.5 model and /pic uses Stable Diffusion. All other commands and the data associated stays with HimBot's server and Discord's Servers.`
			  })
			: await interactionOrMessage.reply({
					content: `Commands such as /ask and /pic do not send your discord data to the AI providers. They will, however, send anything you include in the prompts. /ask uses OpenAI's gpt3.5 model and /pic uses Stable Diffusion. All other commands and the data associated stays with HimBot's server and Discord's Servers.`,
					fetchReply: true
			  });
	}
}
