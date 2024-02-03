package domain

func NewMultiplicationSampleProvider() *MultiplicationSampleProvider {
	return &MultiplicationSampleProvider{}
}

type MultiplicationSampleProvider struct{}

func (p *MultiplicationSampleProvider) Get(idx int) Sample {
	return Sample{idx/10 + 1, idx%10 + 1}
}