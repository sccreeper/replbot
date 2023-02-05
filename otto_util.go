package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/robertkrimen/otto"
)

var (
	ErrCodeTimeout error = errors.New("code took too long to run")
	ErrHalt        error = errors.New("halt")
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

// Mostly taken from otto README.md
func run_code(vm *otto.Otto, code string, timeout int64) (result otto.Value, e error) {

	start := time.Now()

	defer func() {
		duration := time.Since(start)
		if caught := recover(); caught != nil {
			if caught == ErrHalt {
				log.Printf("Some code took to long! Stopping after: %v\n", duration)
				result = otto.NullValue()
				e = ErrCodeTimeout
				return
			}
			panic(caught) // Something else happened, repanic!
		}
	}()

	vm.Interrupt = make(chan func(), 1) // The buffer prevents blocking
	watchdogCleanup := make(chan struct{})
	defer close(watchdogCleanup)

	go func() {
		select {
		case <-time.After(time.Duration(timeout) * time.Second): // Stop after two seconds
			vm.Interrupt <- func() {
				panic(ErrHalt)
			}
		case <-watchdogCleanup:
		}
		close(vm.Interrupt)
	}()

	val, err := vm.Run(code) // Here be dragons (risky code)
	return val, err

}
