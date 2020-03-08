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
			Value: "0",
			FinalState:true,
		}, 
		fsa.State{
			Name:"S1",
			Value:1,
			FinalState:true,
		}, 
		fsa.State{
			Name: "S2",
			Value: 2,
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