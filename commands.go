package main

import dg "github.com/bwmarrin/discordgo"

const (
	JSModeOption      int = 1
	SessionModeOption int = 2
)

const (
	PythonLanguageOption int = 1
	JSLanguageOption     int = 2
)

var (
	commands = []*dg.ApplicationCommand{
		{
			Name:        "start",
			Description: "Start your session",
			Options: []*dg.ApplicationCommandOption{
				{
					Name:        "language",
					Description: "Language to use",
					Required:    true,
					Type:        dg.ApplicationCommandOptionInteger,
					Choices: []*dg.ApplicationCommandOptionChoice{
						{
							Name:  "JavaScript",
							Value: JSLanguageOption,
						},
					},
				},
			},
		},
		{
			Name:        "end",
			Description: "End your session",
		},
		{
			Name:        "evaluate",
			Description: "Evaluate a string in your session, or just evaluate a single string.",
			Options: []*dg.ApplicationCommandOption{
				{
					Name:        "mode",
					Description: "Language to use or evaluate in session",
					Type:        dg.ApplicationCommandOptionInteger,
					Required:    true,
					Choices: []*dg.ApplicationCommandOptionChoice{
						{
							Name:  "javascript",
							Value: JSModeOption,
						},
						{
							Name:  "session",
							Value: SessionModeOption,
						},
					},
				},
				{
					Name:        "code",
					Description: "Code to evaluate/run",
					Type:        dg.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
	}

	command_handlers = map[string]func(s *dg.Session, i *dg.InteractionCreate){
		"evaluate": handle_eval,
		"end":      handle_end,
		"start":    handle_start,
	}
)
