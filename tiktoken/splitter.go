package tiktoken

import "github.com/deluan/pipelm"

func Splitter(model string, options pipelm.SplitterOptions) pipelm.Splitter {
	lenFunc := Len(model)
	options.LenFunc = lenFunc
	options.Separators = []string{"\n"}
	return pipelm.RecursiveTextSplitter(options)
}
