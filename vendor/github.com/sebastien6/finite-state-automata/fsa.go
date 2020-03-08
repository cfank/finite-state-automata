package fsa

import (
	fmt "fmt"
	errors "errors"
)

var (
	errNoAlphabetDefined error = errors.New("no alphabet defined, use AddAlphabet method to define at least one acceptable event trigger")
	errNoStatesDefined error = errors.New("no state defined, use AddState method to define a minimum of two states")
	errNoTransitionDefined error = errors.New("no transition rules defined, use AddTransition method to define at least one valid transition rule")
	errNoFinalState error = errors.New("transitions did not end in a final state")
	fmtErrInvalidAlphabet string = "Invalid character within input string, %s is not listed as part of the allowing alphabet"
	fmtErrNoTransitionRule string = "Invalid, no transition rules defined for event %s from state %s"
)

// State container to state declaration
type State struct {
	Name  string
	Value interface{}
	FinalState bool
}

// transition container to state transition rule
type Transition struct {
	Source string
	Destination string
	Trigger string
}

// FSA Finiate State Automata access layer
type FSA interface {
	AddAlphabet(values ...string)
	AddState(initialState string, states ...State)
	AddTransition(transitionRules ...Transition)
	Event(event string) error
	Validate() error
	FinalState() (string, error)
}

func NewFSA() FSA {
	return  &fsa{
		alphabet: make(map[string]struct{}),
		states: make(map[string]State),
		transitions: make(map[string]string),
	}
}

// fsa Finite State Automata data container
type fsa struct {
	alphabet map[string]struct{}
	states map[string]State
	transitions map[string]string
	currentState string
	previousState string
	transitionLog []string
}

func (f *fsa) AddAlphabet(values ...string) {
	for _, v := range values {
		f.alphabet[v] = struct{}{}
	}
}

// AddState add a new state to the Finite State Automata (FMA) with an associated output value
func (f *fsa) AddState(initialState string, states ...State) {
	for _, v := range states {
		f.states[v.Name] = v	
	}
	f.currentState = initialState
}

// AddTransition add a new transition rule to the Finite State Automata (FMA) for state transition
func (f *fsa) AddTransition(transitionRules ...Transition) {
	for _, v := range transitionRules {
		f.transitions[v.Source+":"+v.Trigger] = v.Destination
	}
}

// setState update previous and current state
func (f *fsa) setState(stateName string) {
	f.previousState = f.currentState
	f.currentState = stateName
	f.transitionLog = append(f.transitionLog, fmt.Sprintf("%s -> %s", f.previousState, f.currentState))
}

func (f *fsa) Event(event string) error {
	if _, ok := f.alphabet[event]; !ok {
		return fmt.Errorf(fmtErrInvalidAlphabet, event)
	}

	key := f.currentState+":"+event
	if _, ok := f.transitions[key]; !ok {
		return fmt.Errorf(fmtErrNoTransitionRule, event, f.currentState)
	}
	
	f.setState(f.transitions[key])
	return nil
}

func (f *fsa) Validate() error {
	if len(f.alphabet) == 0 {
		return errNoAlphabetDefined
	}

	if len(f.states) < 2 {
		return errNoStatesDefined
	}

	if len(f.transitions) < 1 {
		return errNoTransitionDefined
	}
	return nil
}

func (f *fsa) FinalState() (string, error) {
	if !f.states[f.currentState].FinalState {
		return "", errNoFinalState
	} 
	return fmt.Sprintf("%s=%v", f.states[f.currentState].Name, f.states[f.currentState].Value), nil
}