# Finite State Automata

**A simple Finite State Automata in Go**

The automata is taking a list of states, with an associated int value and a boolean to define if the state if a final state or not.
A list of character as event trigger for a state change and a list of transition rule used
to define the transition from a state to another based on a specific character trigger event.

example: 
  - S0 and S1 are states
  - "0" and "1" are alphabet event triggers
  - the following transition rules define a change of state
    S0 -> 1 -> S1
    S0 -> 0 -> S0
    etc...

    the program will take a string as input that includes only the characters listed in the alphabets and change the state accordingly 
    to the transition rule. At the end of the string, of the state is a final state, it will return the final state value.

## How to use

- create a new FSA interface
- define the alphabet of allowing event input character with AddAlphabet()
- define the various states with AddState()
- define the transition rules between the states with AddTransition()
- pass the list of event to change the automata state.
- find the automata final state at the end of the input parsing with FinalState

## code example

```
package main

import (
	fmt "fmt"
	log "log"

	fsa "github.com/sebastien6/finite-state-automata"
)

func main() {
	var err error
	var input string

	// capture user input
	fmt.Print("Enter input (string of 0 and 1 only): ")
	fmt.Scanln(&input)
	
	// create a new FSA interface
	automata := fsa.NewFSA()

	// set list of accepted input characters
	automata.AddAlphabet("0", "1")

	// set list of state
	automata.AddState("S0", 
		fsa.State{
			Name: "S0",
			Value: 0,
			FinalState:true,
		}, 
		fsa.State{
			Name:"S1",
			Value:1,
			FinalState:true,
		}, 
		fsa.State{
			Name: "S2",
			value: 2,
			FinalState: true,
		})

	// set list of transition
	automata.AddTransition(
		fsa.Transition{Source: "S0", Destination: "S0", Trigger: "0"},
		fsa.Transition{Source: "S0", Destination: "S1", Trigger: "1"},
		fsa.Transition{Source: "S1", Destination: "S2", Trigger: "0"},
		fsa.Transition{Source: "S1", Destination: "S0", Trigger: "1"},
		fsa.Transition{Source: "S2", Destination: "S1", Trigger: "0"},
		fsa.Transition{Source: "S2", Destination: "S2", Trigger: "1"},
	)

	// loop through input to transition states
	for idx, c := range input {
		if idx == 0 {
			if err = automata.Validate(); err != nil {
				log.Fatal(err)
			}
		}

		if err = automata.Event(string(c)); err != nil {
			log.Fatal(err)
		}
	}

	// fetch final state
	res, err := automata.FinalState()
	if err !=  nil {
		log.Fatal(err)
	}
	
	fmt.Println(res)
}

```

Running it can produce the following output:

*Enter input (string of 0 and 1 only): 110*

*S0=0*

## Test

run the following command:

```go test github.com/sebastien6/finite-state-automata```

## Run example

Folder example include the code above in the file main.go.
