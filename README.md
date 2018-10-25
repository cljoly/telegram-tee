# telegram-cli-bot
Simple cli tool to send stdin to any Telegram chat, through a bot. It’s a bit like `tee`.

## Getting started

### Set up

First, install the tool with
``` bash
go get -u github.com/leowzukw/telegram-tee
```

Then, you need to control a bot. Set the environment variable `TLGCLI_TOKEN` to
the token of the bot that will write stdin to a chat for you. You may want to [create a new bot](https://core.telegram.org/bots#3-how-do-i-create-a-bot) or use an existing one.

### Use

Then, you need to get the chat ID of the conversation to which you want to send stdin. Just run
``` bash
telegram-tee
```
and write with Telegram to your bot. It will reply with the current chatID.

You can then do
``` bash
echo Hi | telegram-tee <chatID>
```

You can even send to several chatID at the same time.

## TODO

- [ ] Make stdout usable like tee
