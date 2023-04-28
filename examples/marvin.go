package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	. "github.com/deluan/pipelm"
	"github.com/deluan/pipelm/llms/openai"
)

var systemMessage = SystemMessageTemplate(`You are Marvin, the depressed Android from the Hitchhiker's Guide to the Galaxy.`)

func main() {
	ctx := context.Background()

	chain := Chain(
		ChatTemplate{
			systemMessage,
			UserMessageTemplate("{input}"),
		},
		ChatLLM(openai.NewChatModel(openai.Options{Model: "gpt-3.5-turbo", Temperature: 0.8})),
	)
	input := "Introduce yourself."
	reader := bufio.NewReader(os.Stdin)
	for {
		res, err := chain(ctx, Values{"input": input})
		if err != nil {
			panic(err)
		}
		fmt.Println("🤖", strings.TrimSpace(res.String()))

		print("❓  ")
		input, _ = reader.ReadString('\n')

	}
}
