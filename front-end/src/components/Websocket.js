import React, { useEffect, useState } from 'react';
import { displayErrorMessage } from "../components/ErrorMessage";

let socket;

function WebSocketComponent({ firstNameTo, firstNameFrom, conversationHistory }) {
  const [messages, setMessages] = useState(conversationHistory);
  console.log(conversationHistory)
  const [messageInput, setMessageInput] = useState('');

  useEffect(() => {
    function handleMessage(event) {
      const messageData = JSON.parse(event.data);
      console.log(messageData)
      const messageText = messageData.message;
      const senderName = messageData.first_name_from;
      const formattedTime = new Date(messageData.date).toLocaleString();
  
      const formattedMessage = `${formattedTime} - ${senderName}: ${messageText}`;
  
      setMessages((prevMessages) => [...prevMessages, formattedMessage]);
  
      // Automatically scroll to the bottom of the message box
      const messageBox = document.getElementById("message-box");
      messageBox.scrollTop = messageBox.scrollHeight;
    }

    socket = new WebSocket('ws://localhost:3000/ws');

    socket.addEventListener('open', () => {
      console.log('WebSocket connection established.');
    });

    socket.addEventListener('error', (error) => {
      console.error('WebSocket error:', error);
    });

    socket.addEventListener('message', handleMessage);

    return () => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.close();
        console.log('WebSocket connection closed.');
      }
    };
  }, [firstNameTo, firstNameFrom]);

  function handleSendMessage(event) {
    event.preventDefault();
    if (!socket || socket.readyState !== WebSocket.OPEN) {
      console.error('WebSocket connection not open.');
      return;
    }

    const date = new Date();
    const data = {
      message: messageInput,
      first_name_from: firstNameFrom,
      first_name_to: firstNameTo.first_name,
      date: date,
    };

    socket.send(JSON.stringify(data));
    setMessageInput('');

    fetch('/message', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    })
      .then((response) => {
        if (response.ok) {
          console.log(data)
          console.log('Message sent successfully');
        } else {
          return response.json();
        }
      })
      .then((errorMessage) => {
        if (errorMessage) {
          displayErrorMessage(errorMessage);
        }
      })
      .catch((error) => {
        displayErrorMessage(`An error occurred while sending message: ${error.message}`);
      });
  }

  return (
    <div>
      Chat with {firstNameTo.first_name}
      <div>
      <textarea className="chatbox" id="message-box" value={messages.join('\n')} readOnly />
      </div>
      <div className="container">
        <div className="left-container2">
          <input
            id="message-input"
            type="text"
            value={messageInput}
            onChange={(e) => setMessageInput(e.target.value)}
          />
        </div>
        <div className="right-container1">
          <button className="send-button" id="send-button" onClick={handleSendMessage}>
            Send
          </button>
        </div>
      </div>
    </div>
  );
}

export default WebSocketComponent