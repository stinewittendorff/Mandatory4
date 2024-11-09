# Mandatory4

To run the program you have to start each participant in the TokenRing seperately, when starting a new paticipant you need to use a form of the following
 - go run TokenRing.go id, port, nextPort, hastoken
Here the id is the id for the created node, therefore the best choice is just to make it increment with one when running a new node.
The port is chosen by the user, when we ran it we started at 50051 but any port can be used.
The nextPort part is the port of the next part of the ring, in the case of starting at 50051 the nextPort would be localhost:50052. This is again incremented with each new node, except for the last node in the ring that should point back to our starting point.
The HasToken is a boolean that tells whether this node starts with the token, this should be set to false except for the last node added, as a true statement will start the ring.

The Nodes used in our logs where created as follows:
 - go run TokenRing.go 1 55001 localhost:55002 false
 - go run TokenRing.go 2 55002 localhost:55003 false
 - go run TokenRing.go 3 55003 localhost:55001 true
