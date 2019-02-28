package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
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

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	r.Handle("/rpc", s).Methods("POST")

	// This is a cors preflight request, which is sent when trying to invoke jsonrpc.
	r.HandleFunc("/rpc", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Println("handling options request")
		rw.Header().Set("Access-Control-Allow-Headers", "*")
		rw.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081") // this should be a flag
	}).Methods("OPTIONS")

	log.Printf("Running on localhost%s", *addr)
	corsHandler := handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"OPTIONS", "POST"}),
	)(r)

	err := http.ListenAndServe(*addr, handlers.LoggingHandler(os.Stdout, corsHandler))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
