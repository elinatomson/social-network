import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';

function GroupRequests() {
  const [groupRequests, setGroupRequests] = useState([]);

  useEffect(() => {
    fetch('/group-requests')
        .then((response) => response.json())
        .then((data) => {
            if (data === null) {
                setGroupRequests([]); 
            } else {
                setGroupRequests(data); 
                console.log(data)
            }
        })
        .catch((error) => {
        console.error('Error fetching group requests:', error);
        });
    }, []);

    if (!Array.isArray(groupRequests)) {
      return;
    }

    const handleAccept = (groupID, requesterID) => {
      const requestData = {
        group_id: groupID,
        requester_id: requesterID,
        };

      const headers = new Headers();
      headers.append("Content-Type", "application/json");
  
      let requestOptions = {
        body: JSON.stringify(requestData),
        method: "POST",
        headers: headers,
      }

      fetch('/accept-requester', requestOptions)
        .then((response) => response.json())
        .catch((error) => {
          console.error('Error accepting group request:', error);
      });
    }

      const handleDecline = (groupID, requesterID) => {
        const requestData = {
            group_id: groupID,
            requester_id: requesterID,
          };
  
        const headers = new Headers();
        headers.append("Content-Type", "application/json");
    
        let requestOptions = {
          body: JSON.stringify(requestData),
          method: "POST",
          headers: headers,
        }
  
        fetch('/decline-requester', requestOptions)
          .then((response) => response.json())
          .catch((error) => {
            console.error('Error declining group request:', error);
        });
      }

  return (
    <div>
        {groupRequests.length > 0 && <div className="following">Group requests:</div>}
        <div className="user">
        {groupRequests.map((user, index) => (
            <div className="requests" key={index}>
                <div className="container">
                <div className="left-container2">
                    Group:&nbsp;
                    <Link className="link" to={`/group/${user.group_id}`}>
                    {user.group_title}
                    </Link>
                    <br/>
                    Requester:&nbsp;
                    <Link className="link" to={`/user/${user.requester.user_id}`}>
                    {user.requester.first_name} {user.requester.last_name}
                    </Link>
                </div>
                <div className="right-container1">
                    <div className="container">
                    <div className="left-container2">
                        <button className="accept-button" onClick={() => handleAccept(user.user_id, user.group_id)}>
                        Accept
                        </button>
                    </div>
                    <div className="right-container1">
                        <button className="accept-button" onClick={() => handleDecline(user.user_id, user.group_id)}>
                        Decline
                        </button>
                    </div>
                    </div>
                </div>
                </div>
            <hr/>
            </div>
        ))}
        </div>
    </div>
  );
}

export default GroupRequests;
