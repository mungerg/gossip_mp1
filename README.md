## gossip_mp1
=====

Implementation of Gossip Protocol as described in http://publicatio.bibl.u-szeged.hu/1529/1/gossip11.pdf

Sebastian Fernandez, Gabby Munger, Justin Gomez

Description
-----
This program creates a network of independent nodes which spread a message via the gossip protocol, with the intention of using this program to gather data on runtime between the three different protocols (Push, Pull, and Push-Pull) and on systems of varying size.

How To Run
----
### 1.  Clone Git Repository
### 2.  Run using:
`go run main.go`

### 3.  Input the message you would like to send
The type of message(i.e. int, float, boolean, string) does not matter, as the program will read it as a string. There is strange behavior exhibited if the string includes a space, however the program still gives the desired output, as shown below.

<img width="538" alt="Screen Shot 2021-10-12 at 10 21 17 PM" src="https://user-images.githubusercontent.com/90423480/137057123-53356767-1bee-4c7c-bd6d-e85aa095390b.png">


### 4.  Input the specific gossip protocol to use.
There are three protocols to choose from: push, pull, and push-pull. As described by Jelasity:

`In push gossip, susceptible nodes are passive and infective nodes actively infect the population. In pull and push-pull gossip each node is active.`

The program will create a system of nodes that executes the desired gossip protocol to send a message between all of the nodes. The number of nodes in the system can be changed according to the instructions in the block comment at the top of main.go.

### 5. Example of expected output:

This example shows an input of _systems_ as the message, and _a_ to indicate the Push protocol in a system of 100 nodes. The output displays the time taken to complete the gossip protocol.
  
<img width="654" alt="Screen Shot 2021-10-12 at 10 17 57 PM" src="https://user-images.githubusercontent.com/90423480/137057178-73cfffd8-cb41-451f-a64b-dd854bfc5402.png">


This example shows an input of _message_ as the message, and _b_ to indicate the Pull protocol in a system of 100 nodes. The output displays the time taken to complete the gossip protocol.

<img width="652" alt="Screen Shot 2021-10-12 at 10 18 16 PM" src="https://user-images.githubusercontent.com/90423480/137057033-2379ac74-4073-4fbb-b5b1-f6b23843db5e.png">


This example shows an input of _hello_ as the message, and _c_ to indicate the Push-Pull protocol in a system of 20 nodes. The output displays the time taken to complete the gossip protocol.

<img width="672" alt="Screen Shot 2021-10-12 at 10 23 01 PM" src="https://user-images.githubusercontent.com/90423480/137056675-d8b9b553-a7b9-4ec7-9e43-42c317ebfce5.png">

Occasionally an error occurs when running the Push-Pull algorithm, resulting in a deadlock. This can be resolved by running the program again.


Workflow Diagram
----
![Gossip Protocol Workflow Diagram](https://github.com/mungerg/gossip_mp1/blob/main/Gossip%20Protocol.png)
