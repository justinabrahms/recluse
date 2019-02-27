package main

import "net/http"

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
