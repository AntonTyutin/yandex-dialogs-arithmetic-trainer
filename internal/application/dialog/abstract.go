package dialog

import (
	"fmt"
	"regexp"

	"github.com/AntonTyutin/yandex-dialogs-arithmetic-trainer.git/pkg/utilities"
)

type Interpreter struct {
	Recognize func(*Message, *Context) bool
	Interpret func(*Message, *Context) *Response
}

type Greeter interface {
	Greet(*Context) *Response
}

const UndefinedState ContextState = ""
const RootState ContextState = "root"

type ContextState string

func (s ContextState) IsRoot() bool {
	return s == RootState
}
func (s ContextState) Is(state ContextState) bool {
	return s == state
}

func NewContext() *Context {
	return &Context{
		State: RootState,
	}
}

type Context struct {
	State         ContextState   `json:"state"`
	PreviousState ContextState   `json:"previous_state"`
	Data          map[string]any `json:"data"`
}

func (c *Context) SwitchState(state ContextState) *Context {
	c.PreviousState = c.State
	c.State = state
	return c
}

type Message struct {
	Text     string
	Meanings map[string]string
}

func (m *Message) MatchPattern(re *regexp.Regexp) bool {
	m.Meanings = nil
	regexResult := re.FindStringSubmatch(m.Text)

	if regexResult == nil {
		return false
	}

	keys := re.SubexpNames()
	meanings := make(map[string]string)
	for i, key := range keys {
		if key != "" && regexResult[i] != "" {
			meanings[key] = regexResult[i]
		}
	}
	if len(meanings) != 0 {
		m.Meanings = meanings
	}

	return true
}

type Response struct {
	Text   string
	Speach string
}

func (r *Response) Merge(b *Response) *Response {
	newResponse := &Response{}
	if r.Speach != "" || b.Speach != "" {
		newResponse.Speach = fmt.Sprintf("%s. %s", utilities.Coalesce(r.Speach, r.Text), utilities.Coalesce(b.Speach, b.Text))
	}
	newResponse.Text = fmt.Sprintf("%s\n\n%s", r.Text, b.Text)

	return newResponse
}
