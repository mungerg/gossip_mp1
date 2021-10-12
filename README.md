## gossip_mp1
=====

Implementation of Gossip Protocol as described in http://publicatio.bibl.u-szeged.hu/1529/1/gossip11.pdf

Sebastian Fernandez, Gabby Munger, Justin Gomez

Description
-----
This program creates a network of independent nodes which spread a message via the gossip protocol. 

How To Run
----
### 1.  Clone Git Repository
### 2.  Run using:
`go run main.go`

### 3.  Input the message you would like to send
The type of message does not matter because it will be read as a string.
### 4.  Input the specific gossip protocol to use.
There are three protocols to choose from: push, pull, and push-pull. As described by Jelasity:

`In push gossip, susceptible nodes are passive and infective nodes actively infect the population. In pull and push-pull gossip each node is active.`

Currently, while our code executes correctly for the pull protocol, it does not execute correctly for the push or the push-pull protocols. We are still working to change this.



Workflow Diagram
----
![Gossip Protocol Workflow Diagram](https://github.com/mungerg/gossip_mp1/blob/main/Gossip%20Protocol.png)
