package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// initialize global list of node objects and list to track status
var (
	mu          sync.Mutex // allows for locking of variable
	listOfNodes [10]Node
	statusList  [10]bool
)

var wg sync.WaitGroup // makes main wait for go routines to finish

// data structure to hold data being pushed through channel to other node
type pushChanData struct {
	message  string // message being sent
	pushNode *Node  // pointer to Node sent from
}

// data structure to hold Node data
type Node struct {
	id       int               // holds ID of node
	status   bool              // false corresponds to susceptible, true corresponds to infected
	msg      string            // holds message
	pullChan chan int          // for receiving pull requests
	pushChan chan pushChanData // for pushing message and "from" node through
}

// function for creating nodes
func createNode(id int, status bool, msg string) Node {
	pull := make(chan int, 10)
	push := make(chan pushChanData, 10)
	node := Node{id, status, msg, pull, push}
	return node
}

func main() {
	wg.Add(10)

	// ask for input
	message, protocolCode := askInput()

	// initializes 10 nodes for system indexed 0 - 9
	for i := 0; i < 10; i++ {
		if i == 9 { // initial infected node
			listOfNodes[i] = createNode(i, true, message)
			statusList[i] = true
		} else { // remaining susceptible nodes
			listOfNodes[i] = createNode(i, false, "waiting")
			statusList[i] = false
		}
	}

	// start gossip protocol
	time1 := time.Now() // captures start time of protocol
	for i := 0; i < 10; i++ {
		go runNode(&wg, &listOfNodes[i], protocolCode)
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
	error1 := false
	for !error1 {
		fmt.Scanln(&message)
		if message != "switch" {
			error1 = true
		} else {
			fmt.Println("Invalid message. Pick another message to send.")
		}
	}

	var code string
	var printCode string
	// cycles until a, b, or c is selected
	error2 := false
	for !error2 {
		fmt.Println("Enter the letter corresponding to the desired gossip algorithm:")
		fmt.Println("a. Push")
		fmt.Println("b. Pull")
		fmt.Println("c. Push-Pull")
		fmt.Scanln(&code)

		if code == "a" {
			printCode = "Push"
			error2 = true
		} else if code == "b" {
			printCode = "Pull"
			error2 = true
		} else if code == "c" {
			printCode = "Push-Pull"
			error2 = true
		} else {
			fmt.Println("Invalid algorithm code. Try again!")
		}
	}

	fmt.Println("Great! We will send " + message + " using the " + printCode + " algorithm.")
	return message, code
}

func push(currNode *Node) {
	// executes as long as the node is susceptible
	if !currNode.status {
		receptionData := <-currNode.pushChan  // waits for message in pushChan
		receptionData.pushNode.pullChan <- -2 // sends confirmation "-2" to Node pushed from
	}
	// executes when the node becomes infected
	for true {
		// test to see if there are susceptible nodes remaining
		// if all nodes are infected, then we break the loop and do not perform push
		if sumBoolLocks(statusList) == 10 {
			break
		}
		mu.Lock()                    // locks global lists
		pushTo := pickNode(currNode) // choose random node to push to
		if pushTo.status == false {
			pushTo.pushChan <- pushChanData{currNode.msg, currNode} // send message through the receiving node's channel
			confirmation := <-currNode.pullChan                     // waits for confirmation from receiving node
			if confirmation == -2 {                                 // added if statement to avoid unused variable error
				pushTo.status = true         // sets receiving node's status to infected
				statusList[pushTo.id] = true // tells array that node is infected node
			}
		}
		mu.Unlock() // unlocks global lists
	}
}

func pull(currNode *Node) {
	lastNode := false // tells if this is the final node to be infected
	// executes while node is susceptible, picks nodes to request until becomes infected
	for !currNode.status {
		mu.Lock()                      // locks global lists
		pullFrom := pickNode(currNode) // choose random node to pull from
		if pullFrom.status {           // if the pull node is infected, send the message from pullFrom
			pullFrom.pullChan <- currNode.id     // sends id to pullFrom
			receptionData := <-currNode.pushChan // wait for message from pullFrom
			currNode.status = true               // set Node's status to infected
			currNode.msg = receptionData.message
			statusList[currNode.id] = true // tells array that node is infected now

			// checks to see if this is the last node to be infected
			if sumBoolNoLocks(statusList) == 10 {
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
		mu.Unlock() // unlocks global lists
	}
	// executes once node is infected as long as there are still susceptible nodes
	for !lastNode {
		sendToId := <-currNode.pullChan // waits for request in pullChan
		// ends goroutine if it receives signal from last infected node
		if sendToId == -1 {
			break
		} else {
			listOfNodes[sendToId].pushChan <- pushChanData{currNode.msg, currNode} // sends message through pull Node's pushChan
		}
	}
}

//Based on page 7 and 8 of Gossip by Mark Jelasity. When the proportion of nodes is less than .5, pull is faster than push. "In fact, the quadratic convergence phase,
//roughly after st < 0.5, lasts only for O(log N) cycles"
func pushPull(currNode *Node) {
	// start with push protocol while less than half the nodes are infected
	for true {
		// checks if node is susceptible during push
		if !currNode.status {
			receptionData := <-currNode.pushChan // waits for message in pushChan
			// catch for receiving switch message to switch protocols
			if receptionData.message == "switch" {
				break
			} else { // received message, needs to be changed to infected
				currNode.msg = receptionData.message  // sets Node's message to reception string
				receptionData.pushNode.pullChan <- -2 // sends confirmation to Node pushed from
				// check to see if we need to end the push protocol
				mu.Lock() // locks global lists
				if sumBoolNoLocks(statusList) == len(listOfNodes)/2 {
					// inform susceptible goroutines to stop
					for i := 0; i < 10; i++ {
						if i != currNode.id && !listOfNodes[i].status {
							tempNode := listOfNodes[i]
							tempNode.pushChan <- pushChanData{"switch", currNode} // sends invalid id pull request to all nodes except currNode
						}
					}
				}
				mu.Unlock() // unlocks global lists
			}
		} else { // executes when the node becomes infected during push
			mu.Lock() // locks global lists
			// breaks loop if it is time to switch to pull
			if sumBoolNoLocks(statusList) >= len(listOfNodes)/2 {
				mu.Unlock()
				break
			}
			pushTo := pickNode(currNode) // choose random node to push to
			if pushTo.status == false {
				pushTo.pushChan <- pushChanData{currNode.msg, currNode} // send message through the receiving node's channel
				confirmation := <-currNode.pullChan                     // waits for confirmation from receiving node
				if confirmation == -2 {                                 // added if statement to avoid unused variable error
					pushTo.status = true         // sets receiving node's status to infected
					statusList[pushTo.id] = true // tells array that node is infected node
				}
			}
			mu.Unlock() // unlocks global lists
		}
	}
	// switches to pull after half the nodes are infected
	pull(currNode)
}

// used to make goroutine that runs protocol based on input from user
func runNode(wg *sync.WaitGroup, currNode *Node, protocol string) {
	defer wg.Done()
	if protocol == "a" {
		push(currNode)
	} else if protocol == "b" {
		pull(currNode)
	} else if protocol == "c" {
		pushPull(currNode)
	}
}

// picks a random node from the global list listOfNodes.
// If the node picks itself, it keeps running until it picks a node that isn't itself.
// returns node chosen
func pickNode(primeNode *Node) *Node {
	var randomId int
	for true {
		rand.Seed(time.Now().UnixNano()) // makes it so that the random int is not deterministic
		randomId = rand.Intn(10)         // random int between 0 and 9 inclusive
		if randomId != primeNode.id {
			break // cycles through until we find a node that is not the one picking
		}
	}

	pickedNode := &listOfNodes[randomId]
	return pickedNode
}

// returns integer representing how many entries in boolean array are true using locks within function
func sumBoolLocks(list [10]bool) int {
	mu.Lock() // locks global lists
	sum := 0
	for _, entry := range list {
		if entry {
			sum++
		}
	}
	defer mu.Unlock() // unlocks global list when function finishes running
	return sum
}

// returns integer representing how many entries in boolean array are true not using locks within function
func sumBoolNoLocks(list [10]bool) int {
	sum := 0
	for _, entry := range list {
		if entry {
			sum++
		}
	}
	return sum
}
