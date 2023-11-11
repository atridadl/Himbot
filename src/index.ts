import './lib/setup';
import { LogLevel, SapphireClient, BucketScope } from '@sapphire/framework';
import { ActivityType, GatewayIntentBits } from 'discord.js';

const client = new SapphireClient({
	defaultPrefix: '!',
	presence: {
		status: 'online',
		activities: [
			{
				name: 'idk',
				type: ActivityType.Custom
			}
		]
	},
	caseInsensitiveCommands: true,
	logger: {
		level: LogLevel.Debug
	},
	intents: [GatewayIntentBits.DirectMessages, GatewayIntentBits.GuildMessages, GatewayIntentBits.Guilds, GatewayIntentBits.MessageContent],
	loadMessageCommandListeners: true,
	defaultCooldown: {
		// 10s
		delay: 10_000,
		filteredCommands: ['support', 'ping', 'wryna'],
		limit: 2,
		// Yes... I did hardcode myself.
		filteredUsers: ['83679718401904640'],
		scope: BucketScope.User
	}
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
