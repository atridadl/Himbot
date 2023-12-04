import { Command, container } from '@sapphire/framework';
import { AttachmentBuilder, blockQuote, codeBlock } from 'discord.js';
import OpenAI from 'openai';

const openai = new OpenAI({
	apiKey: process.env.OPENAI_API_KEY
});

export class AskCommand extends Command {
	public constructor(context: Command.LoaderContext) {
		super(context, {
			description: 'You can ACTUALLY ask Himbot something! So cool!',
			options: ['prompt']
		});
	}

	public override registerApplicationCommands(registry: Command.Registry) {
		registry.registerChatInputCommand((builder) =>
			builder
				.setName(this.name)
				.setDescription(this.description)
				.addStringOption((option) =>
					option.setName('prompt').setDescription('You can ACTUALLY ask Himbot something! So cool!').setRequired(true)
				)
		);
	}

	public async chatInputRun(interaction: Command.ChatInputCommandInteraction) {
		const prompt = interaction.options.getString('prompt');

		await interaction.reply({ content: '🤔 Thinking... 🤔', fetchReply: true });

		const chatCompletion = await openai.chat.completions.create({
			model: 'gpt-4-1106-preview',
			messages: [
				{
					role: 'user',
					content: prompt
				}
			]
		});

		const content = blockQuote(`> ${prompt}\n${codeBlock(`${chatCompletion.choices[0].message?.content}`)}`);

		const messageAttachment: AttachmentBuilder[] = [];

		if (content.length > 2000) {
			messageAttachment.push(
				new AttachmentBuilder(Buffer.from(`> ${prompt}\n${`${chatCompletion.choices[0].message?.content}`}`, 'utf-8'), {
					name: 'response.txt',
					description: "Himbot's Response"
				})
			);
		}

		return interaction.editReply({
			content:
				content.length < 2000
					? content
					: `Discord only allows messages with 2000 characters or less. Please see your response in the attached file!`,
			files: messageAttachment
		});
	}
}

void container.stores.loadPiece({
	store: 'commands',
	name: 'ask',
	piece: AskCommand
});
