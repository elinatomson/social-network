import React, { useEffect, useState, useCallback } from 'react';
import { displayErrorMessage } from "../components/ErrorMessage";

let socket;

export function WebSocketComponent({ firstNameTo, firstNameFrom }) {
  const [messages, setMessages] = useState([]);
  const [messageInput, setMessageInput] = useState('');

  const handleMessage = useCallback((message) =>  {
    let senderName = firstNameTo;
    if (message.firstnameto === firstNameTo) {
      senderName = firstNameFrom;
    }
    const messageText = message.message;
    const formattedTime = new Date(message.date).toLocaleString();
    setMessages((prevMessages) => [
        ...prevMessages,
        `${formattedTime} - ${senderName}: ${messageText}`,
        ]);
    }, [firstNameTo, firstNameFrom]);

  useEffect(() => {
    socket = new WebSocket('ws://localhost:8080/ws');

    socket.addEventListener('open', () => {
      console.log('WebSocket connection established.');
    });

    socket.addEventListener('error', (error) => {
      console.error('WebSocket error:', error);
    });

    socket.addEventListener('message', (event) => {
      const message = JSON.parse(event.data);
      handleMessage(message);
    });

    return () => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.close();
        console.log('WebSocket connection closed.');
      }
    };
  }, [handleMessage]);

  function handleSendMessage(event) {
    event.preventDefault();
    if (!socket || socket.readyState !== WebSocket.OPEN) {
      console.error('WebSocket connection not open.');
      return;
    }

    const date = new Date();
    const data = {
      message: messageInput,
      nicknamefrom: firstNameFrom,
      nicknameto: firstNameTo,
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
          console.log('Message sent successfully');
        } else {
          return response.text();
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
      <textarea id="message-box" value={messages.join('\n')} readOnly />
      <input
        id="message-input"
        type="text"
        value={messageInput}
        onChange={(e) => setMessageInput(e.target.value)}
      />
      <button id="send-button" onClick={handleSendMessage}>
        Send
      </button>
    </div>
  );
}
