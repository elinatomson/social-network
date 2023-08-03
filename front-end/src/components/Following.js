import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';

function Following() {
  // Assuming you have already fetched the data and stored it in the followingUsers state
  const [followingUsers, setFollowingUsers] = useState([]);

  // Function to fetch following users from the server
  useEffect(() => {
    // Make an API call to fetch the list of following users
    // Replace 'API_ENDPOINT' with your actual backend API endpoint
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
