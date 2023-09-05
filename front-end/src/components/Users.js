import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import WebSocketComponent from './Websocket'; 
import { displayErrorMessage } from "../components/ErrorMessage";

function Users() {
  const [showUsers, setShowUsers] = useState(false);
  const [users, setUsers] = useState([]);
  const [firstNameTo, setFirstNameTo] = useState(null);
  const [firstNameFrom, setFirstNameFrom] = useState(null);

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
        {firstNameTo && <WebSocketComponent firstNameTo={firstNameTo} firstNameFrom={firstNameFrom}/>}
      </div>
    </div>
  );
}
export default Users;
