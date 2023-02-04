package main

import (
	"github.com/bwmarrin/discordgo"
)

func receive_message(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Content[0] != '!' {
		return
	} else if m.Author.ID == s.State.User.ID {
		return
	}

	message_string := m.Content[1:]

	switch message_string {
	case "ping":
		s.ChannelMessageSend(m.ChannelID, "pong")
	}

}
