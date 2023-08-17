import { ApplyOptions } from '@sapphire/decorators';
import { Args, Command } from '@sapphire/framework';
import { AttachmentBuilder, Message, blockQuote, codeBlock } from 'discord.js';
import OpenAI from 'openai';

const openai = new OpenAI({
	apiKey: process.env.OPENAI_API_KEY
});

@ApplyOptions<Command.Options>({
	description: 'You can ACTUALLY ask Himbot something! So cool!',
	options: ['prompt']
})
export class UserCommand extends Command {
	// Register Chat Input and Context Menu command
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

	// Message command
	public async messageRun(message: Message, args: Args) {
		return this.ask(message, args.getOption('prompt') || message.content.split('!ask ')[1]);
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

		const chatCompletion = await openai.chat.completions.create({
			model: 'gpt-4',
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

		if (interactionOrMessage instanceof Message) {
			return askMessage.edit({
				content:
					content.length < 2000
						? content
						: `Discord only allows messages with 2000 characters or less. Please see your response in the attached file!`,
				files: messageAttachment
			});
		}

		return interactionOrMessage.editReply({
			content:
				content.length < 2000
					? content
					: `Discord only allows messages with 2000 characters or less. Please see your response in the attached file!`,
			files: messageAttachment
		});
	}
}
