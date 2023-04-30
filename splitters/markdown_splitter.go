package splitters

import "github.com/deluan/pipelm"

func Markdown(opts ...Option) pipelm.Splitter {
	separators := []string{
		// First, try to split along Markdown headings (starting with level 2)
		"\n## ",
		"\n### ",
		"\n#### ",
		"\n##### ",
		"\n###### ",
		// Note the alternative syntax for headings (below) is not handled here
		// Heading level 2
		// ---------------
		// End of code block
		"```\n\n",
		// Horizontal lines
		"\n\n***\n\n",
		"\n\n---\n\n",
		"\n\n___\n\n",
		// Note that this splitter doesn't handle horizontal lines defined
		// by *three or more* of ***, ---, or ___, but this is not handled
		"\n\n",
		"\n",
		" ",
		"",
	}
	opts = append([]Option{WithSeparators(separators...)}, opts...)
	return RecursiveCharacterText(opts...)
}
