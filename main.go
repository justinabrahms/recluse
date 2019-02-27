package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var addr = flag.String("addr", ":8080", "http service address")

type Post struct {
	Id, Body, Author string
	Children         *[]Post
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()

	go func() {
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
			hub.broadcast <- []byte(data)
			seq += 1
		}
	}()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
