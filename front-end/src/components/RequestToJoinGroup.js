import { displayErrorMessage } from "../components/ErrorMessage";
import { displayMessage } from "../components/ErrorMessage";

function RequestToJoinGroup ({ groupId }) {
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
        if (!response.ok) {
            return response.json().then((data) => {
            throw new Error(data.message);
            })
        } else {
            return response.json();
        }
        })
        .then((data) => {
            displayMessage(`Your request has been sent to the group owner`);
        })
        .catch((error) => {
            displayErrorMessage(`${error.message}`);
        });
    };

    return (
        <div>
            <button className="follow-button" onClick={handleJoinLeaveGroup}>
                Request To Join
            </button>
            <div id="message"></div>
            <div id="error" className="alert"></div>
        </div>
    )
}

export default RequestToJoinGroup;