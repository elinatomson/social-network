import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';

function Followers() {
  const [followersUsers, setFollowersUsers] = useState([]);

  useEffect(() => {
    fetch('/followers')
      .then((response) => response.json())
      .then((data) => {
        if (data && data.followers_users) {
            setFollowersUsers(data.followers_users);
          }
      })
      .catch((error) => {
        console.error('Error fetching following users:', error);
      });
    }, []);

  return (
    <div>
      <div className="following">Followers</div>
      {followersUsers.length === 0 ? (
        <p className="user">No users.</p>
      ) : (
        <div className="user">
          {followersUsers.map((user) => (
            <Link className="link" key={user.user_id} to={`/user/${user.user_id}`}>
              {user.first_name} {user.last_name}
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
export default Followers;
