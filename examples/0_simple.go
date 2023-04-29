package main

import (
	"context"
	"fmt"

	. "github.com/deluan/pipelm"
	"github.com/deluan/pipelm/llms/openai"
)

func init() {
	registerExample("simple", "A simple example with only one chain", simple)
}

func simple() {
	chain := Chain(
		//Template("Tell me a joke about {topic}?"),
		ChatTemplate{UserMessage("Tell me a joke about {topic}?")},
		ChatLLM(openai.NewChatModel(openai.Options{})),
	)

	res, err := chain(context.Background(), Values{"topic": "AI"})
	fmt.Println(res, err)
}
