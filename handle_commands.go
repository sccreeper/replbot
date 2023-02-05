// Handles evaluate, start & end.

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	dg "github.com/bwmarrin/discordgo"
	"github.com/robertkrimen/otto"
)

var (
	js_sessions map[string]Session = map[string]Session{}
)

type Session struct {
	VM          *otto.Otto
	TimeCreated int64
}

func handle_eval(s *dg.Session, i *dg.InteractionCreate) {

	log.Println(gen_interaction_log("evaluate", i))

	options := i.ApplicationCommandData().Options
	mode := options[0].IntValue()
	code_string := options[1].StringValue()

	switch mode {
	case int64(JSModeOption):
		var value_string string

		// Create context

		vm, err := new_otto()
		if err != nil {

			value_string = fmt.Sprintf("🔴 **Error:** `%s`", trunc_err(err.Error()))

		} else {

			value, err := vm.Run(code_string)

			if err != nil {
				value_string = fmt.Sprintf("🔴 **Error:** `%s`", trunc_err(err.Error()))

			} else {

				value_string, err = value.ToString()
				check_error(err)
			}
		}

		s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &dg.InteractionResponseData{
				Content: value_string,
			},
		})
	case int64(SessionModeOption):

		if _, ok := js_sessions[resolve_user(i).ID]; !ok {
			s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &dg.InteractionResponseData{
					Content: "🔴 You do not have a running session.",
				},
			})
		} else {

			var value_string string

			val, err := js_sessions[resolve_user(i).ID].VM.Run(code_string)
			if err != nil {
				value_string = fmt.Sprintf("🔴 **Error:** `%s`", trunc_err(err.Error()))
			} else {

				value_string, _ = val.ToString()
			}

			s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &dg.InteractionResponseData{
					Content: value_string,
				},
			})

		}
	default:
		s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &dg.InteractionResponseData{
				Content: "Unrecognized option. Try restarting your Discord client. If the problem persists contact the maker of this bot.",
			},
		})

	}

}

func handle_start(s *dg.Session, i *dg.InteractionCreate) {

	log.Print(gen_interaction_log("start", i))

	// See if user already has session.

	var error_string string
	var successful bool = true

	if _, ok := js_sessions[i.Member.User.ID]; ok {
		successful = false
		error_string = "🔴 You already have a session. Please end it before starting a new one."
	} else {

		switch i.ApplicationCommandData().Options[0].IntValue() {
		case int64(JSLanguageOption):
			vm, err := new_otto()
			if err != nil {
				successful = false
				error_string = "🔴 There was an error when starting your session."
			} else {
				successful = true
				js_sessions[i.Member.User.ID] = Session{VM: vm, TimeCreated: time.Now().Unix()}
			}
		}

	}

	if !successful {
		s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &dg.InteractionResponseData{
				Content: error_string,
			},
		})
	} else {

		s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &dg.InteractionResponseData{
				Content: "🟢 Session started successfully.",
			},
		})

	}

}

func handle_end(s *dg.Session, i *dg.InteractionCreate) {

	log.Println(gen_interaction_log("end", i))

	var resp_string string

	if _, ok := js_sessions[i.Member.User.ID]; !ok {
		resp_string = "🔴 You do not have a running session."
	} else {
		delete(js_sessions, i.Member.User.ID)
		resp_string = "🟢 Session closed successfully."
	}

	s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &dg.InteractionResponseData{
			Content: resp_string,
		},
	})

}
