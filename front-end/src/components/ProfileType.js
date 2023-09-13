import { displayErrorMessage } from "../components/ErrorMessage";
import React, { useState, useEffect } from "react";

function Follow({ profileType }) {
  const [isPublic, setIsPublic] = useState(profileType);

  const handleProfileType = () => {

    const headers = new Headers();
    headers.append("Content-Type", "application/json");

    let requestOptions = {
      method: "POST",
      headers: headers,
    };

    fetch('/profile-type', requestOptions)
        .then((response) => {
            if (response.ok) {
                setIsPublic(!isPublic); 
            } else {
            return response.json();
            }
        })
        .catch((error) => {
            displayErrorMessage(`${error.message}`);
        });
    };

    useEffect(() => {
        setIsPublic(profileType); 
    }, [profileType]);

    return (
        <div>
            {(isPublic) ? (
                <button className="follow-button" onClick={handleProfileType}>
                Set your profile as private
                </button>
            ) : (
                <button className="follow-button" onClick={handleProfileType}>
                Set your profile as public
                </button>
            )}
            <div id="error" className="alert"></div>
        </div>
    );
}

export default Follow;
