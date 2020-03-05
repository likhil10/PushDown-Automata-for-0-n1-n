package main

import (
	"encoding/json"
	"fmt"
	"os"
	"io/ioutil"
)

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

	//This keeps a count of everytime put method is called
	PutCounter int
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
	pda.PutCounter += 1
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

func eos()  {
	
}

func current_state()  {
	
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

// Defines the type PdaProcessor.
type PdaProcessor struct {
	// Note: field names must begin with a capital in order to be recognized by the JSON Marshaller
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

	// The slice is used to hold the tokens.
	TokenStack []string

	// The slice is used to hold the tokens.
	Stack []string

	// Flag to find if the input is valid
	accepting_transition bool

}

type PdaDriver struct {
	// Note: field names must begin with a capital in order to be recognized by the JSON Marshaller
	Name string
	States []string
	InputAlphabet []string
	StackAlphabet []string
	AcceptingStates []string
	StartState string
	Transitions [][]string
	Eos string

	// Holds the current state.
	CurrentState string

	// Token at the top of the stack.
	CurrentStack string

	// The slice is used to hold the tokens.
	TokenStack []string

	// The slice is used to hold the tokens.
	Stack []string

	// Flag to find if the input is valid
	accepting_transition bool

	char string

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
	pda.Stack = []string{}
}

func Put(pdr *PdaDriver, p *PdaProcessor, s string){
	input_length := len(s)
	transitions := p.Transitions
	transition_length := len(transitions)
	pdr.CurrentState = p.StartState
	pdr.CurrentStack = "null"
	// fmt.Println("transition_length \n transitions \n", transition_length, transitions)
	for i:= 0; i < input_length; i++ {
		fmt.Println("i ran for ", i)
		pdr.char = string(s[i])
		pdr.accepting_transition = false
		pdr.CurrentStack = "null"

		for j := 0 ; j < transition_length ; j++ {	
			// fmt.Println("transitions[j][0]", transitions[j][0])
			// fmt.Println("pdr.CurrentState", pdr.CurrentState)
			// fmt.Println("transitions[j][1]", transitions[j][1])
			// fmt.Println("pdr.char", pdr.char)
			// fmt.Println("transitions[j][2]", transitions[j][2])
			// fmt.Println("pdr.CurrentStack", pdr.CurrentStack)
			if transitions[j][0] == pdr.CurrentState && transitions[j][1] == pdr.char && transitions[j][2] == pdr.CurrentStack { 
				// fmt.Println("IFF BLOCK ")
				pdr.accepting_transition = true
				pdr.CurrentState = transitions[j][3]
				pdr.TokenStack = append(pdr.TokenStack, transitions[j][3])
				fmt.Println("Tocken Stack", pdr.TokenStack)
				// fmt.Println("Updated Current State= ", pdr.CurrentState)
				if transitions[j][4] != "null"{
					Push(pdr, transitions[j][4])
				} else {
					Pop(pdr)
				}
			} else if len(pdr.TokenStack) > 0 && len(pdr.Stack) == 0{
				  if pdr.TokenStack[(len(pdr.TokenStack)-1)] == "q4" {
					pdr.accepting_transition = true
				 }
			}
			Peek(pdr)
			// if pdr.accepting_transition ==  true {
			// 	break
			// }
		}

		// if pdr.accepting_transition ==  false || (i == input_length-1 && (len(pdr.Stack) > 0)) {
		// 	fmt.Println("failure.")
		// // 	break
		// }
	}
	if len(pdr.TokenStack) > 0 {
		if pdr.TokenStack[0] == "q2" && pdr.TokenStack[(len(pdr.TokenStack)-1)] == "q4" {
			fmt.Println("TokenStack[0]", pdr.TokenStack[0])
			fmt.Println("TokenStack[(len(pdr.TokenStack)-1)]", pdr.TokenStack[(len(pdr.TokenStack)-1)])
			fmt.Println("pdr.accepting_transition", pdr.accepting_transition)
			fmt.Println("len(pdr.Stack)", len(pdr.Stack))
			if  pdr.accepting_transition == true  && len(pdr.Stack) == 0 {	
				fmt.Println("SUCCESS")
			} else {
				fmt.Println("HERE 0")
				fmt.Println("FAILURE")
			}
		} else {
			fmt.Println("HERE 1")
			fmt.Println("FAILURE")
		}
	} else {
		// fmt.Println("HERE 2")
		fmt.Println("FAILURE")
	}
}

func Peek(pdr *PdaDriver) {
	if len(pdr.Stack) == 0 && pdr.CurrentState == "q3" {
		pdr.CurrentStack = "$"
		pdr.char = "null"
	} else if len(pdr.Stack) == 0 {
		pdr.CurrentStack = "null"
	} else {
		pdr.CurrentStack = "0"
	}
}

func Push(pdr *PdaDriver, x string)  {
	// fmt.Println("push 1", x)
	if pdr.CurrentState == pdr.StartState && x == "null" {
		// fmt.Println("push 2", x)
		pdr.accepting_transition = false
	} else {
		pdr.Stack = append(pdr.Stack, x)
		fmt.Println("PUSHED")
	}
}

func Pop(pdr *PdaDriver)  {
	fmt.Println("POP 1")
	if len(pdr.Stack) == 0 && pdr.CurrentState == "q4" {
		fmt.Println("pop 2")
		pdr.accepting_transition = true
	} else if len(pdr.Stack) == 0{
		fmt.Println("pop 3")
		pdr.accepting_transition = false
	} else {
		fmt.Println("POPED")
		pdr.Stack = pdr.Stack[:len(pdr.Stack) - 1]
	}
}

// A function that calls panic if it detects an error.
func check(e error) {
	if e != nil{
		panic(e)
	}
}

func main() {
	if len(os.Args) < 2{
		fmt.Println("Error: command-line args must include json spec file path")
		os.Exit(0)
	}
	jsonFilename := string(os.Args[1])
	jsonText, err := ioutil.ReadFile(jsonFilename)
	
	check(err)

	pda := new(PdaProcessor)
	pdr := new(PdaDriver)
	if pda.Open(string(jsonText)){
		// fmt.Println(pda)
	} else {
		fmt.Println("Error: could not open json spec")
	}

	var s string
	fmt.Print("Enter the input string: ")
  	fmt.Scan(&s)
  	Put(pdr, pda, s)
}