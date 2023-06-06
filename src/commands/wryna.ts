import { ApplyOptions } from '@sapphire/decorators';
import { Args, Command } from '@sapphire/framework';
import { Message } from 'discord.js';

@ApplyOptions<Command.Options>({
	description: 'This command was your nickname in highschool!',
	options: ['nickname']
})
export class UserCommand extends Command {
	// Register Chat Input and Context Menu command
	public override registerApplicationCommands(registry: Command.Registry) {
		registry.registerChatInputCommand((builder) =>
			builder //
				.setName(this.name)
				.setDescription(this.description)
				.addStringOption((option) => option.setName('nickname').setDescription('Your nickname in highschool.').setRequired(true))
		);
	}

	// Message command
	public async messageRun(message: Message, args: Args) {
		return this.sendPing(message, args.getOption('nickname') || message.content.split('!wryna ')[1]);
	}

	// Chat Input (slash) command
	public async chatInputRun(interaction: Command.ChatInputCommandInteraction) {
		return this.sendPing(interaction, interaction.options.getString('nickname') || 'NOTHING');
	}

	private async sendPing(
		interactionOrMessage: Message | Command.ChatInputCommandInteraction | Command.ContextMenuCommandInteraction,
		nickname: string
	) {
		interactionOrMessage instanceof Message
			? await interactionOrMessage.channel.send({
					content: `${nickname} was ${interactionOrMessage.author.username}'s nickname in highschool!`
			  })
			: await interactionOrMessage.reply({
					content: `${nickname} was ${interactionOrMessage.user.username}'s nickname in highschool!`,
					fetchReply: true
			  });
	}
}
