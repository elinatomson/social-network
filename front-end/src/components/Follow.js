import { displayErrorMessage } from "../components/ErrorMessage";

function Follow ({ userData, userId }) {
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
        if (!response.ok) {
            return response.json().then((data) => {
            throw new Error(data.message);
            })
        } else {
            return response.json();
        }
        })
        .then((data) => {
        console.log(data)
        })
        .catch((error) => {
        displayErrorMessage(`An error occured while trying to follow this user: ${error.message}`);
        });
    };

    return (
        <button className="follow-button" onClick={handleFollowUnfollow}>
            Follow
        </button>
    )
}

export default Follow;