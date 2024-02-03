package dialog

import (
	"github.com/AntonTyutin/yandex-dialogs-arithmetic-trainer.git/pkg/utilities"
)

var greetings = []*Response{
	{Text: "Давай потренируем устный счёт?"},
}

var StartDialog = Interpreter{
	Recognize: func(message *Message, context *Context) bool {
		return context.State.IsRoot()
	},
	Interpret: func(message *Message, context *Context) *Response {
		reaction := utilities.Random(greetings)
		greet := MultiplicationGreeter(context)

		return reaction.Merge(greet)
	},
}
