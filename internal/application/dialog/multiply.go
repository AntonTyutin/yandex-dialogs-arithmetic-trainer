package dialog

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"

	"github.com/AntonTyutin/yandex-dialogs-arithmetic-trainer.git/internal/domain"
	"github.com/AntonTyutin/yandex-dialogs-arithmetic-trainer.git/pkg/utilities"
)

const MultiplicationEntryState ContextState = "multiplication:entry"
const MultiplicationTrainingState ContextState = "multiplication:training"

var sampleProvider = domain.MultiplicationSampleProvider{}

func MultiplicationGreeter(context *Context) *Response {
	context.SwitchState(MultiplicationEntryState)
	return &Response{
		Text: `Пройдем по таблице умножения. Я задаю пример, а ты говори ответ!

			Если не знаешь ответ, можешь сказать "не знаю" или "сдаюсь". Также можно попросить меня повторить вопрос.

			Начнем?
		`,
	}
}

var askingForRepeatPattern = regexp.MustCompile(`(?i)(повтори|еще раз(ок)?)|как ты сказала|не понял(а)?`)

var RepeatIntro = Interpreter{
	Recognize: func(message *Message, context *Context) bool {
		return context.State.Is(MultiplicationEntryState) && message.MatchPattern(askingForRepeatPattern)
	},
	Interpret: func(message *Message, context *Context) *Response {
		return MultiplicationGreeter(context)
	},
}

var RepeatQuestion = Interpreter{
	Recognize: func(message *Message, context *Context) bool {
		return context.State.Is(MultiplicationTrainingState) && message.MatchPattern(askingForRepeatPattern)
	},
	Interpret: func(message *Message, context *Context) *Response {
		previousSampleIdx := int(context.Data["sampleIdx"].(float64))
		sample := sampleProvider.Get(previousSampleIdx)
		return &Response{
			Text:   fmt.Sprintf("Посчитай %d × %d =", sample.Operand1, sample.Operand2),
			Speach: fmt.Sprintf("Сколько будет %d умножить на %d?", sample.Operand1, sample.Operand2),
		}
	},
}

var helpMultiplicationPattern = regexp.MustCompile(`(?i)^(помощь|справка|(расскажи )?что ты умеешь)$`)

var MultiplicationHelp = Interpreter{
	Recognize: func(message *Message, context *Context) bool {
		// (context.State.Is(MultiplicationEntryState) || context.State.Is(MultiplicationTrainingState)) &&
		return message.MatchPattern(helpMultiplicationPattern)
	},
	Interpret: func(message *Message, context *Context) *Response {
		return &Response{
			Text: `Мы находимся в режиме тренировки таблицы умножения, когда я задаю примеры, а ты говоришь ответы.

				Пока это единственный режим в навыке "Тренировка арифметики".

				Если не знаешь ответ, можешь сказать "не знаю" или "сдаюсь". Также можно попросить меня повторить вопрос.
			`,
		}
	},
}

var readyToStartPattern = regexp.MustCompile(`(?i)^(ну )?(да|(да )?(давай( нач(инай|нем|инаем))?|нач(нем|инаем)|погнали|поехали|окей))$`)
var letsStart = []string{"Начнем!", "Поехали!", "Погнали!"}

var ReadyToTrainMultiplication = Interpreter{
	Recognize: func(message *Message, context *Context) bool {
		return context.State.Is(MultiplicationEntryState) && message.MatchPattern(readyToStartPattern)
	},
	Interpret: func(message *Message, context *Context) *Response {
		response := &Response{Text: utilities.Random(letsStart)}
		sampleIdx := getRandomSampleIdx()
		context.SwitchState(MultiplicationTrainingState)
		context.Data["sampleIdx"] = sampleIdx
		sample := sampleProvider.Get(sampleIdx)
		return response.Merge(
			&Response{
				Text:   fmt.Sprintf("Посчитай %d × %d =", sample.Operand1, sample.Operand2),
				Speach: fmt.Sprintf("Сколько будет %d умножить на %d?", sample.Operand1, sample.Operand2),
			},
		)
	},
}

var multiplicationAnswerPattern = regexp.MustCompile(`(?i)(?P<result>\d+)((\D+(?P<another_result>\d+))?\D+(?P<too_much_results>\d+))?|(?P<dont_know>не знаю|сдаюсь|дальше)`)
var goodAnswerVariants = []string{"Молодец!", "Правильно!", "Верно!", "Ага.", "Хорошо.", "Так держать!"}
var wrongAnswerVariants = []string{"Не-а!", "Не верно!", "Нет!"}
var sayCorrectAnswerVariants = []string{"Правильный ответ – %d!", "Это будет %d.", "Думаю, это %d."}

var CheckMultiplicationResultAndAskAgain = Interpreter{
	Recognize: func(action *Message, context *Context) bool {
		return context.State.Is(MultiplicationTrainingState) && action.MatchPattern(multiplicationAnswerPattern)
	},
	Interpret: func(message *Message, context *Context) *Response {
		previousSampleIdx := int(context.Data["sampleIdx"].(float64))
		if _, doesntKnow := message.Meanings["dont_know"]; doesntKnow {
			response := &Response{Text: fmt.Sprintf(utilities.Random(sayCorrectAnswerVariants), sayCorrect(sampleProvider.Get(previousSampleIdx)))}
			newSampleIdx := getRandomSampleIdx()
			context.Data["sampleIdx"] = newSampleIdx
			return response.Merge(getQuestion(sampleProvider.Get(newSampleIdx)))
		} else if result, ok := message.Meanings["result"]; ok {
			if anotherResult, ok := message.Meanings["another_result"]; ok {
				if _, tooMuch := message.Meanings["too_much_results"]; tooMuch {
					return &Response{Text: "Слишком много вариантов! Выбери один."}
				} else {
					return &Response{Text: fmt.Sprintf("Так %s или %s?", result, anotherResult)}
				}
			}

			answer, _ := strconv.Atoi(result)
			var response *Response
			if isCorrect(sampleProvider.Get(previousSampleIdx), answer) {
				response = &Response{Text: utilities.Random(goodAnswerVariants)}
			} else {
				response = (&Response{Text: utilities.Random(wrongAnswerVariants)}).Merge(
					&Response{Text: fmt.Sprintf(utilities.Random(sayCorrectAnswerVariants), sayCorrect(sampleProvider.Get(previousSampleIdx)))},
				)
			}
			newSampleIdx := getRandomSampleIdx()
			context.Data["sampleIdx"] = newSampleIdx

			return response.Merge(getQuestion(sampleProvider.Get(newSampleIdx)))
		}

		reaction := &Response{Text: "Не поняла ответ."}
		sample := sampleProvider.Get(previousSampleIdx)
		return reaction.Merge(
			&Response{
				Text:   fmt.Sprintf("Еще раз... %d × %d =", sample.Operand1, sample.Operand2),
				Speach: fmt.Sprintf("Так сколько будет %d умножить на %d?", sample.Operand1, sample.Operand2),
			},
		)
	},
}

var firstOperand = []string{"Од+инажды", "Дв+ажды", "Тр+ижды", "Чет+ырежды", "П+ятью", "Ш+естью", "С+емью", "В+осемью", "Д+евятью", "Д+есятью"}

func getQuestion(s domain.Sample) *Response {
	if rand.Intn(2) == 0 {
		return &Response{
			Text:   fmt.Sprintf("%d × %d =", s.Operand1, s.Operand2),
			Speach: fmt.Sprintf("%s %d", firstOperand[s.Operand1-1], s.Operand2),
		}
	}

	return &Response{
		Text:   fmt.Sprintf("%d × %d =", s.Operand1, s.Operand2),
		Speach: fmt.Sprintf("%d на %d", s.Operand1, s.Operand2),
	}
}

func getRandomSampleIdx() int {
	return (rand.Intn(8)+1)*10 + (rand.Intn(8) + 1)
}

func isCorrect(s domain.Sample, result int) bool {
	return sayCorrect(s) == result
}
func sayCorrect(s domain.Sample) int {
	return s.Operand1 * s.Operand2
}
