package main

import (
	"fmt"

	dg "github.com/bwmarrin/discordgo"
)

const (
	err_max int = 128
)

func check_error(e error) {

	if e != nil {
		panic(e)
	}

}

func trunc_err(e string) string {

	if len(e) > err_max-3 {
		return e[:err_max-3] + "..."
	} else {
		return e
	}

}

func gen_interaction_log(command_name string, i *dg.InteractionCreate) string {

	var user_obj *dg.User

	if i.User == nil {
		user_obj = i.Member.User
	} else {
		user_obj = i.User
	}

	return fmt.Sprintf("/%s used by %s in %s", command_name, user_obj.ID, i.ChannelID)

}

func resolve_user(i *dg.InteractionCreate) *dg.User {

	if i.User == nil {
		return i.Member.User
	} else {
		return i.User
	}

}
