package flowllm

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
	vals := Values{}.Merge(values...)

	replaced := regexTemplate.ReplaceAllStringFunc(string(t), func(match string) string {
		variableName := match[1 : len(match)-1]
		variableValue, ok := vals[variableName]
		if ok {
			return fmt.Sprintf("%v", variableValue)
		}
		return match
	})

	vals[DefaultKey] = replaced
	return vals, nil
}

const chatHistoryPlaceholderRole = "_messages_placeholder"

// ChatTemplate is a prompt that can be used with Chat-style LLMs.
// It will format a list of messages, each with a role and a prompt.
type ChatTemplate []MessageTemplate

func (t ChatTemplate) Call(_ context.Context, values ...Values) (Values, error) {
	vals := Values{}.Merge(values...)

	var output strings.Builder
	msgs := t.messages(vals)
	for _, m := range msgs {
		output.WriteString(fmt.Sprintf("%s: %s\n", m.Role, m.Content))
	}
	vals[DefaultKey] = output.String()
	vals[DefaultChatKey] = msgs
	return vals, nil
}

func (t ChatTemplate) messages(values Values) ChatMessages {
	var res ChatMessages
	for _, m := range t {
		if m.Role == chatHistoryPlaceholderRole {
			chatHistory, _ := values[DefaultChatKey].(ChatMessages)
			res = append(res, chatHistory...)
			continue
		}
		content, _ := m.Template.Call(context.TODO(), values)
		res = append(res, ChatMessage{
			Role:    m.Role,
			Content: content.Get(DefaultKey),
		})
	}
	return res
}

// MessageTemplate is a prompt template that can be used with Chat-style LLMs. It is similar to Template,
// but it also specifies the role of the message.
type MessageTemplate struct {
	Template Handler
	Role     string
}

func SystemMessage(template string) MessageTemplate {
	return MessageTemplate{Template: Template(template), Role: "system"}
}

func UserMessage(template string) MessageTemplate {
	return MessageTemplate{Template: Template(template), Role: "user"}
}

func AssistantMessage(template string) MessageTemplate {
	return MessageTemplate{Template: Template(template), Role: "assistant"}
}

// MessageHistoryPlaceholder is a special message template that can be used with the WithMemory handler and
// Chat-style LLMs. It will be replaced with the history of messages in the conversation.
func MessageHistoryPlaceholder(variableName string) MessageTemplate {
	return MessageTemplate{Template: Template(variableName), Role: chatHistoryPlaceholderRole}
}
