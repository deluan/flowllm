package pipelm

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

// Template can be used to format a string with variables. Useful to create prompts.
// It uses a simple template syntax, where variables are enclosed in curly braces.
type Template string

var regexTemplate = regexp.MustCompile(`{(\w+)}`)

func (t Template) Call(_ context.Context, values ...Values) (Values, error) {
	variables := Values{}.Merge(values...)

	replaced := regexTemplate.ReplaceAllStringFunc(string(t), func(match string) string {
		variableName := match[1 : len(match)-1]
		variableValue, ok := variables[variableName]
		if ok {
			return fmt.Sprintf("%v", variableValue)
		}
		return match
	})

	variables[DefaultKey] = replaced
	return variables, nil
}

// ChatTemplate is a prompt that can be used with ChatTemplate style LLMs.
// It will format a list of messages, each with a role and a prompt.
type ChatTemplate []MessageTemplate

func (t ChatTemplate) Call(ctx context.Context, values ...Values) (Values, error) {
	variables := Values{}.Merge(values...)

	var output strings.Builder
	for _, m := range t {
		msg, err := m.Template.Call(ctx, variables)
		if err != nil {
			return nil, err
		}
		output.WriteString(fmt.Sprintf("%s: %s\n", m.Role, msg.Get(DefaultKey)))
	}
	variables[DefaultKey] = output.String()
	return variables, nil
}

type chatPromptHandler interface {
	Messages(Values) []MessageTemplate
}

func (t ChatTemplate) Messages(values Values) []MessageTemplate {
	var res []MessageTemplate
	for _, m := range t {
		if p, ok := m.Template.(chatPromptHandler); ok {
			res = append(res, p.Messages(values)...)
		} else {
			res = append(res, m)
		}
	}
	return res
}

type MessageTemplate struct {
	Template Handler
	Role     string
}

func SystemMessageTemplate(template string) MessageTemplate {
	return MessageTemplate{Template: Template(template), Role: "system"}
}

func UserMessageTemplate(template string) MessageTemplate {
	return MessageTemplate{Template: Template(template), Role: "user"}
}

func AssistantMessageTemplate(template string) MessageTemplate {
	return MessageTemplate{Template: Template(template), Role: "assistant"}
}
