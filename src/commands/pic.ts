import { ApplyOptions } from '@sapphire/decorators';
import { Args, Command } from '@sapphire/framework';
import { Message } from 'discord.js';
import { Configuration, OpenAIApi } from 'openai';

const configuration = new Configuration({
	apiKey: process.env.OPENAI_API_KEY
});
const openai = new OpenAIApi(configuration);

@ApplyOptions<Command.Options>({
	description: 'Make a picture!',
	options: ['prompt'],
	cooldownDelay: 20_000,
	cooldownLimit: 1
})
export class UserCommand extends Command {
	// Register Chat Input and Context Menu command
	public override registerApplicationCommands(registry: Command.Registry) {
		registry.registerChatInputCommand((builder) =>
			builder //
				.setName(this.name)
				.setDescription(this.description)
				.addStringOption((option) => option.setName('prompt').setDescription('Make a picture!').setRequired(true))
		);
	}

	// Message command
	public async messageRun(message: Message, args: Args) {
		return this.pic(message, args.getOption('prompt') || message.content.split('!wryna ')[1]);
	}

	// Chat Input (slash) command
	public async chatInputRun(interaction: Command.ChatInputCommandInteraction) {
		return this.pic(interaction, interaction.options.getString('prompt') || 'NOTHING');
	}

	private async pic(interactionOrMessage: Message | Command.ChatInputCommandInteraction | Command.ContextMenuCommandInteraction, prompt: string) {
		const askMessage =
			interactionOrMessage instanceof Message
				? await interactionOrMessage.channel.send({ content: '🤔 Thinking... 🤔' })
				: await interactionOrMessage.reply({ content: '🤔 Thinking... 🤔', fetchReply: true });

		const imageResponse = await openai.createImage({
			prompt,
			n: 1,
			size: '256x256'
		});

		const content = imageResponse.data.data[0].url || 'ERROR!';

		if (interactionOrMessage instanceof Message) {
			return askMessage.edit({ content });
		}

		return interactionOrMessage.editReply({
			content: content
		});
	}
}
