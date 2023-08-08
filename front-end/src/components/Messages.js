import React, { useEffect, useState } from 'react';
import { WebSocketComponent } from './Websocket.js';

function UsersForChat() {
  const [users, setUsers] = useState([]);

  useEffect(() => {
    fetch('/users')
      .then((response) => response.json())
      .then((usersData) => {
        setUsers(usersData);
      })
      .catch((error) => {
        console.error('An error occurred while fetching users:', error.message);
      });
  }, []);

  function handleUserClick(user) {
    const formContainer = document.getElementById('formContainer');
    formContainer.innerHTML = `
      Chat with ${user}
      <div>
        <textarea id="message-box" rows="10" cols="50" readOnly></textarea>
        <div>
          <input type="text" id="message-input">
          <button id="send-button">Send</button>
        </div>
      </div>
      <p class="align">
        <input id="back" class="buttons" type="button" value="Back to main page">
      </p>
    `;


    WebSocketComponent(user);
  }

  function attachUserClickListeners() {
    const userItems = users.map((user) => (
      <div
        key={user.nickname}
        className={`user ${user.online ? 'online' : 'offline'}`}
        onClick={() => handleUserClick(user.nickname)}
      >
        {user.nickname}
      </div>
    ));

    return userItems;
  }

  return <div className="users-container">{attachUserClickListeners()}</div>;
}

export default UsersForChat;