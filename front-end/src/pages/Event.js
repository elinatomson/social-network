import React, { useState, useEffect, useCallback } from 'react';
import { useNavigate, useParams, Link } from 'react-router-dom';
import { displayErrorMessage } from "../components/ErrorMessage";
import Header from '../components/Header';
import Footer from "../components/Footer";

function Event() {
    const navigate = useNavigate();
    const [eventData, setEventData] = useState({});
    const { eventId } = useParams();
    const [updateEventData, setUpdateEventData] = useState(false);
    const [isGoing, setIsGoing] = useState(false);
    const [notGoing, setNotGoing] = useState(false);
    
    const token = document.cookie
    .split("; ")
    .find((row) => row.startsWith("sessionId="))
    ?.split("=")[1];

    const fetchEventData = useCallback(() => {
        fetch(`/group-event/${eventId}`, {
            headers: {
                Authorization: `${token}`,
            },
        })
        .then((response) => {
            if (!response.ok) {
                return response.json().then((data) => {
                    throw new Error(data.message);
                });
            } else {
                document.cookie = "sessionId=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/group-event/;";
                return response.json();
            }
        })
        .then((data) => {
            setEventData(data);
            if (data.going) {
                setIsGoing(true);
            } else if (data.not_going) {
                setNotGoing(true);
            } 
        })
        .catch((error) => {
            displayErrorMessage(`${error.message}`);
        });
    }, [eventId, token]);

    useEffect(() => {
        if (!token) {
          navigate("/login");
        } else {
          fetchEventData(updateEventData);
        }
    }, [navigate, token, updateEventData, fetchEventData]);
      

    const handleGoing = (participantID) => {
        const requestData = {
            event_id: Number(eventId),
            participant_id: participantID,
        };

        const headers = new Headers();
        headers.append("Content-Type", "application/json");

        let requestOptions = {
            body: JSON.stringify(requestData),
            method: "POST",
            headers: headers,
        }

        fetch('/going', requestOptions)
        .then((response) => response.json())
        .then(() => {
            setIsGoing(true); 
            setNotGoing(false); 
            setUpdateEventData((prev) => !prev); 
        })
        .catch((error) => {
            displayErrorMessage(`${error.message}`);
        });
    }

    const handleNotGoing = (participantID) => {
        const requestData = {
            event_id: Number(eventId),
            participant_id: participantID,
        };

        const headers = new Headers();
        headers.append("Content-Type", "application/json");

        let requestOptions = {
            body: JSON.stringify(requestData),
            method: "POST",
            headers: headers,
        }

        fetch('/not-going', requestOptions)
        .then((response) => response.json())
        .then(() => {
            setIsGoing(false);
            setNotGoing(true); 
            setUpdateEventData((prev) => !prev); 
        })
        .catch((error) => {
            displayErrorMessage(`${error.message}`);
        });
    }

    return (
        <div className="app-container">
            <Header />
            <div className="home">
            {eventData.is_group_member || eventData.is_group_creator ? (
            <div>
                <div id="error" className="alert"></div>
                <div className="container">
                    <div className="left-container">
                        <div>
                            <div className="group-data">Event creator</div>
                            {eventData.event ? (
                                <div className="user">
                                    <Link className="link" to={`/user/${eventData.event.user_id}`}>
                                        {eventData.event.first_name} {eventData.event.last_name}
                                    </Link>
                                </div>
                            ) : null}
                            <div className="group-data">Going</div>
                            {eventData.participants ? (
                                <div className="user">
                                    {eventData.participants.map((participant) => (
                                        participant.going ? (
                                            <div key={participant.participant_id}>
                                                <Link className="link" to={`/user/${participant.participant_id}`}>
                                                    {participant.first_name} {participant.last_name}
                                                </Link>
                                            </div>
                                        ) : null
                                    ))}
                                </div>                       
                            ) : null}
                            <div className="group-data">Not Going</div>
                            {eventData.participants ? (
                                <div className="user">
                                    {eventData.participants.map((participant) => (
                                        !participant.going ? (
                                            <div key={participant.participant_id}>
                                                <Link className="link" to={`/user/${participant.participant_id}`}>
                                                    {participant.first_name} {participant.last_name}
                                                </Link>
                                            </div>
                                        ) : null
                                    ))}
                                </div>                       
                            ) : null}
                        </div>
                    </div>
                    {eventData.event ? (
                        <div className="middle-container">
                            <div className="activity">
                            { eventData.event.title}
                            </div>
                            <p className="nothing">
                            {eventData.event.description}
                            </p>
                            <p className="nothing">
                            When: {eventData.event.time}
                            </p>
                            <div className="container">
                                <div className="left-container2">
                                    <button className="follow-button" onClick={() => handleGoing(eventData.event.event_id, eventData.event.participant_id)}     disabled={isGoing}>
                                    Going
                                    </button>
                                </div>
                                <div className="right-container1">
                                    <button className="follow-button" onClick={() => handleNotGoing(eventData.event.event_id, eventData.event.participant_id)}     disabled={notGoing}>
                                    Not going
                                    </button>
                                </div>
                            </div>
                        </div>
                    ) : null}
                    <div className="right-container">
                    </div>
                </div>
            </div>
            ) : null}
            </div>
        <Footer />
        </div>
    );
};  

export default Event;
