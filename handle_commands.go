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
	History         string
}

// Evaluates a code string
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

			user_session.History += fmt.Sprintf(">>> %s\n", code_string)

			if err != nil {
				user_session.History += fmt.Sprintf("%s\n", trunc_err(err.Error()))
				js_sessions[resolve_user(i).ID] = user_session

				resp_string = fmt.Sprintf("🔴 **Error:** `%s`", trunc_err(err.Error()))
			} else {

				str_val, err := val.ToString()
				check_error(err)

				user_session.History += fmt.Sprintf("%s\n", str_val)
				js_sessions[resolve_user(i).ID] = user_session

				resp_string = fmt.Sprintf("```%s```", str_val)
			}

			default_resp(s, i.Interaction, resp_string, "Evaluating")

		}
	default:
		default_resp(s, i.Interaction, "🔴 Unrecognized option. Try restarting your Discord client. If the problem persists contact the maker of this bot.", "Error")
	}

}

// Starts a session.
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

// Kills a session
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

// Bot about page
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

// Help menu with all commands explained in more detail.
func handle_help(s *dg.Session, i *dg.InteractionCreate) {

	s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
		Type: dg.InteractionResponseChannelMessageWithSource,
		Data: &dg.InteractionResponseData{
			Embeds: []*dg.MessageEmbed{
				{
					Title:       "❓ Help",
					Description: "Help menu",
					Fields: []*dg.MessageEmbedField{
						{Name: "/help", Value: "Help command"},
						{Name: "/evaluate <mode> <code>", Value: "Evaluates code in a REPL session or on it's own."},
						{Name: "/start <language>", Value: "Start a REPL session with a specified language. Sessions timeout after 5 minutes of inactivity."},
						{Name: "/end", Value: "Ends your REPL session."},
						{Name: "/info", Value: "About page for the bot."},
						{Name: "/history", Value: "Displays evaluation history for a session"},
						{Name: "/clear", Value: "Clears the evaluation history for a session. This does not include any values that have been declared."},
					},
					Color: 1976635,
				},
			},
		},
	})

}

// Shows session history
func handle_history(s *dg.Session, i *dg.InteractionCreate) {

	if _, ok := js_sessions[resolve_user(i).ID]; !ok {
		default_resp(s, i.Interaction, "🔴 You do not have a running session.", "Error")
	} else {

		time_started := time.Unix(js_sessions[resolve_user(i).ID].TimeCreated, 0)

		var history_string string

		if js_sessions[resolve_user(i).ID].History == "" {
			history_string = "```No history in session!```"
		} else {
			history_string = fmt.Sprintf("```%s```", js_sessions[resolve_user(i).ID].History)
		}

		s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
			Type: dg.InteractionResponseChannelMessageWithSource,
			Data: &dg.InteractionResponseData{
				Embeds: []*dg.MessageEmbed{
					{
						Title:       "Session History",
						Description: fmt.Sprintf("Session started at %s", time_started.String()),
						Fields: []*dg.MessageEmbedField{
							{Name: "History", Value: history_string},
						},
						Color: 1976635,
					},
				},
			},
		})

	}

}

func handle_clear(s *dg.Session, i *dg.InteractionCreate) {

	if _, ok := js_sessions[resolve_user(i).ID]; !ok {
		default_resp(s, i.Interaction, "🔴 You do not have a running session.", "Error")
	} else {

		user_session := js_sessions[resolve_user(i).ID]
		user_session.History = ""
		js_sessions[resolve_user(i).ID] = user_session

		default_resp(s, i.Interaction, "🟢 History cleared successfully.", "Clearing history")

	}

}
