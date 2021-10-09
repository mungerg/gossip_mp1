package gossip_mp1

import (
	"fmt"
	"math/rand"
)

var listOfNodes []Node

type chanData struct {
	status bool
	msg    string
}

func main() {
	message, code := askInput()
	allInfected := make(chan bool)
	//Need channel for receiving and sending updates. They are only one way?
	var counter int = 1


	//This while loop is an attempt to just get the program constantly running.
	for i:=0;i<10;i++{
		if i==0{
			listOfNodes[i]= createNode(i,true,message)
		} else{
			listOfNodes[i]= createNode(i,false,"Ready")
		}

	}
			//Need to add node to list.
			//There is the potential for this go routine to only run during the if statement/ its iteration. Maybe I can call to something outside that will let it run independently.
			//This might also only do this once
			//Ah maybe I should make another function to run all the nodes at onces. Like there is actually no reason for the nodes to be created in this for loop
			//So like for i in listOFNodes: run each node. Probably be better, still need a way to ensure they are running, and not just exiting once the for loop closes.
			//That in this is the purpose of the boolean while loop, but there is probably a better way to implement it.
			//There is the potential for this go routine to only run during the if statement/ its iteration. Maybe I can call to something outside that will let it run independently.
	for i := 0; i < 10; i++ {
		go runNode(allInfected, listOfNodes[i],code,message)
	}

	complete := make(chan bool, 1)

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

	// I think it would be good to put some type of error handling if the input given is not correct
	// I am working to implement this - J

	return message, printCode
}
func runNode(allInfected <-chan bool, currNode Node, protocol string, message string) {
	lenOfList := len(listOfNodes)
	//Create channels to and
	var x bool = true

	for true {
		//Could be a while loop
		//Ensuring Nodes are passive. We don't want them to run gossip if they are passive.
		if currNode.status == false {
			// reception chanData
			reception := <-currNode.receiveChan
			currNode.status = reception.status
			currNode.msg = reception.msg
		}

		//Gossiping-How can I make it so that if it receives something in allInfected it breaks.
		for allInfected != true {
			//WAIT TIME
			
			//Pick a random node-This will probably need to be changed according to the protocol developed in chapter 4 that I havent read yet
			randomPeer := pickNode(currNode, lenOfList)
			if protocol == "Push" && currNode.status == true {
				push(randomPeer, message)
			}
			if protocol == "Pull" && currNode.status == false {
				pull(randomPeer, currNode.receiveChan, message)
			}
		}
	}
}

func pull(peer Node, receiveChan <-chan chanData, message string) {
	toBeRecieved := chanData(true, message)
	peer.receiveChan <- toBeRecieved
	// I might not be understanding the way the channels are setup.
	// Can I also access the current node and pull the update from the peer node instead?
}
func push(receiverNode Node, message string) {
	toBeSent := chanData{true, message}
	//The idea here is that the receiving Node is being sent this data through their receiving Channel, which right now is called send Channel.
	//This is just a quirk because I didn't fully understand channels when I wrote this.
	receiverNode.sendChan <- toBeSent

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
//This function picks a random node from the global list listOfNodes. If the node picks itself, it keeps running until it picks a node that isnt itself.
func pickNode(primeNode Node, lenOfList int) Node {
	var randomNode int
	x := false
	for x == false {
		randomNode := rand.Intn(lenOfList)
		if randomNode != primeNode.id {
			x = true
		}
	}
	receiver := listOfNodes[randomNode]
	return receiver
}
