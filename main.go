package main

import (
	"fmt"
	"math/rand"
	"strconv"
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
	fmt.Println("sumBool is " + strconv.Itoa(sumBoolLocks(statusList)))

	// start gossip protocol
	time1 := time.Now() // captures start time of protocol
	for i := 0; i < 10; i++ {
		fmt.Println("Inside loop cycle " + strconv.Itoa(i))
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
	fmt.Println("Inside of push function for " + strconv.Itoa(currNode.id))
	// executes as long as the node is susceptible
	if !currNode.status {
		receptionData := <-currNode.pushChan // waits for message in pushChan
		fmt.Println(strconv.Itoa(currNode.id) + " has received data, sending confirmation to " +
			strconv.Itoa(receptionData.pushNode.id))
		receptionData.pushNode.pullChan <- -2 // sends confirmation to Node pushed from
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
			fmt.Println(strconv.Itoa(currNode.id) + " picked node " + strconv.Itoa(pushTo.id))
			pushTo.pushChan <- pushChanData{currNode.msg, currNode} // send message through the receiving node's channel
			fmt.Println(strconv.Itoa(currNode.id) + " is waiting for confirmation")
			confirmation := <-currNode.pullChan // waits for confirmation from receiving node
			fmt.Println(strconv.Itoa(currNode.id) + " received confirmation " + strconv.Itoa(confirmation))
			pushTo.status = true // sets receiving node's status to infected
			fmt.Println(strconv.Itoa(pushTo.id) + " status changed to " + strconv.FormatBool(pushTo.status))
			statusList[pushTo.id] = true // tells array that node is infected node
			fmt.Println(strconv.Itoa(pushTo.id) + " status changed in global list to " +
				strconv.FormatBool(statusList[pushTo.id]))
		}
		mu.Unlock() // unlocks global lists
		fmt.Println(strconv.Itoa(pushTo.id) + " is infected! " + strconv.Itoa(10-sumBoolLocks(statusList)) +
			" left to infect.")
	}
}

func pull(currNode *Node) {
	fmt.Println("Inside pull function for " + strconv.Itoa(currNode.id))
	lastNode := false // tells if this is the final node to be infected
	// executes while node is susceptible, picks nodes to request until becomes infected
	for !currNode.status {
		mu.Lock()                      // locks global lists
		pullFrom := pickNode(currNode) // choose random node to pull from
		fmt.Println("Pullfrom node is", pullFrom.id, "for node", currNode.id)
		if pullFrom.status { // if the pull node is infected, send the message from pullFrom
			fmt.Println(strconv.Itoa(currNode.id)+" sent request to ", pullFrom.id)
			pullFrom.pullChan <- currNode.id     // sends id to pullFrom
			receptionData := <-currNode.pushChan // wait for message from pullFrom
			fmt.Println(strconv.Itoa(currNode.id)+" received message from", receptionData.pushNode.id)
			currNode.status = true // set Node's status to infected
			fmt.Println(strconv.Itoa(currNode.id)+" status is", currNode.status)
			currNode.msg = receptionData.message
			fmt.Println("after message update")
			statusList[currNode.id] = true // tells array that node is infected now
			fmt.Println("after status list updates", currNode.id, "to", currNode.status)
			//fmt.Println(strconv.Itoa(currNode.id) + " is infected! " + strconv.Itoa(10-sumBool(statusList)) +
			//	" left to infect.")

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
		fmt.Println(strconv.Itoa(currNode.id)+" received request from", sendToId)
		// ends goroutine if it receives signal from last infected node
		if sendToId == -1 {
			break
		} else {
			listOfNodes[sendToId].pushChan <- pushChanData{currNode.msg, currNode} // sends message through pull Node's pushChan
			fmt.Println(strconv.Itoa(currNode.id)+" sent request to ", sendToId)
		}
	}
}

//Based on page 7 and 8 of Gossip by Mark Jelasity. When the proportion of nodes is less than .5, pull is faster than push. "In fact, the quadratic convergence phase,
//roughly after st < 0.5, lasts only for O(log N) cycles"
func pushPull(currNode *Node) {
	// start with push protocol while less than half the nodes are infected
	for sumBoolLocks(statusList) < len(listOfNodes)/2 {
		// checks if node is susceptible during push
		if !currNode.status {
			receptionData := <-currNode.pushChan // waits for message in pushChan
			fmt.Println(currNode.id, "received a message")
			// catch for receiving switch message to switch protocols
			if receptionData.message == "switch" {
				fmt.Println(currNode.id, "stopped waiting for push")
				break
			} else { // received message, needs to be changed to infected
				currNode.msg = receptionData.message // sets Node's message to reception string
				fmt.Println(strconv.Itoa(currNode.id) + " has received data, sending confirmation to " +
					strconv.Itoa(receptionData.pushNode.id))
				receptionData.pushNode.pullChan <- -2 // sends confirmation to Node pushed from
				// check to see if we need to end the push protocol
				if sumBoolLocks(statusList) == len(listOfNodes)/2 {
					// inform susceptible goroutines to stop
					fmt.Println("Informing susceptible goroutines to stop")
					for i := 0; i < 10; i++ {
						if i != currNode.id && !listOfNodes[i].status {
							tempNode := listOfNodes[i]
							fmt.Println("telling", tempNode.id, "to stop waiting for push")
							tempNode.pushChan <- pushChanData{"switch", currNode} // sends invalid id pull request to all nodes except currNode
						}
					}
				}
			}
		} else { // executes when the node becomes infected during push
			mu.Lock()                    // locks global lists
			pushTo := pickNode(currNode) // choose random node to push to
			if pushTo.status == false {
				fmt.Println(strconv.Itoa(currNode.id) + " picked node " + strconv.Itoa(pushTo.id))
				pushTo.pushChan <- pushChanData{currNode.msg, currNode} // send message through the receiving node's channel
				fmt.Println(strconv.Itoa(currNode.id) + " is waiting for confirmation")
				confirmation := <-currNode.pullChan // waits for confirmation from receiving node
				fmt.Println(strconv.Itoa(currNode.id) + " received confirmation " + strconv.Itoa(confirmation))
				pushTo.status = true // sets receiving node's status to infected
				fmt.Println(strconv.Itoa(pushTo.id) + " status changed to " + strconv.FormatBool(pushTo.status))
				statusList[pushTo.id] = true // tells array that node is infected node
				fmt.Println(strconv.Itoa(pushTo.id) + " status changed in global list to " +
					strconv.FormatBool(statusList[pushTo.id]))
			}
			mu.Unlock() // unlocks global lists
			fmt.Println(strconv.Itoa(pushTo.id) + " is infected! " + strconv.Itoa(10-sumBoolLocks(statusList)) +
				" left to infect.")
		}
	}
	fmt.Println("ended pull protocol")
	// switches to pull after half the nodes are infected
	pull(currNode)
}

// used to make goroutine that runs protocol based on input from user
func runNode(wg *sync.WaitGroup, currNode *Node, protocol string) {
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

	fmt.Println("random id is", randomId)
	pickedNode := &listOfNodes[randomId]
	fmt.Println("value of picked node is", pickedNode.id, "with status", pickedNode.status)
	return pickedNode
}

// returns integer representing how many entries in boolean array are true using locks within function
func sumBoolLocks(list [10]bool) int {
	mu.Lock() // locks global lists
	fmt.Println("checking sumBool")
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
	fmt.Println("checking sumBool")
	sum := 0
	for _, entry := range list {
		if entry {
			sum++
		}
	}
	return sum
}
