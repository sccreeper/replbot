<img src="./profile.png" width="128"/>

# Discord REPL bot

A bot that allows you to interact with a REPL environment in your server.

Add the bot to your server - [click here](https://discord.com/oauth2/authorize?client_id=1071482433540735027&scope=bot&permissions=534791060544)

---

### Supported languages

- JavaScript
- Lua (future)

### Commands

- `/help` - Help menu
- `/evaluate <mode> <code>` Evaluates code in a REPL session or on it's own.
- `/start <language>` Start a REPL session with a specified language. Sessions timeout after 5 minutes of inactivity.
- `/end` Ends your REPL session.
- `/info` About page for the bot.
- `/history` Displays evaluation history for a session
- `/clear` Clears the evaluation history for a session. This does not include any values that have been declared.

### Credits

- otto - [robertkrimen/otto](https://github.com/robertkrimen/otto)
    - Language runtime used for JS.
- discordgo [bwmarrin/discordgo](https://github.com/bwmarrin/discordgo)
  - Library used to interact with the Discord API.

