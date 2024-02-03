package domain

type Sample struct {
	Operand1 int
	Operand2 int
}

type SampleProvider interface {
	Get(int) Sample
}
