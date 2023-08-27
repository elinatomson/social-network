import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import WebSocketComponent from './Websocket'; 
import { displayErrorMessage } from "../components/ErrorMessage";

function Users() {
  const [showUsers, setShowUsers] = useState(false);
  const [users, setUsers] = useState([]);
  const [messages, setMessages] = useState([]);
  const [firstNameTo, setFirstNameTo] = useState(null);

  const handleToggleUsers = () => {
    setShowUsers(!showUsers);
  };

  function fetchConversationHistory(user) {
    fetch(`/messages?firstNameTo=${user.first_name}`)
      .then(response => response.json())
      .then(messagesData => {  
        if (messagesData && messagesData.length > 0) {
          setMessages(messagesData);
        } else {
          setMessages([]);
        }
      })
      .catch(error => {
        displayErrorMessage(`${error.message}`);
      });
  }

  useEffect(() => {
    fetch('/users')
      .then((response) => response.json())
      .then((data) => {
        if (data) {
          const filteredUsers = data.filter((user) => !user.currentUser);
          setUsers(filteredUsers);
        }
      })
      .catch((error) => {
        displayErrorMessage(`${error.message}`);
      });
    }, []);

    const handleUserClick = (user) => {
        if (firstNameTo === user) {
          setFirstNameTo(null); 
        } else {
          setFirstNameTo(user);
          fetchConversationHistory(user);
        }
    };

  return (
    <div className="users">
      <div className="chat_users" onClick={handleToggleUsers}>All Social Network users for chat</div>
      {showUsers && (
      <div>
        <div id="error" className="alert"></div>
        {users.length === 0 ? (
          <p className="user">No users.</p>
        ) : (
          <div className="user">
            {users.map((user) => (
              <div key={user.user_id}>
                <Link className={`user ${user.online ? 'online' : 'offline'}`} onClick={() => handleUserClick(user)}>
                  {user.first_name} {user.last_name}
                </Link>
              </div>
            ))}
          </div>
        )}
      </div>
      )}

      <div className="chat">
        {firstNameTo && <WebSocketComponent firstNameTo={firstNameTo} conversationHistory={messages}/>}
      </div>
    </div>
  );
}
export default Users;
