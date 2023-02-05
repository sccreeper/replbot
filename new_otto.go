package main

import (
	"fmt"

	"github.com/robertkrimen/otto"
)

// Creates and initializes a new VM.
func new_otto() (*otto.Otto, error) {
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
	if err != nil {
		return nil, err
	}

	vm.Set("setInterval", func(call otto.FunctionCall) otto.Value {
		return msg
	})

	msg1, err := otto.ToValue("setTimeout is not supported.")
	if err != nil {
		return nil, err
	}

	vm.Set("setTimeout", func(call otto.FunctionCall) otto.Value {
		return msg1
	})

	return vm, nil
}
