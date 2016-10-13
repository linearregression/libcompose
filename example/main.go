package main

import (
	"log"

	"golang.org/x/net/context"

	"github.com/codeship/libcompose/docker"
	"github.com/codeship/libcompose/project"
	"github.com/codeship/libcompose/project/options"
)

func main() {
	project, err := docker.NewProject(&docker.Context{
		Context: project.Context{
			ComposeFiles: []string{"docker-compose.yml"},
			ProjectName:  "yeah-compose",
		},
	}, nil)

	if err != nil {
		log.Fatal(err)
	}

	err = project.Up(context.Background(), options.Up{})

	if err != nil {
		log.Fatal(err)
	}
}
