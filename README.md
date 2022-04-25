# Telereto

Telereto, written in Go, is an easy-to-use Telegram bot for Chevereto V4 image hosted site.

## Prerequisite
- Chevereto V4 (of course).
- URL uploads should be enabled on Chevereto side. You can enable it on *your.chevereto.site/dashboard/settings/image-upload*. 
    - Kindly note that this is not a privacy issue in this case because IP of either your bot server or chevereto server 
  must be exposed to Telegram.
- A reverse proxy, e.g. Nginx, is required as a frontend of the bot.

## Config file
An example of the config file, *config.yml.example*, is included in the repository.

## Running
```bash
telereto --config /path/to/yout/config
```