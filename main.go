package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	rpcjson "github.com/gorilla/rpc/json"
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

type PostsListArgs struct {
}

type PostsListReply struct {
	Posts []Post
}

type PostsService struct{}

func (h *PostsService) List(r *http.Request, args *PostsListArgs, reply *PostsListReply) error {
	reply.Posts = []Post{
		Post{
			Id:     "my-id",
			Author: "someone",
			Body:   "my body goes here",
		},
		Post{
			Id:     "anotherId",
			Author: "someone else",
			Body:   "body text",
			Children: &[]Post{
				Post{
					Id:     "first-child-id",
					Author: "that first person",
					Body:   "1st child",
				},
				Post{
					Id:     "second-child-id",
					Author: "that first person",
					Body:   "2nd child",
				},
			},
		},
	}
	return nil
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

	s := rpc.NewServer()
	s.RegisterCodec(rpcjson.NewCodec(), "application/json")
	s.RegisterCodec(rpcjson.NewCodec(), "application/json;charset=UTF-8")
	s.RegisterService(new(PostsService), "posts")

	r := mux.NewRouter()
	r.HandleFunc("/", serveHome)
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	r.Handle("/rpc", s).Methods("POST")

	log.Printf("Running on localhost%s", *addr)
	err := http.ListenAndServe(*addr, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
