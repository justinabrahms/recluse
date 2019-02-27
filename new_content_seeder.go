package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Publishes new posts to a byte channel
func newMessagePublisher(sendchan chan<- []byte) {
	seq := 1
	for true {
		time.Sleep(3 * time.Second)
		p := Post{
			Id:     fmt.Sprintf("%d", seq),
			Body:   fmt.Sprintf("My message %d", seq),
			Author: fmt.Sprintf("some-sha-%d", seq),
		}

		if seq%3 == 0 {
			p.Children = &[]Post{
				Post{
					Id:     fmt.Sprintf("%d", seq),
					Body:   fmt.Sprintf("you are basically wrong %d", seq),
					Author: fmt.Sprintf("jerk-%d", seq),
				},
			}
		}

		data, err := json.Marshal(p)

		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}

		sendchan <- []byte(data)
		seq += 1
	}
}
