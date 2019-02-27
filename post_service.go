package main

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"go.cryptoscope.co/luigi"
	"go.cryptoscope.co/muxrpc"
	"go.cryptoscope.co/ssb"
	"go.cryptoscope.co/ssb/message"
)

type KeyValueRaw struct {
	Key       *ssb.MessageRef       `json:"key"`
	Value     message.LegacyMessage `json:"value"`
	Timestamp int64                 `json:"timestamp"`
}

type Post struct {
	Id, Body, Author string
	Children         *[]Post
}

type PostsListArgs struct {
}

type PostsListReply struct {
	Posts []Post
}

type PostsService struct{}

func (h *PostsService) List(r *http.Request, args *PostsListArgs, reply *PostsListReply) error {
	fmt.Println("List method")
	client, err := initClient("/home/justin/.ssb/secret") // this should be moved somewhere else
	if err != nil {
		// @@@ handle sending error reply
		fmt.Println(err)
		return err
	}

	fmt.Println("fetching from sbot")
	// the map thing tells the go client what the marshal format is.
	// src, err := client.Source(r.Context(), KeyValueRaw{}, muxrpc.Method{"createLogStream"})
	src, err := client.Source(r.Context(), KeyValueRaw{}, muxrpc.Method{"messagesByType"}, "post")
	if err != nil {
		fmt.Println(err)
		return err
	}
	// we could call an API that's 'get post by type'

	fmt.Println("iterating")
	v, err := src.Next(r.Context())
	if luigi.IsEOS(err) {
		fmt.Println("eos")
		return nil
	} else if err != nil {
		fmt.Println(err)
		return errors.Wrapf(err, "createLogStream: failed to drain")
	}
	fmt.Printf("V: %#v\n\n", v)
	msg := v.(KeyValueRaw)

	fmt.Printf("Message: %#v\n\n", msg)
	fmt.Printf("Author: %s\n\n", msg.Value.Author)
	fmt.Printf("User's Timestamp: %d\n\n", msg.Value.Timestamp)
	// lol @ content hash

	content, ok := msg.Value.Content.(map[string]interface{})
	if !ok {
		fmt.Println("error swapping to content hash")
		return fmt.Errorf("couldn't swap it to a content hash: %T", msg.Value.Content)
	}

	ctype, ok := content["type"]
	if !ok {
		fmt.Println("error type param")
		return errors.New("Message didn't have a type param")
	}
	ctext, ok := content["text"]
	if !ok {
		fmt.Println("error text param")
		return errors.New("Message didn't have text")
	}
	// croot, ok := content["root"]
	// if !ok {
	// 	return errors.New("Message didn't have a root param")
	// }

	switch ctype.(string) {
	case "about":
		fmt.Println("It's an about stanza")
	case "post":
		reply.Posts = []Post{
			Post{
				Id:     msg.Key.Ref(),
				Author: msg.Value.Author,
				Body:   ctext.(string),
			},
		}
		// fmt.Println("It's an post stanza")
	default:
		fmt.Println("Unknown")
	}

	// reply.Posts = []Post{
	// 	Post{
	// 		Id:     "my-id",
	// 		Author: "someone",
	// 		Body:   "my body goes here",
	// 	},
	// 	Post{
	// 		Id:     "anotherId",
	// 		Author: "someone else",
	// 		Body:   "body text",
	// 		Children: &[]Post{
	// 			Post{
	// 				Id:     "first-child-id",
	// 				Author: "that first person",
	// 				Body:   "1st child",
	// 			},
	// 			Post{
	// 				Id:     "second-child-id",
	// 				Author: "that first person",
	// 				Body:   "2nd child",
	// 			},
	// 		},
	// 	},
	// }
	return nil
}
