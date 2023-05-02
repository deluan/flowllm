# FlowLLM

[![Build](https://img.shields.io/github/actions/workflow/status/deluan/flowllm/go.yml?branch=main&logo=github)](https://github.com/deluan/flowllm/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/deluan/flowllm)](https://goreportcard.com/report/github.com/deluan/flowllm)
[![GoDoc](https://pkg.go.dev/badge/github.com/deluan/flowllm)](https://pkg.go.dev/github.com/deluan/flowllm)
[![License](https://img.shields.io/github/license/deluan/flowllm)](/LICENSE)

> **NOTICE**: This is still a work-in-progress. The interfaces provided by this project are still subject to change.

FlowLLM is a Go framework for developing applications that leverage the power of language models.
It uses composability and chain of responsibility patterns, and offers tools to access LLMs, build prompts, and
chain calls together. The package also includes parsers and database integrations, making it an ideal solution for
developers working with language models.

FlowLLM is heavily inspired by the [LangChain](https://docs.langchain.com/docs) project.

## Usage

In the example below we use FlowLLM to build a simple chain that generates a company name and slogan based on a
product name. The application uses two different LLMs, one for generating the company name and another for generating
the slogan. The slogan LLM is a chat model, so we use a different template to interact with it. the LLMs are called in
parallel, and the results are then combined into a single output.


```go
package main

import (
    "context"
    "fmt"

    . "github.com/deluan/flowllm"
    "github.com/deluan/flowllm/llms/openai"
)

func main() {
    // Build a chain that will generate a company name and slogan, and then use them 
    // to generate a sentence. Calls to the OpenAI API are made in parallel, and the 
    // results are merged into a single result.
    chain := Chain(
        ParallelChain(
            2,
            Chain(
                Template("What is a good name for a company that makes {product}?"),
                LLM(openai.NewCompletionModel(openai.Options{Model: "text-davinci-003", Temperature: 1})),
                MapOutputTo("name"),
            ),
            Chain(
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

    // Output:
    // The company Rainbow Socks Co makes colorful socks and their slogan is "Life is too short for boring socks – let us add some color to your steps!". <nil>
}

```

For more features and advanced usage, please check the [examples](/examples) folder.

## Installation

To install FlowLLM, use the following command:

```sh
go get -u github.com/deluan/flowllm
```

## Features
- Access to LLMs and their capabilities
- Tools to build prompts and parse outputs
- Database integrations for seamless data storage and retrieval
- Inspired by the langchain project, but striving to stay true to Go idioms and patterns

## Usage
For examples and detailed usage instructions,
please refer to the [documentation](https://pkg.go.dev/github.com/deluan/flowllm) (WIP).
Also check the [examples](/examples) folder.

## Contributing
We welcome contributions from the community!
Please read our [contributing guidelines](https://github.com/deluan/flowllm/blob/main/CONTRIBUTING.md)
for more information on how to get started.

## License
FlowLLM is released under the [MIT License](/LICENSE).
