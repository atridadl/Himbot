import { ApplyOptions } from '@sapphire/decorators';
import { Args, BucketScope, Command } from '@sapphire/framework';
import { AttachmentBuilder, Message } from 'discord.js';

@ApplyOptions<Command.Options>({
	description: 'Make a picture!',
	options: ['prompt'],
	// 10mins
	cooldownDelay: 300_000,
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

		const response = await fetch(`https://api.stability.ai/v1/generation/stable-diffusion-xl-beta-v2-2-2/text-to-image`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Accept: 'application/json',
				Authorization: `Bearer ${process.env.STABILITY_API_KEY}`
			},
			body: JSON.stringify({
				text_prompts: [
					{
						text: prompt
					}
				],
				cfg_scale: 7,
				clip_guidance_preset: 'FAST_BLUE',
				height: 512,
				width: 512,
				samples: 1,
				steps: 50
			})
		});

		interface GenerationResponse {
			artifacts: Array<{
				base64: string;
				seed: number;
				finishReason: string;
			}>;
		}

		if (!response.ok) {
			const content = `Sorry! I goofed up. Please ask my maker HimbothySwaggins about what could have happened!`;

			if (interactionOrMessage instanceof Message) {
				return askMessage.edit({ content });
			}

			return interactionOrMessage.editReply({
				content: content
			});
		} else {
			const responseJSON = (await response.json()) as GenerationResponse;
			const imageAttachment = new AttachmentBuilder(Buffer.from(responseJSON.artifacts[0].base64, 'base64'));

			const content = `Prompt: ${prompt}` || 'ERROR!';

			if (interactionOrMessage instanceof Message) {
				return askMessage.edit({ content, files: [imageAttachment] });
			}

			return interactionOrMessage.editReply({
				content: content,
				files: [imageAttachment]
			});
		}
	}
}
