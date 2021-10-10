package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// initialize global list of node objects and list to track status
var listOfNodes [10]Node
var statusList [10]bool

var wg sync.WaitGroup // makes main wait for go routines to finish

// data structure to hold Node data
type Node struct {
	id       int         // holds ID of node
	status   bool        // false corresponds to susceptible, true corresponds to infected
	msg      string      // holds message
	pullChan chan int    // for receiving pull requests
	pushChan chan string // for pushing message through
}

// function for creating nodes
func createNode(id int, status bool, msg string) Node {
	//Why is this one an int? Is it choosing based off id?
	pull := make(chan int, 256)
	push := make(chan string, 1)
	node := Node{id, status, msg, pull, push}
	return node
}

func main() {
	wg.Add(10)

	// ask for input
	message, protocolCode := askInput()

	// initializes 10 nodes for system indexed 0 - 9
	for i := 0; i < 10; i++ {
		if i == 0 { // initial infected node
			listOfNodes[i] = createNode(i, true, message)
			statusList[i] = true
		} else { // remaining susceptible nodes
			listOfNodes[i] = createNode(i, false, "ready")
			statusList[i] = false
		}
	}
	fmt.Println("sumBool is " + strconv.Itoa(sumBool(statusList)))

	// start gossip protocol
	time1 := time.Now() // captures start time of protocol
	for i := 0; i < 10; i++ {
		fmt.Println("Inside loop cycle " + strconv.Itoa(i))
		go runNode(&wg, listOfNodes[i], protocolCode)
	}
	wg.Wait()
	time2 := time.Now() // captures end time of protocol
	timeDiff := time2.Sub(time1)

	fmt.Println("Gossip Protocol is complete!")
	fmt.Println("The protocol took ", &timeDiff)
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
	for !error {
		fmt.Println("Enter the letter corresponding to the desired gossip algorithm:")
		fmt.Println("a. Push")
		fmt.Println("b. Pull")
		fmt.Println("c. Push-Pull")
		fmt.Scanln(&code)

		if code == "a" {
			printCode = "Push"
			error = true
		} else if code == "b" {
			printCode = "Pull"
			error = true
		} else if code == "c" {
			printCode = "Push-Pull"
			error = true
		} else {
			fmt.Println("Invalid algorithm code. Try again!")
		}
	}

	fmt.Println("Great! We will send " + message + " using the " + printCode + " algorithm.")
	return message, code
}

func push(currNode Node) {
	fmt.Println("Inside of push function for " + strconv.Itoa(currNode.id))
	// executes as long as the node is susceptible
	if !currNode.status {
		reception := <-currNode.pushChan // waits for message in pushChan
		currNode.msg = reception         // sets Node's message to reception string
		currNode.status = true           // sets Node's status to infected
		statusList[currNode.id] = true   // tells array that node is infected now
	}
	fmt.Println(strconv.Itoa(currNode.id) + " is infected! " + strconv.Itoa(10-sumBool(statusList)) +
		" left to infect.")
	// executes when the node becomes infected
	for true {
		// test to see if there are susceptible nodes remaining
		// if all nodes are infected, then we break the loop and do not perform push
		if sumBool(statusList) == 10 {
			break
		}
		pushTo := pickNode(currNode) // choose random node to push to
		if pushTo.status == false {
			fmt.Println(strconv.Itoa(currNode.id) + " picked node " + strconv.Itoa(pushTo.id))
			pushTo.pushChan <- currNode.msg // send message through the receiving node's channel
		}
	}
}

func pull(currNode Node) {
	fmt.Println("Inside pull function for " + strconv.Itoa(currNode.id))
	lastNode := false // tells if this is the final node to be infected
	// executes while node is susceptible, picks nodes to request until becomes infected
	for !currNode.status {
		fmt.Println(currNode.status, currNode.id)
		pullFrom := pickNode(currNode) // choose random node to pull from
		if pullFrom.status {           // if the pull node is infected, send the message from pullFrom
			pullFrom.pullChan <- currNode.id // sends id to pullFrom
			reception := <-currNode.pushChan // wait for message from pullFrom
			currNode.status = true           // set Node's status to infected
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
func runNode(wg *sync.WaitGroup, currNode Node, protocol string) {
	defer wg.Done()
	fmt.Println("Starting runNode for " + strconv.Itoa(currNode.id))
	if protocol == "a" {
		push(currNode)
	}
	if protocol == "b" {
		pull(currNode)
	}
	if protocol == "c" {
		pushPull(currNode)
	}
	fmt.Println("Ending runNode for " + strconv.Itoa(currNode.id))
}

//Based on page 7 and 8 of Gossip by Mark Jelasity. When the proportion of nodes is less than .5, pull is faster than push. "In fact, the quadratic convergence phase,
//roughly after st < 0.5, lasts only for O(log log N) cycles"
func pushPull(currNode Node) {

	if !currNode.status {
		reception := <-currNode.pushChan // waits for message in pushChan
		currNode.msg = reception         // sets Node's message to reception string
		currNode.status = true           // sets Node's status to infected
		statusList[currNode.id] = true   // tells array that node is infected now
	}
	// executes when the node becomes infected
	for true {
		// test to see if there are susceptible nodes remaining
		// if all nodes are infected, then we break the loop and do not perform push
		if sumBool(statusList) == len(listOfNodes)/2 {
			pull(currNode)
			break
		}
		pushTo := pickNode(currNode)    // choose random node to push to
		pushTo.pushChan <- currNode.msg // send message through the receiving node's channel
	}
}

// picks a random node from the global list listOfNodes.
// If the node picks itself, it keeps running until it picks a node that isn't itself.
// returns node chosen
func pickNode(primeNode Node) Node {
	var randomId int
	for true {
		rand.Seed(time.Now().UnixNano()) // makes it so that the random int is not deterministic
		randomId = rand.Intn(10)         // random int between 0 and 9 inclusive
		if randomId != primeNode.id {
			break
		}
	}
	pickedNode := listOfNodes[randomId]
	return pickedNode
}

// returns integer representing how many entries in boolean array are true
func sumBool(list [10]bool) int {
	sum := 0
	for _, entry := range list {
		if entry {
			sum++
		}
	}
	return sum
}
