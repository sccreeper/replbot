package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	dg "github.com/bwmarrin/discordgo"
)

var bot_token string
var remove_commands bool = false
var start_time time.Time

func init() {

	bot_token = os.Getenv("BOT_TOKEN")
	start_time = time.Now()

	flag.BoolVar(&remove_commands, "rmcmd", false, "Remove commands")
	flag.Parse()

}

// Just init the bot.
func main() {

	// Bot auth and login

	if bot_token == "" {
		fmt.Println("No token provided")
		os.Exit(1)
	}

	bot, err := dg.New("Bot " + bot_token)
	if err != nil {
		panic(err)
	}

	bot.Identify.Intents = 1536

	err = bot.Open()
	check_error(err)

	defer bot.Close()

	// Handlers

	bot.AddHandler(receive_message)

	// Create commands

	log.Println("Adding commands...")
	registeredCommands := make([]*dg.ApplicationCommand, len(commands))

	for i, v := range commands {
		cmd, err := bot.ApplicationCommandCreate(bot.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	bot.AddHandler(func(s *dg.Session, i *dg.InteractionCreate) {
		if h, ok := command_handlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	go clear_sessions()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if remove_commands {
		log.Println("Removing commands...")

		for _, v := range registeredCommands {
			err := bot.ApplicationCommandDelete(bot.State.User.ID, "", v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

}
