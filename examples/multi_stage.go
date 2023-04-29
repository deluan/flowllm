package main

import (
	"context"
	"fmt"

	. "github.com/deluan/pipelm"
	"github.com/deluan/pipelm/llms/openai"
)

func main() {
	chain := Chain(
		ParallelChain(
			2,
			Chain(
				Template("What is a good name for a company that makes {product}?"),
				LLM(openai.NewCompletionModel(openai.Options{Model: "text-davinci-003", Temperature: 1})),
				MapOutputTo("name"),
			),
			Chain(
				//Template("What is a good slogan for a company that makes {product}?"),
				//LLM(openai.NewChatModel(openai.Options{Model: "gpt-3.5-turbo", Temperature: 1})),
				ChatTemplate{UserMessage("What is a good slogan for a company that makes {product}?")},
				ChatLLM(openai.NewChatModel(openai.Options{Model: "gpt-3.5-turbo", Temperature: 1})),
				MapOutputTo("slogan"),
			),
		),
		TrimSpace("name", "slogan"),
		TrimSuffix(".", "name"),
		Template("The company {name} makes {product} and their slogan is {slogan}."),
	)

	res, err := chain(context.Background(), Values{"product": "colorful sockets"})
	fmt.Println(res, err)
}
