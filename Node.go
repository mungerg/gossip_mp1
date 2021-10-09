package gossip_mp1

type Node struct{
	id     int     // holds ID of node
	status bool // false corresponds to susceptible, true corresponds to infected
	msg    string  // holds message
	sendChan chan <- chanData
	receiveChan <- chan chanData
}

func createNode(id int, status bool, msg string) Node{
	send := make(chan chanData)
	receive := make(chan chanData)
	node := Node{ id, status, msg, send,receive}
	return node
	}
