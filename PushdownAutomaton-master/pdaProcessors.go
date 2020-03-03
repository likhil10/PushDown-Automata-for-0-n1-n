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
}

// Unmarshals the jsonText string. Returns true if it succeeds.
func (pda *PdaProcessor) Open(jsonText string) (bool){

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
func (pda *PdaProcessor) Reset(){
	pda.CurrentState = pda.StartState
	pda.TokenStack = []string{}
}

func Put(pda *PdaProcessor, s string){
	input_length := len(s)
	transitions := pda.Transitions
	transition_length := len(transitions)
	pda.CurrentState = pda.StartState
	pda.CurrentStack = "null"
	matching := false
	for i:= 0; i < input_length; i++ {
		fmt.Println("i ran for ", i)
		pda.CurrentStack = "null"
		char := string(s[i])
		matching = false

		for j := 0 ; j < transition_length ; j++ {	
			t := transitions[j]
			if t[0] == pda.CurrentState && t[1] == char && t[2] == pda.CurrentStack { 
				fmt.Println("IFF BLOCK ")
				matching = true
				pda.TransitionStack = append(pda.TransitionStack, pda.CurrentState)
				pda.CurrentState = t[3]
				fmt.Println("Tocken Stack", pda.TransitionStack)
				if t[4] != "null" {
					Push(pda, t[4])
					pda.CurrentStack = t[4]
				} else {
					if len(pda.TokenStack) == 0 {
						matching = false
						break
					} else {
						Pop(pda)
						break
					}

				}
			}
			if len(pda.TokenStack) > 0 {
				pda.CurrentStack = pda.TokenStack[len(pda.TokenStack)-1]
			} else {
				break
			}
		}
		if matching == false {
			break
		}
	}
	if len(pda.TransitionStack) > 0 {
		if pda.TransitionStack[len(pda.TransitionStack)-1] == "q3" && len(pda.TokenStack) == 0 && matching == true {
			pda.CurrentStack = pda.Eos
			pda.CurrentState = "q4"
			pda.TransitionStack = append(pda.TransitionStack, pda.CurrentState)
			fmt.Println("pda.TransitionStackFinal", pda.TransitionStack)
		}
	}
	is_accepted(pda, matching)
}

func is_accepted(pda *PdaProcessor, matching bool)  {
	if len(pda.TokenStack) == 0 && matching == true {
		fmt.Println("Accepted: TRUE")
	} else {
		fmt.Println("Accepted: FALSE")
	}
}

func Peek(pda *PdaProcessor) {
	if len(pda.TokenStack) == 0 && pda.CurrentState == "q3" {
		pda.CurrentStack = "$"
		// pda.char = "null"
	} else if len(pda.TokenStack) == 0 {
		pda.CurrentStack = "null"
	} else {
		pda.CurrentStack = "0"
	}
}

func Push(pda *PdaProcessor, x string)  {
	fmt.Println("push ", x)
	pda.TokenStack = append(pda.TokenStack, x)
 	fmt.Println("PUSHED= ", pda.TokenStack)
}

func Pop(pda *PdaProcessor)  {
	pda.TokenStack = pda.TokenStack[:len(pda.TokenStack) - 1]
	fmt.Println("POPPED= ", pda.TokenStack)
}

// A function that calls panic if it detects an error.
func check(e error) {
	if e != nil{
		panic(e)
	}
}

func main() {
	if len(os.Args) < 2{
		fmt.Println("Error: command-line args must include JSON file path")
		os.Exit(0)
	}
	jsonFilename := string(os.Args[1])
	jsonText, err := ioutil.ReadFile(jsonFilename)
	
	pda := new(PdaProcessor)
	
	if pda.Open(string(jsonText)){
		// fmt.Println(pda)
	} else {
		fmt.Println("Error: could not open JSON file")
	}

	check(err)

	var s string
	fmt.Print("Enter the input string: ")
  	fmt.Scan(&s)
  	Put(pda, s)
}