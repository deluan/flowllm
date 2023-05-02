package main

import (
	"context"
	"fmt"

	. "github.com/deluan/flowllm"
	"github.com/deluan/flowllm/llms/openai"
)

func init() {
	registerExample("multi_stage", "A multi-stage pipeline, showcasing sequential and parallel chains", multiStage)
}

func multiStage() {
	// Build a chain that will generate a company name and slogan, and then use them to generate a sentence.
	// Calls to the OpenAI API are made in parallel, and the results are merged into a single result.
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
		// You can modify the LLMs outputs using some string transformation handlers
		TrimSpace("name", "slogan"),
		TrimSuffix(".", "name"),
		Template("The company {name} makes {product} and their slogan is {slogan}."),
	)

	// Run the chain
	res, err := chain(context.Background(), Values{"product": "colorful sockets"})
	fmt.Println(res, err)
}
