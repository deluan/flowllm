package tiktoken

import "github.com/deluan/flowllm"

func Splitter(model string, options flowllm.SplitterOptions) flowllm.Splitter {
	lenFunc := Len(model)
	options.LenFunc = lenFunc
	options.Separators = []string{"\n"}
	return flowllm.RecursiveTextSplitter(options)
}
