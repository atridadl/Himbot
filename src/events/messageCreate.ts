import { events } from "./mod.ts";
import { logger } from "../utils/logger.ts";

const log = logger({ name: "Event: messageCreate" });

events.messageCreate = (bot, message) => {
  log.info(`${message.tag}: ${message.content}`);
};
