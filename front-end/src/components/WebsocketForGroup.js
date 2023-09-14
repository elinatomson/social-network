import React, { useEffect, useState, useRef } from 'react';
import { displayErrorMessage } from "../components/ErrorMessage";
import Picker from '@emoji-mart/react'
import data from '@emoji-mart/data/sets/14/twitter.json'

function WebSocketComponentForGroup({ groupName, firstNameFrom }) {
  const [messages, setMessages] = useState([]);
  const [messageInput, setMessageInput] = useState('');
  const [ws, setWs] = useState(null);
  const [showEmojiPicker, setShowEmojiPicker] = useState(false);

  useEffect(() => {
    const websocket = new WebSocket(`ws://localhost:8080/chatroom/?group=${groupName}`);
    websocket.onopen = () => {
      console.log('WebSocket connected');
      setWs(websocket);
      fetchConversationHistory();
    };

    websocket.onmessage = (event) => {
      const eventData = JSON.parse(event.data);
      const message = eventData.message;
      const from = eventData.first_name_from;
      const to = eventData.first_name_to;
      handleMessage(message, from, to)
    };

    websocket.onclose = () => {
      console.log('WebSocket disconnected');
    };

    return () => {
      websocket.close();
    };

    function fetchConversationHistory() {
        fetch(`/group-conversation-history/?groupName=${groupName}`)
        .then(response => response.json())
        .then(messagesData => { 
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

  }, [groupName, firstNameFrom]);
  
  const handleInputChange = (event) => {
    setMessageInput(event.target.value);
  };

  const handleKeyDown = (event) => {
    if (event.keyCode === 13) {
      sendMessage()
    }
  }
  
  const chatParent = useRef(null);

  useEffect(() => {
    const domNode = chatParent.current;
    if (domNode) {
      domNode.scrollTop = domNode.scrollHeight;
    }
  })

  function handleMessage(message, from, to) {
    const senderName = from;
    const messageText = message;
    const date = new Date();
    const formattedTime = new Date(date).toLocaleString();
    const formattedMessage = `${formattedTime} - ${senderName}: ${messageText}`;
    setMessages((prevMessages) => [...prevMessages, formattedMessage]);  
  }

  const sendMessage = () => {
    if (messageInput.trim() !== '') {
    const data = {
      message: messageInput,
      first_name_from: firstNameFrom,
      first_name_to: groupName,
      date: new Date(),
    };
    ws.send(JSON.stringify(data));
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
        displayErrorMessage(`${error.message}`);
      });
    }
  };

  return (
    <div>
      <div className="chat-messages" ref={chatParent}>
          {messages.map((msg, index) => (
            <div className="chat-message" key={index}>{msg}</div>
          ))}
        </div>
      <div className="chat-container">
        <div className="left-container3">
          <input
            id="message-input"
            type="text"
            className="input-box"
            value={messageInput}
            onChange={handleInputChange}
            onKeyDown={(e) => handleKeyDown(e) }
          />
          <button
            className="emoji-button"
            onClick={() => setShowEmojiPicker(!showEmojiPicker)}
          >
            😃
          </button>
        </div>
        <div className="right-container1">
          <button className="chat-send-button" id="send-button" onClick={sendMessage}>
            Send
          </button>
        </div>
        <div className="emoji-picker">
          {showEmojiPicker && (
            <Picker data={data}
              set="twitter" 
              onEmojiSelect={(emoji) => {
                setMessageInput((prevMessageInput) => prevMessageInput + emoji.native);
                setShowEmojiPicker(false);
              }}
            />
          )}
        </div>
      </div>
    </div>
  );
}

export default WebSocketComponentForGroup