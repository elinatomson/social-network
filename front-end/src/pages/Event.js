import React, { useState, useEffect } from 'react';
import { useNavigate, useParams, Link } from 'react-router-dom';
import { displayErrorMessage } from "../components/ErrorMessage";
import Header from '../components/Header';
import Footer from "../components/Footer";

function Event() {
    const navigate = useNavigate();
    const [eventData, setEventData] = useState({});
    const { eventId } = useParams();
    //const [isMember, setIsMember] = useState(false);
    
    const token = document.cookie
    .split("; ")
    .find((row) => row.startsWith("sessionId="))
    ?.split("=")[1];

    useEffect(() => {
        if (!token) {
        navigate("/login");
        } else {
        fetch(`/group-event/${eventId}`, {
            headers: {
            Authorization: `${token}`,
            },
        })
        .then((response) => {
            if (!response.ok) {
            return response.json().then((data) => {
                throw new Error(data.message);
            })
            } else {
            return response.json();
            }
        })
        .then((data) => {
            setEventData(data);
            console.log(data)
            //const currentUserID = data.userID;
            //const groupMembers = data.group_members || [];
            //const groupCreator = data.group.user_id;
            //setIsMember(groupMembers.includes(currentUserID) || currentUserID === groupCreator);
        })
        .catch((error) => {
            displayErrorMessage(`An error occured while displaying group: ${error.message}`);
        });
        }
    }, [navigate, eventId, token]);

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
            .catch((error) => {
                console.error('Error with going:', error);
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
            .catch((error) => {
                console.error('Error with not going:', error);
        });
    }


    return (
        <div className="app-container">
            <Header />
            <div className="home">
            <div>
                <div id="error" className="alert"></div>
                <div className="container">
                    <div className="left-container">
                        <div className="users">
                            <div className="following">Event creator</div>
                            {eventData.event ? (
                                <div className="user">
                                    <Link className="link" to={`/user/${eventData.event.user_id}`}>
                                        {eventData.event.first_name} {eventData.event.last_name}
                                    </Link>
                                </div>
                            ) : null}
                            <div className="following">Going</div>
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
                            <div className="following">Not Going</div>
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
                            <div className="container">
                                <div className="left-container2">
                                    <button className="follow-button" onClick={() => handleGoing(eventData.event.event_id, eventData.event.participant_id)}>
                                    Going
                                    </button>
                                </div>
                                <div className="right-container1">
                                    <button className="follow-button" onClick={() => handleNotGoing(eventData.event.event_id, eventData.event.participant_id)}>
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
            </div>
        <Footer />
        </div>
    );
};  

export default Event;
