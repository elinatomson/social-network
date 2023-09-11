import { displayErrorMessage } from "../components/ErrorMessage";
import { useNavigate } from "react-router-dom";
import React, { useState, useEffect } from "react";

function Follow({ userData, userId, pendingFollower }) {
  const [followContent, setFollowContent] = useState("");
  const [isPending, setIsPending] = useState(pendingFollower);
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
    };

    fetch('/follow', requestOptions)
      .then((response) => {
        if (response.ok) {
          setFollowContent("");
          setIsPending(!isPending); // Toggle isPending state
          navigate(`/user/${userId}`, { state: { followContent } });
        } else {
          return response.json();
        }
      })
      .catch((error) => {
        displayErrorMessage(`${error.message}`);
      });
  };

  useEffect(() => {
    setIsPending(pendingFollower); // Update the button text based on pendingFollower prop
  }, [pendingFollower]);

  return (
    <div>
      {isPending ? (
        <button className="follow-button" onClick={handleFollowUnfollow}>
          Unfollow
        </button>
      ) : (
        <button className="follow-button" onClick={handleFollowUnfollow}>
          Follow
        </button>
      )}
      <div id="error" className="alert"></div>
    </div>
  );
}

export default Follow;
