package gossip_mp1

import (
	"fmt"
)

type node struct {
	id     int     // holds ID of node
	status boolean // false corresponds to susceptible, true corresponds to infected
	msg    string  // holds message
}

func main() {
	message, code := askInput()

	complete := make(chan bool, 1)

}

// asks for user input of web addresses
func askInput() (string, string, string) {
	fmt.Println("Hello! Type a word that you would like to send")
	var message string
	fmt.Scanln(&message)

	fmt.Println("Enter the letter corresponding to the desired gossip algorithm:")
	fmt.Println("a. Push")
	fmt.Println("b. Pull")
	fmt.Println("c. Push-Pull")
	var code string
	fmt.Scanln(&code)

	var printCode string
	if code == "a" {
		printCode = "Push"
	} else if code == "b" {
		printCode = "Pull"
	} else {
		printCode = "Push-Pull"
	}

	fmt.Println("Great! We will send" + message + "using the" + printCode + "algorithm.")

	return message, code
}

// initialize the value for each node,
//let each node start as susceptible and hold empty message
func initNode(node) {
	node.status = false
	node.msg = ""
}

func infect(node.msg, c)
