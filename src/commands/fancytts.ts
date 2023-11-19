import { ApplyOptions } from '@sapphire/decorators';
import { Args, BucketScope, Command } from '@sapphire/framework';
import { AttachmentBuilder, Message, MessageFlags } from 'discord.js';
import OpenAI from 'openai';

const openai = new OpenAI({
	apiKey: process.env.OPENAI_API_KEY
});

// @ts-ignore
@ApplyOptions<Command.Options>({
	description: 'Generate "HD" TTS every 2 minutes!',
	options: ['prompt'],
	cooldownDelay: 200_000,
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
				.addStringOption((option) => option.setName('prompt').setDescription('The prompt you will use to generate audio!').setRequired(true))
		);
	}

	// Message command
	public async messageRun(message: Message, args: Args) {
		return this.tts(message, args.getOption('prompt') || 'NOTHING');
	}

	// Chat Input (slash) command
	public async chatInputRun(interaction: Command.ChatInputCommandInteraction) {
		return this.tts(interaction, interaction.options.getString('prompt') || 'NOTHING');
	}

	private async tts(interactionOrMessage: Message | Command.ChatInputCommandInteraction | Command.ContextMenuCommandInteraction, prompt: string) {
		const askMessage =
			interactionOrMessage instanceof Message
				? await interactionOrMessage.channel.send({ content: '🤔 Thinking... 🤔' })
				: await interactionOrMessage.reply({ content: '🤔 Thinking... 🤔', fetchReply: true });

		try {
			enum voice {
				alloy = 'alloy',
				echo = 'echo',
				fable = 'fable',
				onyx = 'onyx',
				nova = 'nova',
				shimmer = 'shimmer'
			}

			const voices = [voice.alloy, voice.echo, voice.fable, voice.onyx, voice.nova, voice.shimmer];
			const mp3 = await openai.audio.speech.create({
				model: 'tts-1-hd',
				voice: voices[Math.floor(Math.random() * voices.length)],
				input: prompt
			});
			const mp3Buffer = Buffer.from(await mp3.arrayBuffer());

			const mp3Attachment: AttachmentBuilder[] = [];

			mp3Attachment.push(
				new AttachmentBuilder(Buffer.from(new Uint8Array(mp3Buffer)), {
					name: 'himbot_response.mp3',
					description: `An TTS message generated by Himbot using the prompt: ${prompt}`
				})
			);

			const content = `Prompt: ${prompt}:`;

			if (interactionOrMessage instanceof Message) {
				return askMessage.edit({
					content,
					files: mp3Attachment,
					flags: '8192'
				});
			}

			return interactionOrMessage.editReply({
				content,
				files: mp3Attachment,
				options: {
					flags: MessageFlags.IsVoiceMessage.valueOf()
				}
			});
		} catch (error) {
			const content = "Sorry, I can't complete the prompt for: " + prompt + '\n' + error;

			if (interactionOrMessage instanceof Message) {
				return askMessage.edit({
					content
				});
			}

			return interactionOrMessage.editReply({
				content
			});
		}
	}
}