package main

import (
	"log"
)

func ReleaseConnection(c *Connection) {

	log.Println("ReleaseConnection:", c.name)

	index := -1
	for k, v := range allConnections {
		if c == v {
			index = k
			break
		}
	}

	if index >= 0 {
		allConnections = append(allConnections[:index], allConnections[index+1:]...)
	}

	log.Println("Connections count:", len(allConnections))
}

func BroadcastToAll(message []byte) {
	for _, connection := range allConnections {
		BroadcastTo(connection, message)
	}
}

func BroadcastTo(to *Connection, message []byte) {
	to.ws.Write(message)
}
