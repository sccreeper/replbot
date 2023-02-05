// Handles evaluate

package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	dg "github.com/bwmarrin/discordgo"
	"github.com/robertkrimen/otto"
)

var sessions map[int64]otto.Otto

type Session struct {
	VM          otto.Otto
	TimeCreated uint64
}

func init() {
	sessions = make(map[int64]otto.Otto)
}

func handle_eval(s *dg.Session, i *dg.InteractionCreate) {

	options := i.ApplicationCommandData().Options
	mode := options[0].IntValue()
	code_string := options[1].StringValue()

	if mode == int64(JSModeOption) {

		// Create context

		vm := otto.New()

		// Redirect log output to just return value.

		vm.Set("__log__", func(call otto.FunctionCall) otto.Value {

			output_string := ""

			for _, v := range call.ArgumentList {
				output_string = fmt.Sprintf("%s %s", output_string, v.String())
			}

			v, _ := otto.ToValue(output_string)

			return v

		})

		vm.Run("console.log = __log__;")

		// Set window interval & timeout

		msg, err := otto.ToValue("setInterval is not supported.")
		check_error(err)
		vm.Set("setInterval", func(call otto.FunctionCall) otto.Value {
			return msg
		})

		msg, err = otto.ToValue("setTimeout is not supported.")
		check_error(err)
		vm.Set("setTimeout", func(call otto.FunctionCall) otto.Value {
			return msg
		})

		value, err := vm.Run(code_string)

		var value_string string

		if err != nil {

			value_string = fmt.Sprintf("**Error:** `%s`", err.Error())

		} else {

			value_string, err = value.ToString()
			check_error(err)
		}

		s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &dg.InteractionResponseData{
				Content: value_string,
			},
		})

	} else {
		s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &dg.InteractionResponseData{
				Content: "Unrecognized option. Try restarting your Discord client. If the problem persists contact the maker of this bot.",
			},
		})
	}

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
