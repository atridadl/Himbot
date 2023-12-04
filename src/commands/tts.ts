import { BucketScope, Command, container } from '@sapphire/framework';
import { AttachmentBuilder, MessageFlags } from 'discord.js';
import OpenAI from 'openai';

const openai = new OpenAI({
	apiKey: process.env.OPENAI_API_KEY
});

export class TTSCommand extends Command {
	public constructor(context: Command.LoaderContext) {
		super(context, {
			description: 'Generate TTS every minute!',
			options: ['prompt'],
			cooldownDelay: 100_000,
			cooldownLimit: 1,
			// Yes... I did hardcode myself.
			cooldownFilteredUsers: ['83679718401904640'],
			cooldownScope: BucketScope.User
		});
	}

	// Register Chat Input and Context Menu command
	public override registerApplicationCommands(registry: Command.Registry) {
		registry.registerChatInputCommand((builder) =>
			builder
				.setName(this.name)
				.setDescription(this.description)
				.addStringOption((option) => option.setName('prompt').setDescription('The prompt you will use to generate audio!').setRequired(true))
		);
	}

	// Chat Input (slash) command
	public async chatInputRun(interaction: Command.ChatInputCommandInteraction) {
		const prompt = interaction.options.getString('prompt') || 'NOTHING';

		await interaction.reply({ content: '🤔 Thinking... 🤔', fetchReply: true });

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
				model: 'tts-1',
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

			return interaction.editReply({
				content,
				files: mp3Attachment,
				options: {
					flags: MessageFlags.IsVoiceMessage.valueOf()
				}
			});
		} catch (error) {
			const content = "Sorry, I can't complete the prompt for: " + prompt + '\n' + error;

			return interaction.editReply({
				content
			});
		}
	}
}

void container.stores.loadPiece({
	store: 'commands',
	name: 'tts',
	piece: TTSCommand
});
