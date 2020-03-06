package main

import (
	"encoding/json"
	"fmt"
	"os"
	"io/ioutil"
)

// {"name":"HelloPDA",
// "states":["q1","q2","q3","q4"],
// "input_alphabet":["0","1"],
// "stack_alphabet":["0","1"],
// "accepting_states":["q1","q4"],
// "start_state":"q1",
// "transitions":[
//  ["q1","null","null","q2","$"], first transition state.
// 	["q2","0","null","q2","0"],
// 	["q2","1","0","q3","null"],
// 	["q3","1","0","q3","null"],
// 	["q3","null","$","q4","null"]],
// "eos":"$"}

// Structure of type PdaProcessor.
type PdaProcessor struct {
	// Note: field names must begin with capital letter for JSON
	Name string `json:"name"`
	States []string `json:"states"`
	InputAlphabet []string `json:"input_alphabet"`
	StackAlphabet []string `json:"stack_alphabet"`
	AcceptingStates []string `json:"accepting_states"`
	StartState string `json:"start_state"`
	Transitions [][]string `json:"transitions"`
	Eos string `json:"eos"`

	// Holds the current state.
	CurrentState string

	// Token at the top of the stack.
	CurrentStack string

	// This slice is used to hold the transition states tokens.
	TransitionStack []string

	// This slice is used to hold the token stack.
	TokenStack []string

	// This keeps a count of everytime put method is called
	PutCounter int

	// This keeps a count for everytime a transition  is changed
	TransitionCounter int

	// This checks if the input is accepted by the PDA
	IsAccepted bool
}

// Unmarshals the jsonText string. Returns true if it succeeds.
func (pda *PdaProcessor) open(jsonText string) (bool){

	if err := json.Unmarshal([]byte(jsonText), &pda); err != nil {
		check(err)
	}

	// Validate input.	
	if len(pda.Name) == 0 || len(pda.States) == 0 || len(pda.InputAlphabet) == 0 || 
	len(pda.StackAlphabet) == 0 || len(pda.AcceptingStates) == 0 || len(pda.StartState) == 0 ||
	len(pda.Transitions) == 0 || len(pda.Eos) == 0 {
		return false
	}

	return true
}

// Sets the CurrentState to StartState and assigns Stack a new empty slice
func reset(pda *PdaProcessor){
	pda.CurrentState = pda.StartState
	pda.TokenStack = []string{}
}

func put(pda *PdaProcessor, char string){
	pda.PutCounter += 1
	transitions := pda.Transitions
	transition_length := len(transitions)
	if pda.PutCounter == 1 {
		putForTFirst(pda)
	}

	for j := 1 ; j < transition_length ; j++ {	
		t := transitions[j]

		if t[0] == pda.CurrentState && t[1] == char && t[2] == pda.CurrentStack { 
			pda.IsAccepted = true
			pda.TransitionStack = append(pda.TransitionStack, pda.CurrentState)
			pda.TransitionCounter += 1
			pda.CurrentState = t[3]
			pda.TransitionStack = append(pda.TransitionStack, pda.CurrentState)			

			if t[4] != "null" {
				push(pda, t[4])
				pda.CurrentStack = t[4]
			} else {
				if len(pda.TokenStack) == 0 {
					pda.IsAccepted = false
					break
				} else {
					pop(pda)
					break
				}
			}
		}

		if len(pda.TokenStack) > 1 {
			pda.CurrentStack = pda.TokenStack[len(pda.TokenStack)-1]
		} else {
			break
		}
	}
}

func putForTFirst(pda *PdaProcessor)  {
	if pda.Transitions[0][0] == pda.CurrentState {
		pda.TransitionStack = append(pda.TransitionStack, pda.CurrentState)
		pda.CurrentState = pda.Transitions[0][3]
		push(pda, pda.Transitions[0][4])
		pda.TransitionCounter += 1
	}
}

func is_accepted(pda *PdaProcessor) bool {
	if len(pda.TokenStack) == 0 && pda.IsAccepted == true {
		return true
	} else {
		return false
	}
}

func peek(pda *PdaProcessor, n int) []string {
	if len(pda.TokenStack) > 0 {
		if len(pda.TokenStack) < n {
			return pda.TokenStack
		} else if len(pda.TokenStack) > n {
			x := len(pda.TokenStack) - (n-1)
			return pda.TokenStack[x-1:]
		} else if len(pda.TokenStack) == n {
			return pda.TokenStack[:n]
		}
	}
	return pda.TokenStack
}

func push(pda *PdaProcessor, x string)  {
	pda.TokenStack = append(pda.TokenStack, x)
}

func pop(pda *PdaProcessor)  {
	pda.TokenStack = pda.TokenStack[:len(pda.TokenStack) - 1]
}

// A function that calls panic if it detects an error.
func check(e error) {
	if e != nil{
		panic(e)
	}
}

func eos(pda *PdaProcessor)  {
	if len(pda.TransitionStack) > 0 && pda.TransitionStack[0] == "q1" && pda.TransitionStack[len(pda.TransitionStack)-1] == "q4" {
		fmt.Println("Reached the end of string.")
	} else {
		fmt.Println("Did not reach the end of string but EOS was called.")
	}
}

func current_state(pda *PdaProcessor) string {
	return pda.CurrentState
}

func close()  {
	
}

func main() {
	if len(os.Args) < 2{
		fmt.Println("Error: command-line args must include JSON file path")
		os.Exit(0)
	}
	jsonFilename := string(os.Args[1])
	jsonText, err := ioutil.ReadFile(jsonFilename)
	
	pda := new(PdaProcessor)
	
	if pda.open(string(jsonText)){
		// fmt.Println(pda)
	} else {
		fmt.Println("Error: could not open JSON file")
	}

	reset(pda)

	check(err)

	var s string
	var currentState string
	var accepted bool
	var peekStack []string

	fmt.Print("Enter the input string: ")
  	fmt.Scan(&s)
	  
	input_length := len(s)

	pda.CurrentState = pda.StartState
	pda.CurrentStack = "null"
	pda.IsAccepted = false

	for i:= 0; i < input_length; i++ {
		char := string(s[i])
		pda.CurrentStack = "null"
		pda.IsAccepted = false
		if len(char) == 1 {
			put(pda, char)
		}
		if pda.IsAccepted == false {
			break
		}
	}
	if len(pda.TransitionStack) > 0 {
		if pda.TransitionStack[len(pda.TransitionStack)-1] == "q3" && len(pda.TokenStack) == 1 && pda.IsAccepted == true {
			pda.CurrentStack = pda.TokenStack[0]
			pda.CurrentState = "q4"
			pda.TransitionStack = append(pda.TransitionStack, pda.CurrentState)
			pop(pda)
		}
	}

	accepted = is_accepted(pda)
	fmt.Println("Accepted: ", accepted)	
	if accepted {
		fmt.Println("Number of times put was called: ", pda.PutCounter)
		fmt.Println("Number of times Transitions were changed: ", pda.TransitionCounter)
	}

	peekStack = peek(pda, 5)
	fmt.Println("Peek Stack: ",peekStack)

	eos(pda)

	currentState = current_state(pda)
	fmt.Println("Current State: ", currentState)
}