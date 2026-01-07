## Description
The project was developed to learn the Golang and standard packages. A Telegram bot for saving user notes was implemented.
## Telegram Bot Functionality
1. Saves user notes by sending a message in the chat.
2. View notes using the /check /check3 commands
3. Delete notes using /delete
4. Get help using /help

## Bot Architecture
Main parts: head, fetcher, and processer.
1. head - Manages the fetcher and processer, an infinite loop.
2. fetcher - Accesses the Telegram API and retrieves messages sent to the bot.
3. processor - Processes user messages, accesses the database, and sends messages from the bot.
4. sqLite3 was chosen as the database.
The project was designed with the idea of ​​extending to other messengers, so interfaces were implemented.

## Launching the application
1. Clone the git repository: 
```bash
git clone https://github.com/DimitryAn/ProjectTelegramBot.git
```
2. Download the Golang compiler: https://go.dev/doc/install
3. In the project folder, run 
```bash
go mod download
```
This will install the driver for sqLite.
4. Build the project: 
```bash
go build
```
5. To launch the project, you need to pass flags with your bot's token and the messenger host.
For example, for Telegram: 
```bash
./bot -tgToken 'you're token' -host 'api.telegram.org'
```
