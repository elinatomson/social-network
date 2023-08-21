import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';

function GroupInvitations() {
  const [groupInvitations, setGroupInvitations] = useState([]);

  useEffect(() => {
    fetch('/group-invitations')
        .then((response) => response.json())
        .then((data) => {
            if (data === null) {
                setGroupInvitations([]); 
            } else {
                setGroupInvitations(data); 
                console.log(data)
            }
        })
        .catch((error) => {
        console.error('Error fetching group requests:', error);
        });
    }, []);

    if (!Array.isArray(groupInvitations)) {
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

      fetch('/accept-group-invitation', requestOptions)
        .then((response) => response.json())
        .catch((error) => {
          console.error('Error accepting group invitation:', error);
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
  
        fetch('/decline-group-request', requestOptions)
          .then((response) => response.json())
          .catch((error) => {
            console.error('Error declining group request:', error);
        });
      }

  return (
    <div>
        {groupInvitations.length > 0 && <div className="following">Group invitations:</div>}
        <div className="user">
        {groupInvitations.map((user, index) => (
            <div className="requests" key={index}>
                <div className="container">
                <div className="left-container2">
                    Group:&nbsp;
                    <Link className="link" to={`/group/${user.group_id}`}>
                    {user.group_title}
                    </Link>
                    <br/>
                </div>
                <div className="right-container1">
                    <div className="container">
                    <div className="left-container2">
                        <button className="accept-button" onClick={() => handleAccept(user.group_id, user.invited_user.user_id)}>
                        Accept
                        </button>
                    </div>
                    <div className="right-container1">
                        <button className="accept-button" onClick={() => handleDecline(user.group_id, user.invited_user.user_id)}>
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

export default GroupInvitations;
