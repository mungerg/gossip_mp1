package gossip_mp1

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
	"math/rand"
)

var listOfNodes []node

type node struct {
	id     int     // holds ID of node
	status bool // false corresponds to susceptible, true corresponds to infected
	msg    string  // holds message
}

func main() {
	message, code := askInput()

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
func gossip( nodeChan chan <-string,  allInfected <-chan bool, primeNode node, protocol string, message string){
	lenOfList  := len(listOfNodes)
	for allInfected != true{
		//WAIT TIME

		//Pick a random node-This will probably need to be changed according to the protocol developed in chapter 4 that I havent read yet
		receiver := pickNode(primeNode, lenOfList)
		if protocol == "Push" && primeNode.status==true {
			push(receiver, nodeChan, message)
		}
		//Need something to do the updating for the receiver. Like right now you could just change the status of the receiver right here

	}

}
func pickNode(primeNode node, lenOfList int) node {
	var randomNode int
	x := false
	for x == false {
		randomNode :=rand.Intn(lenOfList)
		if randomNode != primeNode.id {
			x = true
		}
	}
	receiver := listOfNodes[randomNode]
	return receiver
}
func push(receiver node, nodeChan chan <-string, message string){

}


