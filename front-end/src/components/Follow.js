import { displayErrorMessage } from "../components/ErrorMessage";
import { useNavigate } from "react-router-dom";
import React, { useState } from "react";

function Follow ({ userData, userId }) {
    const [followContent, setFollowContent] = useState("");
    const navigate = useNavigate();

    const handleFollowUnfollow = () => {
        const followData = {
            following_id: userId, 
            request_pending: userData,
        };

        const headers = new Headers();
        headers.append("Content-Type", "application/json");

        let requestOptions = {
        body: JSON.stringify(followData),
        method: "POST",
        headers: headers,
        }

        fetch('/follow', requestOptions)
        .then((response) => {
            if (response.ok) {
                setFollowContent("");
                navigate(`/user/${userId}`, { state: { followContent } });
            } else {
                return response.json();
            }
        })
        .catch((error) => {
            displayErrorMessage(`${error.message}`);
        });
    };

    return (
        <button className="follow-button" onClick={handleFollowUnfollow}>
            Follow
        </button>
    )
}

export default Follow;