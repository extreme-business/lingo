package dataflow

import (
	"fmt"
	"reflect"
	"sync"
)

type Process interface {
	Run()
	Input(string) (interface{}, error)
	Output(string) (interface{}, error)
}

type Graph struct {
	processes map[string]Process
	wg        sync.WaitGroup
}

func NewGraph() *Graph {
	return &Graph{
		processes: make(map[string]Process),
		wg:        sync.WaitGroup{},
	}
}

func (g *Graph) Add(name string, process Process) {
	g.processes[name] = process
}

func (g *Graph) Run() {
	for _, process := range g.processes {
		go process.Run()
	}
}

// Connect attempts to connect the output of one process to the input of another within the graph.
func (g *Graph) Connect(from, fromPort, to, toPort string) error {
	fromProcess, ok := g.processes[from]
	if !ok {
		return fmt.Errorf("process %q not found", from)
	}

	toProcess, ok := g.processes[to]
	if !ok {
		return fmt.Errorf("process %q not found", to)
	}

	out, err := fromProcess.Output(fromPort)
	if err != nil {
		return fmt.Errorf("failed to get output from %q: %v", from, err)
	}

	in, err := toProcess.Input(toPort)
	if err != nil {
		return fmt.Errorf("failed to get input from %q: %v", to, err)
	}

	if reflect.TypeOf(out) != reflect.TypeOf(in) {
		fmt.Println("The types of the channels are the same")
	}

	inVal := reflect.ValueOf(in)
	outVal := reflect.ValueOf(out)

	// Ensure src and dst are both Channels
	if inVal.Kind() != reflect.Chan || outVal.Kind() != reflect.Chan {
		return fmt.Errorf("both src and dst must be channels")
	}

	// Check if the channels are of the same type
	if inVal.Type() != outVal.Type() {
		return fmt.Errorf("channels are of different types")
	}

	// Check if the src channel can receive and dst channel can send
	if inVal.Type().ChanDir()&reflect.RecvDir == 0 || outVal.Type().ChanDir()&reflect.SendDir == 0 {
		return fmt.Errorf("invalid channel directions")
	}

	// Perform the transfer in a goroutine to avoid blocking
	go func() {
		for {
			v, ok := outVal.Recv()
			if !ok {
				break
			}

			inVal.Send(v)
		}
	}()

	return nil
}

func (g *Graph) Send(name, port string, data interface{}) error {
	g.wg.Add(1)
	defer g.wg.Done()

	process, ok := g.processes[name]
	if !ok {
		return fmt.Errorf("process %q not found", name)
	}

	in, err := process.Input(port)
	if err != nil {
		return fmt.Errorf("failed to get input from %q: %v", name, err)
	}

	inVal := reflect.ValueOf(in)
	if inVal.Kind() != reflect.Chan {
		return fmt.Errorf("input %q is not a channel", port)
	}

	inVal.Send(reflect.ValueOf(data))
	return nil
}
