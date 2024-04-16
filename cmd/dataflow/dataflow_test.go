package dataflow

import (
	"errors"
	"fmt"
	"testing"
)

type EmitString struct {
	In  <-chan string
	Out *FanOut[string]
}

func NewEmitString(In chan string) *EmitString {
	return &EmitString{
		In:  In,
		Out: NewFanOut[string](),
	}
}

func (a *EmitString) Run() {
	go a.Out.Run()
	for v := range a.In {
		a.Out.In <- v
	}
	a.Out.Close()
}

func (a *AppendString) Input(name string) (interface{}, error) {
	return nil, errors.New("no input evailable")
}

func (a *AppendString) Output(name string) (interface{}, error) {
	switch name {
	case "Out":
		return a.Out.New(), nil
	default:
		return nil, errors.New("unknown output")
	}
}

type AppendString struct {
	str string
	In  chan string
	Out *FanOut[string]
}

func NewAppendString(str string) *AppendString {
	return &AppendString{
		str: str,
		In:  make(chan string),
		Out: NewFanOut[string](),
	}
}

func (a *AppendString) Run() {
	go a.Out.Run()
	for v := range a.In {
		a.Out.In <- v + a.str
	}

	a.Out.Close()
}

func (a *AppendString) Input(name string) (interface{}, error) {
	switch name {
	case "In":
		return a.In, nil
	default:
		return nil, errors.New("unknown input")
	}
}

func (a *AppendString) Output(name string) (interface{}, error) {
	switch name {
	case "Out":
		return a.Out.New(), nil
	default:
		return nil, errors.New("unknown output")
	}
}

type PrintString struct {
	In chan string
}

func NewPrintString() *PrintString {
	return &PrintString{
		In: make(chan string),
	}
}

func (p *PrintString) Run() {
	for v := range p.In {
		fmt.Println(v)
	}
}

func (p *PrintString) Input(name string) (interface{}, error) {
	switch name {
	case "In":
		return p.In, nil
	default:
		return nil, errors.New("unknown input")
	}
}

func (p *PrintString) Output(name string) (interface{}, error) {
	return nil, errors.New("unknown output")
}

func TestGraph(t *testing.T) {
	g := NewGraph()

	g.Add("append", NewAppendString(" World!"))
	g.Add("print", NewPrintString())

	if err := g.Connect("append", "Out", "print", "In"); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}

	if err := g.Send("append", "In", "hello"); err != nil {
		t.Fatalf("Send() error = %v", err)
	}
}
