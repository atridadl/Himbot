import {
  ApplicationCommandOptionTypes,
  ApplicationCommandTypes,
  InteractionResponseTypes,
} from "../../deps.ts";

import { createCommand } from "./mod.ts";

// Thanks for the bot name idea Wryna!

createCommand({
  name: "wryna",
  description: "What was your nickname in highschool?",
  type: ApplicationCommandTypes.ChatInput,
  scope: "Global",
  options: [
    {
      type: ApplicationCommandOptionTypes.String,
      name: "input",
      description: "Text you would like to send to this command.",
      required: true,
    },
  ],
  execute: async (bot, interaction) => {
    const input = interaction.data?.options?.find(
      (option) => option.name === "input"
    );
    await bot.helpers.sendInteractionResponse(
      interaction.id,
      interaction.token,
      {
        type: InteractionResponseTypes.ChannelMessageWithSource,
        data: {
          content: `${input?.value} was ${interaction.user.username}'s nickname in highschool.`,
        },
      }
    );
  },
});
