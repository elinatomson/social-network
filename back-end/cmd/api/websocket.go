package main

import (
	"back-end/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader    = websocket.Upgrader{}
	connections = make(map[string]*websocket.Conn)
	mutex       = sync.Mutex{}
)

func (app *application) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	// in order to enable full-duplex communication and support WebSocket-specific features, the HTTP connection needs to be upgraded to a WebSocket connection.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}

	_, _, firstName, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the name"), http.StatusInternalServerError)
		return
	}
	// Add the WebSocket connection to the connections map
	//map is used to maintain active WebSocket connections.
	mutex.Lock()
	connections[firstName] = conn
	mutex.Unlock()
	//The function enters a loop to continuously read messages from the client.
	for {
		// Read message from the client
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message:", err)
			break
		}
		fmt.Print(message)
		//It then unmarshals the received message into a Message struct, which includes the recipient user's nickname.
		var msg models.Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Println("Failed to unmarshal message:", err)
			break
		}
		//Finally, it calls the handleMessage function, passing the recipient user's nickname, writer user's nickname and the message as parameters to handle the received message.
		app.handleMessage(r, w, msg.FirstNameFrom, msg.FirstNameTo, msg)
	}
	// Remove the WebSocket connection from the connections map when the connection is closed
	//The function uses a mutex to protect concurrent access to the connections map to ensure thread safety.
	mutex.Lock()
	delete(connections, firstName)
	mutex.Unlock()
}
func (app *application) handleMessage(r *http.Request, w http.ResponseWriter, receiverFirstName string, senderFirstName string, message models.Message) {
	// Check if the recipient user has an active WebSocket connection
	mutex.Lock()
	recipientConn, recipientFound := connections[receiverFirstName]
	mutex.Unlock()
	// Check if the sender user has an active WebSocket connection
	_, _, senderFirstName, _, err := app.database.DataFromSession(r)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("Failed to get the name"), http.StatusInternalServerError)
		return
	}
	mutex.Lock()
	senderConn, senderFound := connections[senderFirstName]
	mutex.Unlock()

	chatMessage := models.Message{
		Type:          "chat",
		Message:       message.Message,
		FirstNameFrom: senderFirstName,
		FirstNameTo:   receiverFirstName,
		Date:          message.Date,
	}

	if recipientFound {
		// Send the message to the recipient user's WebSocket connection
		data, err := json.Marshal(chatMessage)
		if err != nil {
			log.Println("Failed to marshal message:", err)
			return
		}
		err = recipientConn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("Failed to write message to recipient:", err)
		}
		fmt.Print(err)
	} else {
		log.Println("No active WebSocket connection found for recipient:", receiverFirstName)
	}
	if senderFound {
		// Send the message to the sender's WebSocket connection
		data, err := json.Marshal(chatMessage)
		if err != nil {
			log.Println("Failed to marshal message:", err)
			return
		}
		err = senderConn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("Failed to write message to sender:", err)
		}
	} else {
		log.Println("No active WebSocket connection found for sender:", senderFirstName)
	}
}
