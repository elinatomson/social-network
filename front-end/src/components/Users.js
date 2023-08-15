import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import WebSocketComponent from './Websocket'; 
import { displayErrorMessage } from "../components/ErrorMessage";

function Users() {
  const [users, setUsers] = useState([]);
  const [messages, setMessages] = useState([]);
  const [firstNameTo, setFirstNameTo] = useState(null);

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
        displayErrorMessage('An error occurred while displaying messages: ' + error.message);
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
        console.error('Error fetching following users:', error);
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
        <div className="following">All Social Network users</div>
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
        <div className="chat">
            {firstNameTo && <WebSocketComponent firstNameTo={firstNameTo} conversationHistory={messages}/>}
        </div>
    </div>
  );
}
export default Users;
