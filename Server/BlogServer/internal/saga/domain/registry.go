package domain

type Registry interface {
	Register(name string, steps []Step)
	GetDefinition(name string) []Step
	GetStepByIndex(name string, index int32) *Step
	GetNextStep(name string, currentIndex int32) *Step
}
