package main

import (
	"fmt"
	"math/rand"
	"time"
)

// initialize global list of node objects and list to track status
var listOfNodes [10]Node
var statusList []bool

func main() {
	// ask for input
	message, protocolCode := askInput()

	// initializes 10 nodes for system indexed 0 - 9
	for i:=0; i < 10; i++ {
		if i == 0 { // initial infected node
			listOfNodes[i] = createNode(i, true, message)
			statusList[i] = true
		} else { // remaining susceptible nodes
			listOfNodes[i] = createNode(i, false, "ready")
			statusList[i] = false
		}
	}

	// start gossip protocol
	time1 := time.Now() // captures start time of protocol
	for i := 0; i < 10; i++ {
		go runNode( listOfNodes[i], protocolCode)
	}
	time2 := time.Now() // captures end time of protocol
	timeDiff := time2.Sub(time1)

	fmt.Println("Gossip Protocol is complete!")
	fmt.Println("The protocol took %d",  &timeDiff)
}

// asks for user input of message and desired gossip algorithm
func askInput() (string, string) {
	fmt.Println("Hello! Type a word that you would like to send")
	var message string
	fmt.Scanln(&message)

	var code string
	var printCode string
	// cycles until a, b, or c is selected
	error := false
	for (!error) {
		fmt.Println("Enter the letter corresponding to the desired gossip algorithm:")
		fmt.Println("a. Push")
		fmt.Println("b. Pull")
		fmt.Println("c. Push-Pull")
		fmt.Scanln(&code)

		if code == "a"{
			printCode = "Push"
			error = true
		} else if code == "b"{
			printCode = "Pull"
			error = true
		} else if code == "c"{
			printCode = "Push-Pull"
			error = true
		} else {
			fmt.Println("Invalid algorithm code. Try again!")
		}
	}

	fmt.Println("Great! We will send" + message + "using the" + printCode + "algorithm.")
	return message, code
}

func push(currNode Node) {
	// executes as long as the node is susceptible
	if (!currNode.status) {
		reception := <-currNode.pushChan // waits for message in pushChan
		currNode.msg = reception // sets Node's message to reception string
		currNode.status = true // sets Node's status to infected
		statusList[currNode.id] = true // tells array that node is infected now
	}
	// executes when the node becomes infected
	for true {
		// test to see if there are susceptible nodes remaining
		// if all nodes are infected, then we break the loop and do not perform push
		if sumBool(statusList) == 10 {
			break
		}
		pushTo := pickNode(currNode) // choose random node to push to
		pushTo.pushChan <- currNode.msg // send message through the receiving node's channel
	}
}

func pull(currNode Node) {
	lastNode := false // tells if this is the final node to be infected
	// executes while node is susceptible, picks nodes to request until becomes infected
	for (!currNode.status) {
		pullFrom := pickNode(currNode) // choose random node to pull from
		if pullFrom.status { // if the pull node is infected, send the message from pullFrom
			pullFrom.pullChan <- currNode.id // sends id to pullFrom
			reception := <- currNode.pushChan // wait for message from pullFrom
			currNode.status = true // set Node's status to infected
			currNode.msg = reception
			statusList[currNode.id] = true // tells array that node is infected now

			// checks to see if this is the last node to be infected
			if sumBool(statusList) == 10 {
				lastNode = true
				// need to inform other goroutines to stop
				for i := 0; i < 10; i++ {
					if i != currNode.id {
						tempNode := listOfNodes[i]
						tempNode.pullChan <- -1 // sends invalid id pull request to all nodes except currNode
					}
				}
			}
		}
	}
	// executes once node is infected as long as there are still susceptible nodes
	for !lastNode {
		sendToId := <-currNode.pullChan // waits for request in pullChan
		// ends goroutine if it receives signal from last infected node
		if sendToId == -1 {
			break
		} else {
			listOfNodes[sendToId].pushChan <- currNode.msg // sends message through pull Node's pushChan
		}
	}
}

// used to make goroutine that runs protocol based on input from user
func runNode(currNode Node, protocol string) {
	if protocol == "a" {
		push(currNode)
	}
	if protocol == "b" {
		pull(currNode)
	}
	//if protocol == "c" {
	//	pushPull(currNode)
	//}
}

// picks a random node from the global list listOfNodes.
// If the node picks itself, it keeps running until it picks a node that isn't itself.
// returns node chosen
func pickNode(primeNode Node) (Node) {
	var randomId int
	x := false
	for x {
		rand.Seed(time.Now().UnixNano()) // makes it so that the random int is not deterministic
		randomId = rand.Intn(10) - 1 // random int between 0 and 9 inclusive
		if randomId != primeNode.id {
			x = true
		}
	}
	pickedNode := listOfNodes[randomId]
	return pickedNode
}

// returns integer representing how many entries in boolean array are true
func sumBool(list []bool) int{
	sum := 0
	for _, entry := range list {
		if entry {
			sum++
		}
	}
	return sum
}
