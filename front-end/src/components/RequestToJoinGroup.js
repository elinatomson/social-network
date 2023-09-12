import { displayErrorMessage } from "../components/ErrorMessage";
import React, { useState, useEffect } from "react";

function RequestToJoinGroup ({ groupId, pendingRequest }) {
    const [isPending, setIsPending] = useState(pendingRequest);

    const handleJoinLeaveGroup = () => {
        const requestData = {
            group_id: groupId, 
        };

        const headers = new Headers();
        headers.append("Content-Type", "application/json");

        let requestOptions = {
            body: JSON.stringify(requestData),
            method: "POST",
            headers: headers,
        }

        fetch('/request-to-join-group', requestOptions)
        .then((response) => {
        if (response.ok) {
            setIsPending(!isPending); 
        } else {
            return response.json();
        }
        })
        .catch((error) => {
            displayErrorMessage(`${error.message}`);
        });
    };

    useEffect(() => {
        setIsPending(pendingRequest); 
    }, [pendingRequest]);

    return (
        <div>
            {isPending ? (
                <button className="follow-button" onClick={handleJoinLeaveGroup}>
                    Cancel join request
                </button>
            ) : (
                <button className="follow-button" onClick={handleJoinLeaveGroup}>
                    Request To Join
                </button>
            )}
            <div id="error" className="alert"></div>
        </div>
    )
}

export default RequestToJoinGroup;