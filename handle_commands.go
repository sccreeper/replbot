// Handles evaluate

package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	dg "github.com/bwmarrin/discordgo"
)

func handle_eval(s *dg.Session, i *dg.InteractionCreate) {

	options := i.ApplicationCommandData().Options
	mode := options[0].IntValue()
	code_string := options[1].StringValue()

	fmt.Println(mode, code_string)

	s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &dg.InteractionResponseData{
			Content: "Hello World",
		},
	})

}

func handle_start(s *dg.Session, i *dg.InteractionCreate) {

	s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &dg.InteractionResponseData{
			Content: "Hello World",
		},
	})

}

func handle_end(s *dg.Session, i *dg.InteractionCreate) {

	s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &dg.InteractionResponseData{
			Content: "Hello World",
		},
	})

}
