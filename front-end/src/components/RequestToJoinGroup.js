import { displayErrorMessage } from "../components/ErrorMessage";

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
        console.log(data)
        })
        .catch((error) => {
        displayErrorMessage(`An error occured while trying to follow this user: ${error.message}`);
        });
    };

    return (
        <button className="follow-button" onClick={handleJoinLeaveGroup}>
            Request To Join
        </button>
    )
}

export default RequestToJoinGroup;