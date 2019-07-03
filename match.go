package go_cypherdsl

import "errors"

type MatchBuilder struct {
	//match steps are a linked list, store the first node for further iteration
	firstStep *matchStep

	//current step so we can append easily
	currentStep *matchStep

	//save all of the errors for the end
	errors      []error
}

type VertexStep struct {
	builder *MatchBuilder
}

type EdgeStep struct {
	builder *MatchBuilder
}

type PStep struct{
	builder *MatchBuilder
}

type matchStep struct {
	Vertices []V
	Edge     *E

	P bool

	OtherOperation string

	Next *matchStep
}

func (m *matchStep) ToCypher() (string, error){

}

func (m *MatchBuilder) ToCypher() (string, error){
	if m.firstStep == nil{
		return "", errors.New("no steps process")
	}

	if len(m.errors) != 0{
		errStr := ""
		for _, err := range m.errors{
			errStr += ";" + err.Error()
		}

		return "", errors.New("incurred one or many errors: " + errStr)
	}



	return "", nil //todo this function
}

func (m *MatchBuilder) P() *PStep{
	newStep := &matchStep{
		P: true,
	}

	//its the first step
	if m.currentStep == nil{
		m.firstStep = newStep
		m.currentStep = newStep
	} else {
		m.currentStep.Next = newStep
		m.currentStep = newStep
	}

	return &PStep{
		builder: m,
	}
}

func (p *PStep) V(vertices ...V) *MatchBuilder {
	if vertices == nil || len(vertices) == 0 {
		if p.builder.errors == nil{
			p.builder.errors = []error{}
		}
		p.builder.errors = append(p.builder.errors, errors.New("vertices can not be nil or empty"))
	}

	newStep := &matchStep{
		Vertices: vertices,
	}

	p.builder.currentStep.Next = newStep
	p.builder.currentStep = newStep

	return p.builder
}

func (m *MatchBuilder) V(vertices ...V) *VertexStep {
	if vertices == nil || len(vertices) == 0{
		if m.errors == nil{
			m.errors = []error{}
		}
		m.errors = append(m.errors, errors.New("vertices can not be nil or empty"))
	}

	newStep := &matchStep{
		Vertices: vertices,
	}

	//its the first step
	if m.currentStep == nil{
		m.firstStep = newStep
		m.currentStep = newStep
	} else {
		m.currentStep.Next = newStep
		m.currentStep = newStep
	}

	return &VertexStep{
		builder: m,
	}
}


func (e *EdgeStep) V(vertices ...V) *MatchBuilder{
	if vertices == nil || len(vertices) == 0{
		if e.builder.errors == nil{
			e.builder.errors = []error{}
		}
		e.builder.errors = append(e.builder.errors, errors.New("vertices can not be nil or empty"))
	}

	newStep := &matchStep{
		Vertices: vertices,
	}

	e.builder.currentStep.Next = newStep
	e.builder.currentStep = newStep

	return e.builder
}


func (v *VertexStep) E(edge E) *EdgeStep {
	newStep := &matchStep{
		Edge: &edge,
	}

	v.builder.currentStep.Next = newStep
	v.builder.currentStep = newStep

	return &EdgeStep{
		builder: v.builder,
	}
}