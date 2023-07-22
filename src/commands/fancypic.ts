import { ApplyOptions } from '@sapphire/decorators';
import { Args, BucketScope, Command } from '@sapphire/framework';
import { AttachmentBuilder, Message } from 'discord.js';

// This is literally the worlds messiest TS code. Please don't judge me...

@ApplyOptions<Command.Options>({
	description: 'Make a picture... but high res!',
	options: ['prompt', 'number of pictures'],
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
				.addStringOption((option) =>
					option.setName('prompt').setDescription('The prompt you will use to generate an image!').setRequired(true)
				)
				.addStringOption((option) =>
					option
						.setName('amount')
						.setDescription('The number of images you would like to generate. Maximum 2.')
						.setChoices(
							...[
								{ name: '1', value: '1' },
								{ name: '2', value: '2' }
							]
						)
				)
		);
	}

	// Message command
	public async messageRun(message: Message, args: Args) {
		const amount = Math.abs(Number(args.getOption('amount')));
		return this.picHr(message, args.getOption('prompt') || 'Scold me for not passing any prompt in.', amount <= 2 ? amount : 1);
	}

	// Chat Input (slash) command
	public async chatInputRun(interaction: Command.ChatInputCommandInteraction) {
		const amount = Number(interaction.options.getString('amount'));
		return this.picHr(interaction, interaction.options.getString('prompt') || 'NOTHING', amount || 1);
	}

	private async picHr(
		interactionOrMessage: Message | Command.ChatInputCommandInteraction | Command.ContextMenuCommandInteraction,
		prompt: string,
		amount: number
	) {
		const askMessage =
			interactionOrMessage instanceof Message
				? await interactionOrMessage.channel.send({ content: '🤔 Thinking... 🤔' })
				: await interactionOrMessage.reply({ content: '🤔 Thinking... 🤔', fetchReply: true });

		const creditCountResponse = await fetch(`https://api.stability.ai/v1/user/balance`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${process.env.STABILITY_API_KEY}`
			}
		});

		const balance = (await creditCountResponse.json()).credits || 0;

		if (balance > 5) {
			const imageGenResponse = await fetch(`https://api.stability.ai/v1/generation/stable-diffusion-xl-1024-v0-9/text-to-image`, {
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
					cfg_scale: 6,
					clip_guidance_preset: 'FAST_BLUE',
					height: 1024,
					width: 1024,
					samples: amount,
					steps: 32,
					seed: Number(String(interactionOrMessage.member?.user.id).substring(0, 5)) || 0
				})
			});

			interface GenerationResponse {
				artifacts: Array<{
					base64: string;
					seed: number;
					finishReason: string;
				}>;
			}

			if (!imageGenResponse.ok) {
				const content = `Sorry, I can't complete the prompt for: ${prompt}`;

				if (interactionOrMessage instanceof Message) {
					return askMessage.edit({ content });
				}

				return interactionOrMessage.editReply({
					content: content
				});
			} else {
				const responseJSON = (await imageGenResponse.json()) as GenerationResponse;
				const imageAttachment: AttachmentBuilder[] = [];

				for (let i = 0; i < responseJSON.artifacts.length; i++) {
					imageAttachment.push(
						new AttachmentBuilder(Buffer.from(responseJSON.artifacts[i].base64, 'base64'), {
							name: 'response.jpg',
							description: "Himbot's Response"
						})
					);
				}

				const newCreditCountResponse = await fetch(`https://api.stability.ai/v1/user/balance`, {
					method: 'GET',
					headers: {
						'Content-Type': 'application/json',
						Authorization: `Bearer ${process.env.STABILITY_API_KEY}`
					}
				});

				const newBalance = (await newCreditCountResponse.json()).credits || 0;

				const content =
					`Credits Used: ${balance - newBalance}\nPrompt: ${prompt}${
						balance <= 300
							? `\n\n⚠️I am now at ${balance} credits. If you'd like to help fund this command, please type "/support" for details!`
							: ''
					}` || 'ERROR!';

				if (interactionOrMessage instanceof Message) {
					return askMessage.edit({ content, files: imageAttachment });
				}

				return interactionOrMessage.editReply({
					content,
					files: imageAttachment
				});
			}
		} else {
			const content = `Oops! We're out of credits for this. If you'd like to help fund this command, please type "/support" for details!`;

			if (interactionOrMessage instanceof Message) {
				return askMessage.edit({ content });
			}

			return interactionOrMessage.editReply({
				content: content
			});
		}
	}
}
