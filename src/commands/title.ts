import { ApplyOptions } from '@sapphire/decorators';
import { Args, Command } from '@sapphire/framework';
import { Message } from 'discord.js';

@ApplyOptions<Command.Options>({
	description: 'This command is the title of your sextape.',
	options: ['title']
})
export class UserCommand extends Command {
	// Register Chat Input and Context Menu command
	public override registerApplicationCommands(registry: Command.Registry) {
		registry.registerChatInputCommand((builder) =>
			builder //
				.setName(this.name)
				.setDescription(this.description)
				.addStringOption((option) => option.setName('title').setDescription('The title of your sextape.').setRequired(true))
		);
	}

	// Message command
	public async messageRun(message: Message, args: Args) {
		return this.titleHandler(message, args.getOption('title') || message.content.split('!title ')[1]);
	}

	// Chat Input (slash) command
	public async chatInputRun(interaction: Command.ChatInputCommandInteraction) {
		return this.titleHandler(interaction, interaction.options.getString('title') || 'NOTHING');
	}

	private async titleHandler(
		interactionOrMessage: Message | Command.ChatInputCommandInteraction | Command.ContextMenuCommandInteraction,
		title: string
	) {
		interactionOrMessage instanceof Message
			? await interactionOrMessage.channel.send({
					content: `${title}: Title of ${interactionOrMessage.author.username}'s sex tape!`
			  })
			: await interactionOrMessage.reply({
					content: `${title}: Title of ${interactionOrMessage.user.username}'s sex tape!`,
					fetchReply: true
			  });
	}
}
