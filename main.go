package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	rpcjson "github.com/gorilla/rpc/json"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()

	go newMessagePublisher(hub.broadcast)

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
