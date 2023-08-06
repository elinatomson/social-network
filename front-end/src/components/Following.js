import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';

function Following() {
  const [followingUsers, setFollowingUsers] = useState([]);

  useEffect(() => {
    fetch('/following')
      .then((response) => response.json())
      .then((data) => {
        if (data && data.following_users) {
            setFollowingUsers(data.following_users);
          }
      })
      .catch((error) => {
        console.error('Error fetching following users:', error);
      });
    }, []);


  return (
    <div>
      <div className="following">Following</div>
      {followingUsers.length === 0 ? (
        <p className="user">No users.</p>
      ) : (
        <div className="user">
          {followingUsers.map((user) => (
            <div key={user.user_id}>
              <Link className="link" to={`/user/${user.user_id}`}>
                {user.first_name} {user.last_name}
              </Link>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
export default Following;
