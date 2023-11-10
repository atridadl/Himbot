import { ApplyOptions } from '@sapphire/decorators';
import { Args, BucketScope, Command } from '@sapphire/framework';
import { AttachmentBuilder, Message } from 'discord.js';
import OpenAI from 'openai';

const openai = new OpenAI({
	apiKey: process.env.OPENAI_API_KEY
});

// @ts-ignore
@ApplyOptions<Command.Options>({
	description: 'Shhhh!',
	options: ['prompt'],
	// 10mins
	cooldownDelay: 400_000,
	cooldownLimit: 1,
	// Yes... I did hardcode myself.
	cooldownFilteredUsers: ['83679718401904640'],
	cooldownScope: BucketScope.User
})
export class UserCommand extends Command {
	// Register Chat Input and Context Menu command
	public override registerApplicationCommands(registry: Command.Registry) {
		registry.registerChatInputCommand((builder) =>
			builder
				.setName(this.name)
				.setDescription(this.description)
				.addStringOption((option) =>
					option.setName('prompt').setDescription('The prompt you will use to generate an image!').setRequired(true)
				)
		);
	}

	// Message command
	public async messageRun(message: Message, args: Args) {
		return this.pic(message, args.getOption('prompt') || 'Scold me for not passing any prompt in.');
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

		const response = await openai.images.generate({
			model: 'dall-e-3',
			prompt,
			n: 1,
			size: '1024x1024'
		});

		const imageUrl = response.data[0].url || '';
		// get an array buffer
		const imageBuffer = await fetch(imageUrl).then((r) => r.arrayBuffer());

		if (response.data.length === 0) {
			const content = `Sorry, I can't complete the prompt for: ${prompt}`;

			if (interactionOrMessage instanceof Message) {
				return askMessage.edit({ content });
			}

			return interactionOrMessage.editReply({
				content: content
			});
		} else {
			const imageAttachment: AttachmentBuilder[] = [];

			imageAttachment.push(
				new AttachmentBuilder(Buffer.from(new Uint8Array(imageBuffer)), {
					name: 'response.jpg',
					description: "Himbot's Response"
				})
			);

			const content = `Prompt: ${prompt}:`;

			if (interactionOrMessage instanceof Message) {
				return askMessage.edit({ content, files: imageAttachment });
			}

			return interactionOrMessage.editReply({
				content,
				files: imageAttachment
			});
		}
	}
}
