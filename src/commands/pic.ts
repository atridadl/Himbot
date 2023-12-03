import { ApplyOptions } from '@sapphire/decorators';
import { Args, BucketScope, Command } from '@sapphire/framework';
import { Message } from 'discord.js';
import Replicate from 'replicate';

const replicate = new Replicate({
	auth: process.env.REPLICATE_API_TOKEN
});

// @ts-ignore
@ApplyOptions<Command.Options>({
	description: 'Generate an image using Stability AI! Cooldown 1 Minute to prevent spam!',
	options: ['prompt'],
	// 10mins
	cooldownDelay: 100_000,
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
		return this.picHr(message, args.getOption('prompt') || 'Scold me for not passing any prompt in.');
	}

	// Chat Input (slash) command
	public async chatInputRun(interaction: Command.ChatInputCommandInteraction) {
		return this.picHr(interaction, interaction.options.getString('prompt') || 'NOTHING');
	}

	private async picHr(interactionOrMessage: Message | Command.ChatInputCommandInteraction | Command.ContextMenuCommandInteraction, prompt: string) {
		const askMessage =
			interactionOrMessage instanceof Message
				? await interactionOrMessage.channel.send({ content: '🤔 Thinking... 🤔' })
				: await interactionOrMessage.reply({ content: '🤔 Thinking... 🤔', fetchReply: true });

		let result = (await replicate.run('stability-ai/sdxl:39ed52f2a78e934b3ba6e2a89f5b1c712de7dfea535525255b1aa35c5565e08b', {
			input: {
				prompt,
				disable_safety_checker: true,
				refine: 'expert_ensemble_refiner',
				num_inference_steps: 50,
				scheduler: 'K_EULER',
				width: 1024,
				height: 1024
			}
		})) as string[];

		if (result.length <= 0) {
			const content = `Sorry, I can't complete the prompt for: ${prompt}`;

			if (interactionOrMessage instanceof Message) {
				return askMessage.edit({ content });
			}

			return interactionOrMessage.editReply({
				content: content
			});
		} else {
			const content = `Prompt: ${prompt}\n URL: ${result[0]}`;

			if (interactionOrMessage instanceof Message) {
				return askMessage.edit({ content });
			}

			return interactionOrMessage.editReply({
				content
			});
		}
	}
}
