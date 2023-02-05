// Handles evaluate, start & end.

package main

import (
	"fmt"
	"log"
	"time"

	dg "github.com/bwmarrin/discordgo"
	"github.com/robertkrimen/otto"
)

var (
	js_sessions map[string]Session = map[string]Session{}
)

const (
	MaxCodeLength int = 512
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

	if len(code_string) > MaxCodeLength {
		default_resp(s, i.Interaction, fmt.Sprintf("🔴 Code is too long (max allowed is %d). Yours was %d char(s) long.", MaxCodeLength, len(code_string)))
	}

	switch mode {
	case int64(JSModeOption):
		var value_string string

		// Create context

		vm, err := new_otto()
		if err != nil {

			value_string = fmt.Sprintf("🔴 **Error:** `%s`", trunc_err(err.Error()))

		} else {

			value, err := run_code(vm, code_string, 2)

			if err != nil {
				value_string = fmt.Sprintf("🔴 **Error:** `%s`", trunc_err(err.Error()))

			} else {

				value_string, err = value.ToString()
				check_error(err)
			}
		}

		default_resp(s, i.Interaction, value_string)

	case int64(SessionModeOption):

		if _, ok := js_sessions[resolve_user(i).ID]; !ok {

			default_resp(s, i.Interaction, "🔴 You do not have a running session.")

		} else {

			var resp_string string

			val, err := run_code(js_sessions[resolve_user(i).ID].VM, code_string, 2)

			if err != nil {
				resp_string = fmt.Sprintf("🔴 **Error:** `%s`", trunc_err(err.Error()))
			} else {

				resp_string, _ = val.ToString()
			}

			default_resp(s, i.Interaction, resp_string)

		}
	default:
		default_resp(s, i.Interaction, "Unrecognized option. Try restarting your Discord client. If the problem persists contact the maker of this bot.")
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

		default_resp(s, i.Interaction, error_string)

	} else {

		default_resp(s, i.Interaction, "🟢 Session started successfully.")

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

	default_resp(s, i.Interaction, resp_string)

}
