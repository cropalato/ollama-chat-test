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

	"github.com/ollama/ollama/api"
)

func isInputFromPipe() bool {
    fileInfo, _ := os.Stdin.Stat()
    return fileInfo.Mode() & os.ModeCharDevice == 0
}

func chatResp(resp api.ChatResponse) error {
    fmt.Print(resp.Message.Content)
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
        Content: "create a short git convetional commit message for the following changes\n" + gitDiff,
      },
    },
  }

  err = client.Chat(ctx, req, chatResp)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("\n")

}
