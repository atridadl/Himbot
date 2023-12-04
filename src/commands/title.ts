import { Command, container } from '@sapphire/framework';

export class TitleCommand extends Command {
	public constructor(context: Command.LoaderContext) {
		super(context, {
			description: 'This command is the title of your sextape.',
			options: ['title']
		});
	}

	// Register Chat Input and Context Menu command
	public override registerApplicationCommands(registry: Command.Registry) {
		registry.registerChatInputCommand((builder) =>
			builder //
				.setName(this.name)
				.setDescription(this.description)
				.addStringOption((option) => option.setName('title').setDescription('The title of your sextape.').setRequired(true))
		);
	}

	// Chat Input (slash) command
	public async chatInputRun(interaction: Command.ChatInputCommandInteraction) {
		const title = interaction.options.getString('title') || 'NOTHING';

		await interaction.reply({
			content: `${title}: Title of ${interaction.user.username}'s sex tape!`,
			fetchReply: true
		});
	}
}

void container.stores.loadPiece({
	store: 'commands',
	name: 'title',
	piece: TitleCommand
});
