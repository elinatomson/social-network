import React, { useState } from 'react';
import { displayErrorMessage } from "./ErrorMessage";
import Search from "../components/Search";

function InviteNewMember({ groupId }) {
  const [searchResults, setSearchResults] = useState(null);
  const [selectedUser, setSelectedUser] = useState(null);

  const handleUserSelection = (userId) => {
    setSelectedUser(userId) 
  };

  const handleSubmit = (e) => {
    e.preventDefault();

  const invitationData = {
    group_id: groupId, 
    member_id: selectedUser,
    };

    console.log(invitationData)

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
          setSelectedUser(null);
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
    }

  return (
    <form onSubmit={handleSubmit}>
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
                      type="radio"
                      name="selected_user_id"
                      value={result.user_id}
                      onChange={() => handleUserSelection(result.user_id)}
                      checked={selectedUser === result.user_id}
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
  </form>
  );
}

export default InviteNewMember;
