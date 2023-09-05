import React, { useEffect, useState } from 'react';
import { displayErrorMessage } from "../components/ErrorMessage";

let socket;

function WebSocketComponent({ firstNameTo, firstNameFrom }) {
  const [messages, setMessages] = useState([]);
  const [messageInput, setMessageInput] = useState('');

  useEffect(() => {
    function fetchConversationHistory() {
      fetch(`/conversation-history/?firstNameTo=${firstNameTo.first_name}`)
      .then(response => response.json())
      .then(messagesData => {  
        console.log(messagesData)
        if (messagesData && messagesData.length > 0) {
          const messages = messagesData.map(message => {
            var formattedDate = new Date(message.date).toLocaleString();
            return `${formattedDate} - ${message.first_name_from}: ${message.message}`;
          });
          setMessages(messages);
        } else {
          setMessages([]);
        }
      })
      .catch(error => {
        displayErrorMessage(`${error.message}`);
      });
    }

    function handleMessage(event) {
      const messageData = JSON.parse(event.data);
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
      fetchConversationHistory({ firstNameTo })
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
  }, [firstNameTo]);

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

    const formattedTime = date.toLocaleString();
    const formattedMessage = `${formattedTime} - ${firstNameFrom}: ${messageInput}`;
    setMessages((prevMessages) => [...prevMessages, formattedMessage]);

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