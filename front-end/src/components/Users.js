import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import WebSocketComponent from '../components/Websocket'; 
import { displayErrorMessage } from "../components/ErrorMessage";
import Envelope from './../images/envelope.PNG';

function Users() {
  const [showUsers, setShowUsers] = useState(false);
  const [users, setUsers] = useState([]);
  const [firstNameTo, setFirstNameTo] = useState(null);
  const [firstNameFrom, setFirstNameFrom] = useState(null);
  const [unreadMessageCounts, setUnreadMessageCounts] = useState({});

  const handleToggleUsers = () => {
    setShowUsers(!showUsers);
  };

  useEffect(() => {
    fetch('/users')
      .then((response) => response.json())
      .then((data) => {
        if (data) {
          const filteredUsers = data.filter((user) => !user.currentUser);
          const currentUser = data.find((user) => user.currentUser);
          if (currentUser) {
            setFirstNameFrom(currentUser.first_name);
          }
          setUsers(filteredUsers);
          unreadMessageCount()
        }
      })
      .catch((error) => {
        displayErrorMessage(`${error.message}`);
      });
    }, []);

    function unreadMessageCount() {
      fetch(`/unread-messages`)
      .then(response => response.json())
      .then(messagesData => {  
        if (messagesData && messagesData.length > 0) {
          const counts = {};
          messagesData.forEach((message) => {
            const sender = message.first_name_from; 
            counts[sender] = (counts[sender] || 0) + 1;
          });
          setUnreadMessageCounts(counts);
        }
      })
      .catch(error => {
        displayErrorMessage(`${error.message}`);
      });
    }

    function messagesAsRead(firstNameFrom) {
      fetch(`/mark-messages-as-read/?firstNameFrom=${firstNameFrom}`)
      .then(response => {
        if (response.ok) {
          console.log('All messages marked as read.');
          setUnreadMessageCounts(prevCounts => ({
            ...prevCounts,
            [firstNameFrom]: "",
          }));
        }
      })
      .catch(error => {
        displayErrorMessage(`${error.message}`);
      });
    }

    const handleUserClick = (user) => {
        if (firstNameTo === user) {
          setFirstNameTo(null); 
        } else {
          setFirstNameTo(user);
          messagesAsRead(user.first_name);
        }
    };

    const handleChatClose = () => {
      setFirstNameTo(null); 
    };
  

  return (
    <div>
      <div className="chat_users" onClick={handleToggleUsers}>Open list of users to chat
        {Object.values(unreadMessageCounts).some(count => count > 0) && (
            <img className="envelope" src={Envelope} alt="envelope"></img>
          )}
      </div>
      {showUsers && (
      <div>
        <div id="error" className="alert"></div>
        {users.length === 0 ? (
          <p className="user">No users.</p>
        ) : (
          <div className="user">
            {users.map((user) => (
              <div key={user.user_id}>
                <Link className="link-btn" onClick={() => handleUserClick(user)}>
                  {user.first_name} {user.last_name}&nbsp;
                    {unreadMessageCounts[user.first_name] && (
                      <span className="unread">({unreadMessageCounts[user.first_name]} unread messages)</span>
                    )}
                </Link>
              </div>
            ))}
          </div>
        )}
      </div>
      )}
      <div className="chat">
        {firstNameTo && <WebSocketComponent firstNameTo={firstNameTo} firstNameFrom={firstNameFrom} closeChat={handleChatClose}/>}
      </div>
    </div>
  );
}
export default Users;
