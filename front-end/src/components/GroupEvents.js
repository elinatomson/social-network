import React, { useState, useEffect } from "react";
import { useNavigate, useLocation, Link } from 'react-router-dom';
import { displayErrorMessage } from "./ErrorMessage";

function GroupEvents({ groupId }) {
  const [allEvents, setAllEvents] = useState([]);
  const navigate = useNavigate();
  const location = useLocation();

  const { eventContent } = location.state || {};

  useEffect(() => {
    fetch(`/group-events?groupId=${groupId}`)
      .then((response) => response.json())
      .then((data) => {
        if (data) {
          setAllEvents(data);
        }
      })
      .catch((error) => {
        displayErrorMessage(`${error.message}`);
      });
  }, [groupId,navigate, eventContent]);

  return (
    <div>
      {allEvents.length === 0 ? (
        <p>No Events.</p>
      ) : (
        <div>
          {allEvents.map((event) => (
            <div key={event.event_id}>
              <Link className="link" to={`/group-event/${event.event_id}`}>
                <div>
                Event: <span>{event.title}</span>
                </div>
                <div>
                When: <span>{event.time}</span>
              </div>
              </Link>
              <p/>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default GroupEvents;
