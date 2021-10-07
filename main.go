package gossip_mp1

import (
	"fmt"
	"math/rand"
)

var listOfNodes []Node

type node struct {
	id     int     // holds ID of node
	status bool // false corresponds to susceptible, true corresponds to infected
	msg    string  // holds message
}
type chanData struct{
	status bool
	msg string
}
func main() {
	message, code := askInput()
	allInfected:= make(chan bool)
	//Need channel for receiving and sending updates. They are only one way?
	for i:=0; i<10;i++{
		send := make(chan chanData)
		receive := make(chan chanData)
		x := createNode(i,false,"Ready", send, receive)

		//This is going to stop running the Node once the forloop exits. So more of a while statement. Have a counter creating these nodes up until, it hits a limit, but keep the while loop going.
		go runNode(allInfected, x, code, message)
	}


}

// asks for user input of web addresses
func askInput() (string, string) {
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

	return message, printCode
}
func runNode( allInfected <-chan bool, currNode Node, protocol string, message string){
	lenOfList  := len(listOfNodes)
	//Create channels to and
	for true{
		if currNode.status == false{
			// reception chanData
			reception := <- currNode.receiveChan
			currNode.status = reception.status
			currNode.msg = reception.msg
		}

		//Gossiping
		for allInfected != true{
			//WAIT TIME

			//Pick a random node-This will probably need to be changed according to the protocol developed in chapter 4 that I havent read yet
			randomPeer := pickNode(currNode, lenOfList)
			if protocol == "Push" && currNode.status==true {
				push(randomPeer, currNode.sendChan, message)
			}
			if protocol =="Pull" && currNode.status ==false{
				pull(randomPeer,currNode.receiveChan,message)
			}

		}



	}

}

func pull(peer Node, receiveChan <-chan chanData, message string) {
	
}
/*func gossip( receiveChan  <- chan chanData, sendChan chan<- chanData,   allInfected <-chan bool, primeNode node, protocol string, message string){
	if(primeNode.status ==  false){

	}
	lenOfList  := len(listOfNodes)
	for allInfected != true{
		//WAIT TIME

		//Pick a random node-This will probably need to be changed according to the protocol developed in chapter 4 that I havent read yet
		randomPeer := pickNode(primeNode, lenOfList)
		if protocol == "Push" && primeNode.status==true {
			push(randomPeer, sendChan, message)
		}
		//Need something to do the updating for the receiver. Like right now you could just change the status of the receiver right here
		x := <- receiveChan

	}

}*/
func pickNode(primeNode Node, lenOfList int) Node {
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
func push(receiverNode Node, nodeChan chan <-chanData, message string){
	 toBeSent := chanData{true,message}
	 receiverNode.sendChan <- toBeSent //Ok I don't fully understand go channels then. Ah its less the sender is sending the message through the sendChannel and to the receive channel
	 //More the receiverNode is being sent the message through their own send channel.
	 
}



