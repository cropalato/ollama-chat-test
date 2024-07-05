//
// main.go
// Copyright (C) 2024 rmelo <Ricardo Melo <rmelo@ludia.com>>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
  "regexp"

	"github.com/ollama/ollama/api"
)

var msgResp string

func isInputFromPipe() bool {
    fileInfo, _ := os.Stdin.Stat()
    return fileInfo.Mode() & os.ModeCharDevice == 0
}

func chatResp(resp api.ChatResponse) error {
    msgResp += resp.Message.Content
    return nil
  }

func main() {
  client, err := api.ClientFromEnvironment()
  if err != nil {
    log.Fatal(err)
  }

  ctx := context.Background()

  if ! isInputFromPipe() {
    fmt.Printf("This command look for a pipe input.\n")
    m, err := client.List(ctx)
    if err != nil {
      log.Fatal(err)
    }
    for _, model := range m.Models {
      fmt.Printf("Model: %s\n", model.Model)
    }
    os.Exit(0)
  }


  stdin, err := io.ReadAll(os.Stdin)
  if err != nil {
    log.Fatal(err)
  }
  gitDiff := string(stdin)
  
  req := &api.ChatRequest{
    //Model: "mistral",
    Model: "deepseek-coder-v2",
    Messages: []api.Message{
      api.Message{
        Role: "user",
        Content: "Use the following git diff output to create a short version of a git conventional commit message:\n" + gitDiff,
      },
    },
  }
    err = client.Chat(ctx, req, chatResp)
    if err != nil {
      log.Fatal(err)
    }
    re := regexp.MustCompile(`.*\x60\x60\x60\s(?P<msg>.+)\s\x60\x60\x60.*`)
    matches := re.FindStringSubmatch(msgResp)
    if matches != nil {
      fmt.Println(matches[1])
    }
}
