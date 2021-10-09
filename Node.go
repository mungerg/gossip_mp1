package gossip_mp1

type Node struct{
	id     int     // holds ID of node
	status bool // false corresponds to susceptible, true corresponds to infected
	msg    string  // holds message
	sendChan chan <- chanData
	receiveChan <- chan chanData
}

func createNode(id int, status bool, msg string,sendChan chan <- chanData, receiveChan <- chan chanData) Node{
	node := Node{ id, status, msg, sendChan,receiveChan}
	return node
	}
