import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';

function FollowRequests() {
  const [followRequests, setFollowRequests] = useState([]);

  useEffect(() => {
    fetch('/follow-requests')
        .then((response) => response.json())
        .then((data) => {
            if (data === null) {
                setFollowRequests([]); 
            } else {
                setFollowRequests(data); 
            }
        })
        .catch((error) => {
        console.error('Error fetching follow requests:', error);
        });
    }, [followRequests]);

    const handleAccept = (followerID, followingID) => {
      const requestData = {
        follower_id: followerID,
        following_id: followingID,
        };

      const headers = new Headers();
      headers.append("Content-Type", "application/json");
  
      let requestOptions = {
        body: JSON.stringify(requestData),
        method: "POST",
        headers: headers,
      }

      fetch('/accept-follower', requestOptions)
        .then((response) => response.json())
        .catch((error) => {
          console.error('Error accepting follow request:', error);
      });
  };

  return (
    <div>
        {followRequests.length > 0 && <div className="following">Follow requests:</div>}
        <div className="user">
        {followRequests.map((user) => (
            <div key={user.user_id}>
            <Link className="link" to={`/user/${user.user_id}`}>
            {user.first_name} {user.last_name}
            </Link>
            <button className="accept-button" onClick={() => handleAccept(user.user_id, user.following_id)}>
              Accept
            </button>
            </div>
        ))}
        </div>
    </div>
  );
}
export default FollowRequests;
