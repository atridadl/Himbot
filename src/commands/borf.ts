import { Command, container } from '@sapphire/framework';

export class BorfCommand extends Command {
	public constructor(context: Command.LoaderContext) {
		super(context, {
			description: 'Borf! Borf!'
		});
	}

	// Register Chat Input and Context Menu command
	public override registerApplicationCommands(registry: Command.Registry) {
		// Register Chat Input command
		registry.registerChatInputCommand({
			name: this.name,
			description: this.description
		});
	}

	// Chat Input (slash) command
	public async chatInputRun(interaction: Command.ChatInputCommandInteraction) {
		const dogResponse = await fetch('https://dog.ceo/api/breeds/image/random');
		const dogData = (await dogResponse.json()) as { message: string; status: string };

		await interaction.reply({
			content: dogData.status === 'success' ? dogData.message : 'Error: I had troubles fetching perfect puppies for you... :(',
			fetchReply: true
		});
	}
}

void container.stores.loadPiece({
	store: 'commands',
	name: 'borf',
	piece: BorfCommand
});
