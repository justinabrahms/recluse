package main

import (
	"fmt"
	"net/http"
	"os/user"

	"github.com/pkg/errors"
	"go.cryptoscope.co/luigi"
	"go.cryptoscope.co/muxrpc"
	"go.cryptoscope.co/ssb"
	"go.cryptoscope.co/ssb/message"
)

const MAX_POSTS_IN_REPLY = 20

var client SbotClient

func init() {
	var err error
	client, err = initClient(fmt.Sprintf("/home/%s/.ssb/secret", user.Current()))
	if err != nil {
		panic(errors.Wrap(err, "This may mean you don't have an sbot running?"))
	}
}

type KeyValueRaw struct {
	Key       *ssb.MessageRef       `json:"key"`
	Value     message.LegacyMessage `json:"value"`
	Timestamp float64               `json:"timestamp"`
}

type Post struct {
	Id, Body, Author string
	Children         *[]Post
}

type PostsListArgs struct {
	Count int
}

type PostsListReply struct {
	Posts []Post
}

type PostsService struct{}

func (h *PostsService) List(r *http.Request, args *PostsListArgs, reply *PostsListReply) error {
	var count int
	if args.Count > MAX_POSTS_IN_REPLY || args.Count == 0 {
		count = MAX_POSTS_IN_REPLY
	} else {
		count = args.Count
	}

	// src, err := client.Source(r.Context(), KeyValueRaw{}, muxrpc.Method{"createLogStream"})
	src, err := client.Source(r.Context(), KeyValueRaw{}, muxrpc.Method{"messagesByType"}, "post")
	if err != nil {
		fmt.Println(err)
		return err
	}

	posts := []Post{}
	for len(posts) < count {
		v, err := src.Next(r.Context())
		if luigi.IsEOS(err) {
			fmt.Println("eos")
			break
		} else if err != nil {
			fmt.Println(err)
			return errors.Wrapf(err, "createLogStream: failed to drain")
		}

		msg := v.(KeyValueRaw)

		fmt.Printf("Message: %#v\n\n", msg)
		fmt.Printf("Author: %s\n\n", msg.Value.Author)
		fmt.Printf("User's Timestamp: %d\n\n", msg.Value.Timestamp)

		content, ok := msg.Value.Content.(map[string]interface{})
		if !ok {
			fmt.Println("error swapping to content hash bc of type: %T", msg.Value.Content)
			continue
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
			posts = append(posts, Post{
				Id:     msg.Key.Ref(),
				Author: msg.Value.Author,
				Body:   ctext.(string),
			})
		default:
			fmt.Println("Unknown")
		}
	}
	reply.Posts = posts

	return nil
}
