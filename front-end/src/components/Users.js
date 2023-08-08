import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import Chat from './Chat';

function Users() {
  const [users, setUsers] = useState([]);
  const [selectedUser, setSelectedUser] = useState(null);


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
        if (selectedUser === user) {
          setSelectedUser(null); 
        } else {
          setSelectedUser(user);
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
                    <Link className="link" onClick={() => handleUserClick(user)}>
                    {user.first_name} {user.last_name}
                    </Link>
                </div>
            ))}
            </div>
        )}
        <div className="chat">
            {selectedUser && <Chat selectedUser={selectedUser} />}
        </div>
    </div>
  );
}
export default Users;
