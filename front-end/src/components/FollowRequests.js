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
    }, []);

  return (
    <div>
        {followRequests.length > 0 && <div className="following">Follow requests:</div>}
        <div className="user">
        {followRequests.map((user) => (
            <div>
            <Link className="link" key={user.user_id} to={`/user/${user.user_id}`}>
            {user.first_name} {user.last_name}
            </Link>
            </div>
        ))}
        </div>
    </div>
  );
}
export default FollowRequests;
