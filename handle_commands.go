// Handles evaluate, start & end.

package main

import (
	"fmt"
	"log"
	"math"
	"runtime"
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
	VM              *otto.Otto
	TimeCreated     int64
	LastInteraction int64
}

func handle_eval(s *dg.Session, i *dg.InteractionCreate) {

	log.Println(gen_interaction_log("evaluate", i))

	options := i.ApplicationCommandData().Options
	mode := options[0].IntValue()
	code_string := options[1].StringValue()

	if len(code_string) > MaxCodeLength {
		default_resp(s, i.Interaction, fmt.Sprintf("🔴 Code is too long (max allowed is %d). Yours was %d char(s) long.", MaxCodeLength, len(code_string)), "Code too long")
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

				str_val, err := value.ToString()
				check_error(err)
				value_string = fmt.Sprintf("```%s```", str_val)

			}
		}

		default_resp(s, i.Interaction, value_string, "Evaluate")

	case int64(SessionModeOption):

		if _, ok := js_sessions[resolve_user(i).ID]; !ok {

			default_resp(s, i.Interaction, "🔴 You do not have a running session.", "Error")

		} else {

			user_session := js_sessions[resolve_user(i).ID]
			user_session.LastInteraction = time.Now().Unix()
			js_sessions[resolve_user(i).ID] = user_session

			var resp_string string

			val, err := run_code(js_sessions[resolve_user(i).ID].VM, code_string, 2)

			if err != nil {
				resp_string = fmt.Sprintf("🔴 **Error:** `%s`", trunc_err(err.Error()))
			} else {

				str_val, err := val.ToString()
				check_error(err)
				resp_string = fmt.Sprintf("```%s```", str_val)
			}

			default_resp(s, i.Interaction, resp_string, "Evaluating")

		}
	default:
		default_resp(s, i.Interaction, "🔴 Unrecognized option. Try restarting your Discord client. If the problem persists contact the maker of this bot.", "Error")
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
				js_sessions[i.Member.User.ID] = Session{VM: vm, TimeCreated: time.Now().Unix(), LastInteraction: time.Now().Unix()}
			}
		}

	}

	if !successful {

		default_resp(s, i.Interaction, error_string, "Starting session")

	} else {

		default_resp(s, i.Interaction, "🟢 Session started successfully.", "Starting session")

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

	default_resp(s, i.Interaction, resp_string, "Ending session")

}

func handle_info(s *dg.Session, i *dg.InteractionCreate) {

	log.Println(gen_interaction_log("info", i))

	var mem_stats runtime.MemStats
	runtime.ReadMemStats(&mem_stats)

	time_difference := time.Since(start_time)

	hours := time_difference.Hours()
	minutes := math.Mod(time_difference.Minutes(), 60)
	seconds := math.Mod(time_difference.Seconds(), 60)

	default_resp(s, i.Interaction,
		fmt.Sprintf(
			`
__About__

Provides a REPL environment to experiment with various scripting languages.
Made by oscarcp.

__Bot stats__

**Current sessions:** %d
**Memory usage:** %.2fmb
**Uptime:** %s

	
	`,
			len(js_sessions),
			float64(mem_stats.Sys)/math.Pow10(6),
			fmt.Sprintf("%d hour(s) %d minute(s) %d second(s)",
				int(hours),
				int(minutes),
				int(seconds),
			),
		),
		"Bot Info",
	)

}
