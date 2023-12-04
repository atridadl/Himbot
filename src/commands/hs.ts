import { Command, container } from '@sapphire/framework';

export class HighSchoolCommand extends Command {
	public constructor(context: Command.LoaderContext) {
		super(context, {
			description: 'This command was your nickname in highschool!',
			options: ['nickname']
		});
	}

	// Register Chat Input and Context Menu command
	public override registerApplicationCommands(registry: Command.Registry) {
		registry.registerChatInputCommand((builder) =>
			builder //
				.setName(this.name)
				.setDescription(this.description)
				.addStringOption((option) => option.setName('nickname').setDescription('Your nickname in highschool.').setRequired(true))
		);
	}

	// Chat Input (slash) command
	public async chatInputRun(interaction: Command.ChatInputCommandInteraction) {
		const nickname = interaction.options.getString('nickname') || 'NOTHING';
		await interaction.reply({
			content: `${nickname} was ${interaction.user.username}'s nickname in highschool!`,
			fetchReply: true
		});
	}
}

void container.stores.loadPiece({
	store: 'commands',
	name: 'hs',
	piece: HighSchoolCommand
});
