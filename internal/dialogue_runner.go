package internal

import (
	"embed"
	"errors"
	"fmt"
	"github.com/DrJosh9000/yarn"
	"github.com/DrJosh9000/yarn/bytecode"
	"github.com/Frabjous-Studios/ebitengine-game-template/internal/debug"
	"github.com/razor-1/localizer-cldr/resources/language"
	"strings"
	"sync"
)

const yarnFile = "gamedata/yarn/bin/game"

// yarnBin contains all yarn output from this compilation process.
//
//go:embed gamedata/yarn/bin
var yarnBin embed.FS

type RunnerState uint8

const (
	RunnerStopped RunnerState = iota // RunnerStopped
	RunnerRunning                    // RunnerRunning is set for a runner that's running.
	RunnerWaiting                    // RunnerWaiting indicates the runner is waiting for the player to select an option.
)

// DialogueRunner runs YarnSpinner and any commands from the script. It buffers lines delivered and handles blocking the
// YarnSpinner thread as expected by the game.
type DialogueRunner struct {
	program     *bytecode.Program
	stringTable *yarn.StringTable

	gameState    *GameState
	runState     RunnerState          // runState is manipulated by handler
	CurrNodeName string               // CurrNodeName is the name of the currently running node.
	vm           *yarn.VirtualMachine // vm is the Yarn virtual machine.
	handler      *blockingHandler

	mut *sync.RWMutex

	running bool
}

func NewDialogueRunner() (*DialogueRunner, error) {
	program, st, err := yarn.LoadFilesFS(yarnBin, yarnFile+".yarnc", language.EN_US)
	if err != nil {
		return nil, err
	}
	return &DialogueRunner{
		program:     program,
		stringTable: st,
		runState:    RunnerStopped,
		mut:         &sync.RWMutex{},
	}, nil
}

// Start starts the runner, which blocks the current thread until a fatal error occurs.
func (r *DialogueRunner) Start(state *GameState) error {
	if r.handler != nil {
		return errors.New("runner already running")
	}
	r.handler = new(blockingHandler)
	r.handler.runner = r
	r.handler.choices = make(chan int)
	r.runState = RunnerRunning
	r.gameState = state

	if state.Vars == nil {
		state.Vars = make(yarn.MapVariableStorage)
	}

	r.vm = &yarn.VirtualMachine{
		Program: r.program,
		Handler: r.handler,
		Vars:    state.Vars,
	}
	return r.vm.Run(state.CurrentNode)
}

func (r *DialogueRunner) GameState() *GameState {
	r.gameState.Vars = r.vm.Vars.(yarn.MapVariableStorage)
	r.gameState.CurrentNode = r.CurrNodeName
	return r.gameState
}

// Lines returns the set of lines currently displayed, or nil and an error if the runner has not yet started.
func (r *DialogueRunner) Lines() ([]*yarn.AttributedString, error) {
	r.mut.RLock()
	defer r.mut.RUnlock()
	if r.runState == RunnerStopped {
		return nil, nil
	}
	if len(r.handler.buffered) == 0 {
		return nil, nil
	}
	var result []*yarn.AttributedString
	for _, line := range r.handler.buffered {
		as, err := r.stringTable.Render(line)
		if err != nil {
			return nil, fmt.Errorf("error rendering line '%s': %w", line.ID, err)
		}
		result = append(result, as)
	}
	return result, nil
}

// Options returns the current set of options which dialogue runner is waiting on the player to choose from. Returns
// nil if the RunnerState is not RunnerWaiting
func (r *DialogueRunner) Options() ([]*yarn.AttributedString, error) {
	var err error
	r.mut.RLock()
	defer r.mut.RUnlock()
	if r.runState == RunnerStopped {
		return nil, nil
	}

	result := make([]*yarn.AttributedString, len(r.handler.options)) // defensive copy
	for i, opt := range r.handler.options {
		result[i], err = r.stringTable.Render(opt.Line)
		if err != nil {
			return nil, fmt.Errorf("error rendering option '%d': %w", opt.ID, err)
		}
	}
	return result, nil
}

// Choose selects the option from the list of current options.
func (r *DialogueRunner) Choose(choice int) error {
	r.mut.RLock()
	defer r.mut.RUnlock()
	if r.runState != RunnerWaiting {
		return errors.New("choose not called in the waiting runState")
	}
	r.handler.options = nil     // unset the options
	r.handler.choices <- choice // send and unblock the runner
	return nil
}

// Command runs a command from the game. Commands are all run on the runner thread.
func (r *DialogueRunner) Command(command string) error {
	command = strings.TrimSpace(command)
	tokens := strings.Split(command, " ")
	if len(tokens) == 0 {
		return fmt.Errorf("bad command: %s", command)
	}
	switch tokens[0] {
	default:
		return fmt.Errorf("unknown command %s", tokens[0])
	}
}

// blockingHandler implements yarn.DialogueHandler. It is tightly coupled with DialogueRunner (on purpose).
type blockingHandler struct {
	runner *DialogueRunner

	buffered []yarn.Line
	options  []yarn.Option

	choices chan int // choices is used to send choices back to the runner
}

func (r *blockingHandler) NodeStart(nodeName string) error {
	r.runner.mut.Lock()
	defer r.runner.mut.Unlock()

	r.runner.runState = RunnerRunning
	r.runner.CurrNodeName = nodeName
	debug.Printf("starting node: %v", nodeName)
	return nil
}

func (r *blockingHandler) PrepareForLines(_ []string) error {
	return nil
}

func (r *blockingHandler) Line(line yarn.Line) error {
	r.runner.mut.Lock()
	defer r.runner.mut.Unlock()

	r.buffered = append(r.buffered, line)
	debug.Printf("line: %s", line)
	return nil
}

func (r *blockingHandler) Options(options []yarn.Option) (int, error) {
	func() {
		r.runner.mut.Lock()
		defer r.runner.mut.Unlock()

		r.options = options
		r.runner.runState = RunnerWaiting
		debug.Printf("options: %v", options)
	}()

	choice := <-r.choices

	func() {
		r.runner.mut.Lock()
		defer r.runner.mut.Unlock()

		r.runner.runState = RunnerRunning
		r.buffered = nil // clear the buffered lines
	}()
	return choice, nil
}

func (r *blockingHandler) Command(command string) error {
	return r.runner.Command(command)
}

func (r *blockingHandler) NodeComplete(nodeName string) error {
	debug.Printf("stopping node: %v", nodeName)
	return nil
}

func (r *blockingHandler) DialogueComplete() error {
	debug.Println("dialogue complete")
	return nil
}
