package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
)

var allConnections []*Connection
var allMessages []string

type Connection struct {
	name string
	ws   *websocket.Conn
}

func (connection *Connection) StartListening() {

	buf := make([]byte, 512)

	for {
		n, err := connection.ws.Read(buf[0:])
		if err != nil {
			ReleaseConnection(connection)
			PostMessage(fmt.Sprintf("(%s disconnected)", connection.name))
			break
		} else {
			HandleInputMessage(connection, buf[0:n])
		}
	}
}

func HandleNewConnection(ws *websocket.Conn) {

	log.Println("New connection!")

	connection := Connection{"Unnamed", ws}
	allConnections = append(allConnections, &connection)
	connection.StartListening()
}

func HandleInputMessage(from *Connection, data []byte) {

	var inputJson map[string]string
	json.Unmarshal(data, &inputJson)

	log.Println("New input json:", inputJson)

	var newMessageToAll string

	switch inputJson["action"] {

	case "post_message":
		newMessageToAll = fmt.Sprintf("%s: %s", from.name, inputJson["message"])

	case "update_name":
		oldName := from.name
		from.name = inputJson["name"]
		newMessageToAll = fmt.Sprintf("(%s -> %s)", oldName, from.name)

	case "get_history":
		byteMessage, _ := json.Marshal(allMessages)
		BroadcastTo(from, byteMessage)
	}

	if len(newMessageToAll) > 0 {
		PostMessage(newMessageToAll)
	}
}

func PostMessage(message string) {
	byteMessage, _ := json.Marshal([]string{message})
	BroadcastToAll(byteMessage)
	allMessages = append(allMessages, message)
}
