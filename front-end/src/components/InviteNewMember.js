import React, { useState } from 'react';
import { displayErrorMessage } from "./ErrorMessage";
import Search from "../components/Search";

function InviteNewMember() {
  const [searchResults, setSearchResults] = useState(null);
  const [selectedUsers, setSelectedUsers] = useState([]);

  const handleUserSelection = (userId) => {
    setSelectedUsers((prevSelectedUsers) => {
    if (prevSelectedUsers.includes(userId)) {
        // User is already selected, so remove from the selection
        return prevSelectedUsers.filter((id) => id !== userId);
    } else {
        // User is not selected, so add to the selection
        return [...prevSelectedUsers, userId];
    }
    });
  };

  const handleSubmit = (e) => {
    e.preventDefault();
  }

  // Convert the selectedUsers array to a comma-separated string
  const selectedUserIdString = selectedUsers.join(",");

  const invitationData = {
    selected_user_id: selectedUserIdString,
    };

    const headers = new Headers();
    headers.append("Content-Type", "application/json");

    let requestOptions = {
      body: JSON.stringify(invitationData),
      method: "POST",
      headers: headers,
    }

    fetch("/invite", requestOptions)
      .then((response) => {
        if (response.ok) {
          setSearchResults(null);
          setSelectedUsers([]);
        } else {
          return response.json(); 
        }
      })
      .then((errorMessage) => {
        if (errorMessage) {
          displayErrorMessage(errorMessage.error); 
        }
      })
      .catch((error) => {
        displayErrorMessage(`${error.message}`);
      });

  return (
    <div onSubmit={handleSubmit}>
      <div className="following">
        Invite new member:
      </div>
      <Search setSearchResults={setSearchResults} />
      <div className="search-results">
          {searchResults !== null && searchResults.length > 0 && (
          searchResults.map((result) => (
              <div key={result.user_id} className="search-result-item">
                  <label htmlFor="selected_user_id"></label>
                  <input
                      type="checkbox"
                      name="selected_user_id"
                      value={result.user_id}
                      onChange={() => handleUserSelection(result.user_id)}
                      checked={selectedUsers.includes(result.user_id)}
                  />
                  {result.first_name} {result.last_name}
              </div>
              ))
          )}
          {searchResults !== null && searchResults.length === 0 && (
          <p>No results found for your search query.</p>
          )}
      </div>
    <div id="error" className="alert"></div>
    <button className="log-out-button" type="submit">Invite</button>
  </div>
  );
}

export default InviteNewMember;
