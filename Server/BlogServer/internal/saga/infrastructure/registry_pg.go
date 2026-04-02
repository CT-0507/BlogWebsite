package infrastructure

import "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/domain"

type Registry struct {
	definitions map[string][]domain.Step
}

func NewRegistry() domain.Registry {
	return &Registry{
		definitions: make(map[string][]domain.Step),
	}
}

func (r *Registry) Register(name string, steps []domain.Step) {
	r.definitions[name] = steps
}

func (r *Registry) GetDefinition(name string) []domain.Step {
	return r.definitions[name]
}

func (r *Registry) GetStepByIndex(name string, index int32) *domain.Step {
	steps := r.definitions[name]
	if steps == nil {
		return nil
	}
	if index >= int32(len(steps)) {
		return nil
	}
	return &steps[index]
}

func (r *Registry) GetNextStep(name string, currentIndex int32) *domain.Step {
	steps := r.definitions[name]
	if steps == nil {
		return nil
	}
	if currentIndex+1 >= int32(len(steps)) {
		return nil
	}
	return &steps[currentIndex+1]
}
