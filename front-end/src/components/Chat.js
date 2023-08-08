import React, { useState } from 'react';
import { WebSocketComponent } from './Websocket'; // Import the webSoc function

function Chat({ selectedUser }) {
  const [messageInput, setMessageInput] = useState('');

  // You can call the webSoc function here to initialize the chat

  return (
    <div className="chat">
      Chat with {selectedUser.first_name}
      <div>
        <textarea className="chatbox" value={messageInput} onChange={(e) => setMessageInput(e.target.value)} ></textarea>
        <div className="container">
            <div className="left-container2">
                <input type="text" id="message-input" value={messageInput} onChange={(e) => setMessageInput(e.target.value)}/>
            </div>
            <div className="right-container1">
                <button className="send-button">Send</button>
            </div>
        </div>
    </div>
    </div>
  );
}

export default Chat;
