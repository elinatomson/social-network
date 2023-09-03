import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { displayErrorMessage } from "../components/ErrorMessage";

function GroupEventNotifications() {
  const [eventNotifications, setEventNotifications] = useState([]);

  useEffect(() => {
    fetch('/group-event-notifications')
        .then((response) => response.json())
        .then((data) => {
            if (data === null) {
                setEventNotifications([]); 
            } else {
                setEventNotifications(data); 
            }
        })
        .catch((error) => {
          displayErrorMessage(`${error.message}`);
        });
    }, []);

    if (!Array.isArray(eventNotifications)) {
      return null;
    }

    const handleClick = (eventID) => {
        const requestData = {
            event_id: eventID,
        };
        const headers = new Headers();
        headers.append("Content-Type", "application/json");
    
        let requestOptions = {
            body: JSON.stringify(requestData),
            method: "POST",
            headers: headers,
        }

        fetch('/group-event-seen', requestOptions)
        .then((response) => {        
            if (response.ok) {
            setEventNotifications((prevInvitation) =>
                prevInvitation.filter((invitation) => invitation.event_id !== eventID)
            );
            } else {
            return response.json(); 
            }
        })
        .catch((error) => {
            displayErrorMessage(`${error.message}`);
        });
    }

  return (
    <div>
        {eventNotifications.length > 0 && <div className="following">Event notifications:</div>}
        <div id="error" className="alert"></div>
        <div className="user">
        {eventNotifications.map((event, index) => (
            <div className="requests" key={index}>
                <div>
                    New event&nbsp;
                    <Link className="link" to={`/group-event/${event.event_id}`} onClick={() => handleClick(event.event_id)}>
                    {event.event_title}
                    </Link>&nbsp;
                    in group&nbsp;
                    <Link className="link" to={`/group/${event.group_id}`}>
                    {event.group_title}
                    </Link>
                    <br/>
                </div>
            <hr/>
            </div>
        ))}
        </div>
    </div>
  );
}

export default GroupEventNotifications;
