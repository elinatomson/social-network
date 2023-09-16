import { displayErrorMessage } from "../components/ErrorMessage";
import React, { useState } from "react";

function Follow({ userData, userId, updateFollowers, isPublic }) {
  const [isPending, setIsPending] = useState("");
  const [isFollowing, setIsFollowing] = useState("");

  fetch(`/follower-check?userId=${userId}`)
    .then((response) => response.json())
    .then((data) => {
        if (data) {
            setIsPending(data.is_pending)
            setIsFollowing(data.is_following)
        }
    })
  .catch((error) => {
      displayErrorMessage(`${error.message}`);
  });

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
                setIsPending((prevIsPending) => !prevIsPending);
                setIsFollowing((prevIsFollowing) => !prevIsFollowing);
                if (isPublic) {
                    updateFollowers((prevFollowers) =>
                      isPending ? prevFollowers.concat(userId) : prevFollowers.filter((id) => id !== userId)
                    );
                }
            } else {
            return response.json();
            }
        })
        .catch((error) => {
            displayErrorMessage(`${error.message}`);
        });
    };
    
    return (
        <div>
          {isPending || (!isPending && isFollowing) ? (
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
