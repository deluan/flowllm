package main

import (
	"context"
	"fmt"

	. "github.com/deluan/pipelm"
	"github.com/deluan/pipelm/llms/openai"
)

func main() {
	chain := Chain(
		ChatTemplate{UserMessageTemplate("Tell me a joke about {topic}?")},
		ChatLLM(openai.NewChatModel(openai.Options{})),
	)

	res, err := chain(context.Background(), Values{"topic": "AI"})
	fmt.Println(res, err)
}
