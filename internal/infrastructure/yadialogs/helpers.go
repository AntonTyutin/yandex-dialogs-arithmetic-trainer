package yadialogs

import (
	"fmt"
	"maps"

	"github.com/AntonTyutin/yandex-dialogs-arithmetic-trainer.git/internal/application/dialog"
	"github.com/AntonTyutin/yandex-dialogs-arithmetic-trainer.git/pkg/yadialogs"
)

func CreateContext(input *yadialogs.Input) *dialog.Context {
	if input.Session.New {
		return &dialog.Context{
			State:         dialog.RootState,
			PreviousState: dialog.UndefinedState,
			Data:          map[string]any{},
		}
	}

	dialogState := input.State.Session["dialog_state"]
	if dialogState == nil {
		return &dialog.Context{
			State:         dialog.RootState,
			PreviousState: dialog.UndefinedState,
			Data:          map[string]any{},
		}
	}

	previousState := input.State.Session["dialog_previous_state"]
	switch previousState.(type) {
	case string:
	default:
		previousState = dialog.UndefinedState
	}

	var contextData map[string]any
	switch context := input.State.Session["context"].(type) {
	case map[string]any:
		contextData = maps.Clone(context)
	default:
		contextData = map[string]any{}
	}

	return &dialog.Context{
		State:         dialog.ContextState(dialogState.(string)),
		PreviousState: dialog.ContextState(previousState.(string)),
		Data:          contextData,
	}
}

func CreateMessage(input *yadialogs.Input) *dialog.Message {
	return &dialog.Message{Text: input.Request.Command}
}

func CreateResponse(response *dialog.Response, context *dialog.Context) *yadialogs.Output {
	sessionData := yadialogs.StateValues{}
	sessionData["dialog_state"] = string(context.State)
	sessionData["dialog_previous_state"] = string(context.PreviousState)
	sessionData["context"] = context.Data
	return &yadialogs.Output{
		Response: yadialogs.Response{
			Text: response.Text,
			Tts:  response.Speach,
		},
		SessionState: sessionData,
		Version:      "1.0",
	}
}

func Interpret(interpreters []dialog.Interpreter, message *dialog.Message, context *dialog.Context) *dialog.Response {
	interpret := makeInterpreter(interpreters)
	response := interpret(message, context)

	if response == nil {
		var unrecognizedPhrasesCount int
		switch count := context.Data["interpreter_unrecognized_phrases_count"].(type) {
		case int:
			unrecognizedPhrasesCount = count
		default:
			unrecognizedPhrasesCount = 0
		}

		response = &dialog.Response{
			Text: fmt.Sprintf("Прости, я не знаю, как понимать \"%s\".", message.Text),
		}
		if repeat := interpret(&dialog.Message{Text: "повтори"}, context); repeat != nil {
			response = response.Merge(repeat)
		}
		if unrecognizedPhrasesCount > 1 {
			if help := interpret(&dialog.Message{Text: "помощь"}, context); help != nil {
				response = response.Merge(help)
			}
		}

		unrecognizedPhrasesCount++

		context.Data["interpreter_unrecognized_phrases_count"] = unrecognizedPhrasesCount
	} else {
		delete(context.Data, "interpreter_unrecognized_phrases_count")
	}

	return response
}

func makeInterpreter(interpreters []dialog.Interpreter) func(*dialog.Message, *dialog.Context) *dialog.Response {
	return func(message *dialog.Message, context *dialog.Context) *dialog.Response {
		for _, interpreter := range interpreters {
			if interpreter.Recognize(message, context) {
				return interpreter.Interpret(message, context)
			}
		}

		return nil
	}
}
