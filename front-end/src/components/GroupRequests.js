import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { displayErrorMessage } from "./ErrorMessage";

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
            }
        })
        .catch((error) => {
          displayErrorMessage(`${error.message}`);
        });
    }, []);

    if (!Array.isArray(groupRequests)) {
      return null;
    }

    const handleAccept = (groupID, memberID) => {
      const requestData = {
        group_id: groupID,
        member_id: memberID,
      };

      const headers = new Headers();
      headers.append("Content-Type", "application/json");
  
      let requestOptions = {
        body: JSON.stringify(requestData),
        method: "POST",
        headers: headers,
      }

      fetch('/accept-group-request', requestOptions)
        .then((response) => {        
          if (response.ok) {
            setGroupRequests((prevRequests) =>
              prevRequests.filter((request) => request.group_id !== groupID)
          );
          } else {
            return response.json(); 
          }
        })
        .catch((error) => {
          displayErrorMessage(`${error.message}`);
      });
    }

      const handleDecline = (groupID, memberID) => {
        const requestData = {
            group_id: groupID,
            member_id: memberID,
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
          .then(() => {
            setGroupRequests((prevRequests) =>
              prevRequests.filter((request) => request.group_id !== groupID)
            );
          })
          .catch((error) => {
            displayErrorMessage(`${error.message}`);
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
                    <Link className="link" to={`/user/${user.member.user_id}`}>
                    {user.member.first_name} {user.member.last_name}
                    </Link>
                </div>
                <div className="right-container1">
                    <div className="container">
                    <div className="left-container2">
                        <button className="accept-button" onClick={() => handleAccept(user.group_id, user.member.user_id)}>
                        Accept
                        </button>
                    </div>
                    <div className="right-container1">
                        <button className="accept-button" onClick={() => handleDecline(user.group_id, user.member.user_id)}>
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
