import './lib/setup';
import { LogLevel, SapphireClient } from '@sapphire/framework';
import { ActivityType, GatewayIntentBits } from 'discord.js';

const client = new SapphireClient({
	defaultPrefix: '!',
	presence: {
		status: 'online',
		activities: [
			{
				name: "I'm just here for the vibes my dudes...",
				type: ActivityType.Playing
			}
		]
	},
	caseInsensitiveCommands: true,
	logger: {
		level: LogLevel.Debug
	},
	intents: [GatewayIntentBits.DirectMessages, GatewayIntentBits.GuildMessages, GatewayIntentBits.Guilds, GatewayIntentBits.MessageContent],
	loadMessageCommandListeners: true
});

const main = async () => {
	try {
		client.logger.info('Logging in');
		await client.login();
		client.logger.info('logged in');
	} catch (error) {
		client.logger.fatal(error);
		client.destroy();
		process.exit(1);
	}
};

main();
