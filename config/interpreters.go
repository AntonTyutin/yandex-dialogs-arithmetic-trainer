package config

import (
	"github.com/AntonTyutin/yandex-dialogs-arithmetic-trainer.git/internal/application/dialog"
)

var Interpreters = []dialog.Interpreter{
	dialog.StartDialog,
	dialog.ReadyToTrainMultiplication,
	dialog.CheckMultiplicationResultAndAskAgain,
	dialog.RepeatIntro,
	dialog.RepeatQuestion,
	dialog.MultiplicationHelp,
}
