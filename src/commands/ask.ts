import { ApplyOptions } from '@sapphire/decorators';
import { Args, Command } from '@sapphire/framework';
import { Message, blockQuote, codeBlock } from 'discord.js';
import { Configuration, OpenAIApi } from 'openai';

const configuration = new Configuration({
	apiKey: process.env.OPENAI_API_KEY
});
const openai = new OpenAIApi(configuration);

@ApplyOptions<Command.Options>({
	description: 'You can ACTUALLY ask Himbot something! So cool!',
	options: ['prompt']
})
export class UserCommand extends Command {
	// Register Chat Input and Context Menu command
	public override registerApplicationCommands(registry: Command.Registry) {
		registry.registerChatInputCommand((builder) =>
			builder //
				.setName(this.name)
				.setDescription(this.description)
				.addStringOption((option) =>
					option.setName('prompt').setDescription('You can ACTUALLY ask Himbot something! So cool!').setRequired(true)
				)
		);
	}

	// Message command
	public async messageRun(message: Message, args: Args) {
		return this.ask(message, args.getOption('prompt') || message.content.split('!wryna ')[1]);
	}

	// Chat Input (slash) command
	public async chatInputRun(interaction: Command.ChatInputCommandInteraction) {
		return this.ask(interaction, interaction.options.getString('prompt') || 'NOTHING');
	}

	private async ask(interactionOrMessage: Message | Command.ChatInputCommandInteraction | Command.ContextMenuCommandInteraction, prompt: string) {
		const askMessage =
			interactionOrMessage instanceof Message
				? await interactionOrMessage.channel.send({ content: '🤔 Thinking... 🤔' })
				: await interactionOrMessage.reply({ content: '🤔 Thinking... 🤔', fetchReply: true });

		const chatCompletion = await openai.createChatCompletion({
			model: 'gpt-3.5-turbo',
			messages: [
				{
					role: 'user',
					content: `${prompt}. In your response, please replace "As an AI language model" with "As HimBot, the coolest bot ever"`
				}
			],
			max_tokens: 420
		});

		const content = blockQuote(`> ${prompt}\n${codeBlock(`${chatCompletion.data.choices[0].message?.content}`)}`);

		if (interactionOrMessage instanceof Message) {
			return askMessage.edit({ content: content.length <= 2000 ? content : 'Sorry... AI no work good...' });
		}

		return interactionOrMessage.editReply({
			content: content.length <= 2000 ? content : 'Sorry... AI no work good...'
		});
	}
}
