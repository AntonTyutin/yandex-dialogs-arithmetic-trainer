package main

import (
	"fmt"

	"github.com/AntonTyutin/yandex-dialogs-arithmetic-trainer.git/config"
	helpers "github.com/AntonTyutin/yandex-dialogs-arithmetic-trainer.git/internal/infrastructure/yadialogs"
	"github.com/AntonTyutin/yandex-dialogs-arithmetic-trainer.git/pkg/yadialogs"
)

func YandexDialogsHandler(input *yadialogs.Input) (*yadialogs.Output, error) {
	context := helpers.CreateContext(input)
	message := helpers.CreateMessage(input)
	fmt.Printf("[before interpret] context: %+v; message: %+v\n", context, message)
	response := helpers.Interpret(config.Interpreters, message, context)
	fmt.Printf("[after interpret] context: %+v; message: %+v; response: %+v", context, message, response)
	return helpers.CreateResponse(response, context), nil
}
