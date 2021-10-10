package main

type Node struct{
	id     int     // holds ID of node
	status bool // false corresponds to susceptible, true corresponds to infected
	msg    string  // holds message
	pullChan chan int // for receiving pull requests
	pushChan chan string // for pushing message through
}

func createNode(id int, status bool, msg string) Node{
	//Why is this one an int? Is it choosing based off id?
	pull := make(chan int)
	push := make(chan string)
	node := Node{ id, status, msg, pull,push}
	return node
	}
