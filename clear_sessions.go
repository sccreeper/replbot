package main

import (
	"log"
	"sync"
	"time"
)

var lock = sync.RWMutex{}

const timeout_limit int64 = 5 * 60

func clear_sessions() {

	for {

		for k, v := range js_sessions {

			lock.Lock()

			if (time.Now().Unix() - v.LastInteraction) > timeout_limit {
				delete(js_sessions, k)
				log.Printf("Timeout in session '%s'", k)
			}

			lock.Unlock()

		}

		time.Sleep(3 * time.Second)

	}

}
