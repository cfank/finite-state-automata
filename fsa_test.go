package fsa

import (
	testing "testing"
	reflect "reflect"
	fmt "fmt"

	assert "github.com/stretchr/testify/assert"
)

func TestAddAlphabet(t *testing.T) {
	expected := make(map[string]struct{})
	expected["0"] = struct{}{}
	expected["1"] = struct{}{}
	automata := fsa{
		alphabet: make(map[string]struct{}),
	}
	automata.AddAlphabet("0", "1")

	if eq := reflect.DeepEqual(expected, automata.alphabet); !eq {
		t.Error("invalid alphabet map")
	}
}

func TestAddState(t *testing.T) {
	expected := make(map[string]State)
	expected["S0"] = State{"S0", 0, true}
	expected["S1"] = State{"S1", 1, true}
	expected["S2"] = State{"S2", 2, true}
	automata := fsa{
		states: make(map[string]State),
	}
	automata.AddState("S0", State{"S0", 0, true}, State{"S1", 1, true}, State{"S2", 2, true})
	if eq := reflect.DeepEqual(expected, automata.states); !eq {
		t.Error("invalid states map")
	}
}

func TestAddTransition(t *testing.T) {
	expected := make(map[string]string)
	expected["S0:0"] = "S0"
	expected["S0:1"] = "S1"
	automata := fsa{
		transitions: make(map[string]string),
	}
	automata.AddTransition(
		Transition{"S0", "S0", "0"},
		Transition{"S0", "S1", "1"},
	)
	if eq := reflect.DeepEqual(expected, automata.transitions); !eq {
		t.Error("invalid transitions map")
	}
}

func TestEvent(t *testing.T) {
	automata := fsa{
		alphabet: make(map[string]struct{}),
		states: make(map[string]State),
		transitions: make(map[string]string),
	}
	automata.AddAlphabet("0", "1")
	automata.AddState("S0", State{"S0", 0, true}, State{"S1", 1, true})
	automata.AddTransition(
		Transition{Source: "S0", Destination: "S0", Trigger: "0"},
		Transition{Source: "S0", Destination: "S1", Trigger: "1"},
	)
	err := automata.Event("1")
	assert.NoError(t, err, "error should be nil")

	err = automata.Event("2")
	assert.EqualError(t, err, fmt.Sprintf(fmtErrInvalidAlphabet, "2"))

	err = automata.Event("1")
	assert.EqualError(t, err, fmt.Sprintf(fmtErrNoTransitionRule, "1", "S1"))
}

func TestValidate(t *testing.T) {
	automata := fsa{
		alphabet: make(map[string]struct{}),
		states: make(map[string]State),
		transitions: make(map[string]string),
	}
	err := automata.Validate()
	assert.EqualError(t, err, errNoAlphabetDefined.Error())

	automata.AddAlphabet("0", "1")
	err = automata.Validate()
	assert.EqualError(t, err, errNoStatesDefined.Error())

	automata.AddState("S0", State{"S0", 0, true}, State{"S1", 1, true})
	err = automata.Validate()
	assert.EqualError(t, err, errNoTransitionDefined.Error())

}

func TestFinalState(t *testing.T) {
	automata := fsa{
		alphabet: make(map[string]struct{}),
		states: make(map[string]State),
		transitions: make(map[string]string),
	}
	automata.AddAlphabet("0", "1")
	automata.AddState("S0", State{"S0", 0, false}, State{"S1", 1, true})
	automata.AddTransition(
		Transition{Source: "S0", Destination: "S1", Trigger: "1"},
		Transition{Source: "S1", Destination: "S0", Trigger: "0"},
	)

	automata.Event("1")
	res, _ := automata.FinalState()
	assert.Equal(t, res, "S1=1")

	automata.Event("0")
	_, err := automata.FinalState()
	assert.EqualError(t, err, errNoFinalState.Error())
}

func TestE2E(t *testing.T) {
	var err error
	newAutomata := func() fsa {
		automata := fsa{
			alphabet: make(map[string]struct{}),
			states: make(map[string]State),
			transitions: make(map[string]string),
		}
		automata.AddAlphabet("0", "1")
		automata.AddState("S0", State{"S0", 0, true}, State{"S1", 1, true}, State{"S2", 2, true})
		automata.AddTransition(
			Transition{Source: "S0", Destination: "S0", Trigger: "0"},
			Transition{Source: "S0", Destination: "S1", Trigger: "1"},
			Transition{Source: "S1", Destination: "S2", Trigger: "0"},
			Transition{Source: "S1", Destination: "S0", Trigger: "1"},
			Transition{Source: "S2", Destination: "S1", Trigger: "0"},
			Transition{Source: "S2", Destination: "S2", Trigger: "1"},
		)
		return automata
	}
	
	cases :=  []struct{
		input string
		result string
		err string
	}{
		{"110", "S0=0", ""},
		{"1010", "S1=1", ""},
		{"10100", "S2=2", ""},
		{"11A", "S0=0", "Invalid character within input string, A is not listed as part of the allowing alphabet"},
	}

	for idx, c := range cases {
		automata := newAutomata()
		for _, i := range c.input {
			err = automata.Event(string(i))
			if err != nil {
				break
			}	
		}

		if idx < 3 {
			assert.NoError(t, err, c.err)
		} else if idx == 3 {
			assert.EqualError(t, err, c.err)
		}

		res, _ := automata.FinalState()
		if idx < 3 {
			assert.Equal(t, c.result, res)
		}
	}
}
