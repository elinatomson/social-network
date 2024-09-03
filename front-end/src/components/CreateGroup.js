import React, { useState } from "react";
import { displayErrorMessage } from "./ErrorMessage";
import { useNavigate } from "react-router-dom"
import Search from "../components/Search";

function CreateGroup() {
    const [showFields, setShowFields] = useState(false);
    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [searchResults, setSearchResults] = useState(null);
    const [selectedUsers, setSelectedUsers] = useState([]);
    const [errors, setErrors] =useState([])
    const [groupContent, setGroupContent] = useState("");

    const navigate = useNavigate();

    const handleTitleChange = (e) => {
        setTitle(e.target.value);
      };
    
    const handleDescriptionChange = (e) => {
        setDescription(e.target.value);
    };

    const handleToggleFields = () => {
        setShowFields(!showFields);
    };

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

        setErrors([]);

        const newErrors = [];
    
        if (title.trim() === "") {
          newErrors.push("title");
        } else if (title.length > 10) {
          newErrors.push("title_length");
        }

        if (description.trim() === "") {
          newErrors.push("description");
        } else if (description.length > 100) {
          newErrors.push("description_length");
        }
    
        setErrors(newErrors)
    
        if (newErrors.length > 0) {
          return;
        }

    const selectedUserIdString = selectedUsers.join(",");
    
    const groupData = {
    title: title,
    description: description,
    selected_user_id: selectedUserIdString,
    };

    const headers = new Headers();
    headers.append("Content-Type", "application/json");

    let requestOptions = {
      body: JSON.stringify(groupData),
      method: "POST",
      headers: headers,
    }

    fetch("/create-group", requestOptions)
      .then((response) => {
        if (response.ok) {
          setGroupContent("");
          setTitle("");
          setDescription("");
          setSearchResults(null);
          setSelectedUsers([]);
          setShowFields(false);
          navigate(`/main`, { state: { groupContent } });
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
  };

  return (
    <div className="posting">
        <h2 className="center" onClick={handleToggleFields}>New group</h2>
        {showFields && (
            <form className="login-form" onSubmit={handleSubmit}>
                <input value={title} onChange={handleTitleChange} placeholder="Title" name="title"/>
                {errors.includes("title") && (
                    <p className="alert">Please fill in the title.</p>
                )}
                {errors.includes("title_length") && (
                  <p className="alert">Title is too long (max 10 characters).</p>
                )}
                <input value={description} onChange={handleDescriptionChange} placeholder="Description" name="description"/>
                {errors.includes("description") && (
                    <p className="alert">Please fill in the description.</p>
                )}
                {errors.includes("description_length") && (
                  <p className="alert">Description is too long (max 100 characters).</p>
                )}
                <div>
                    <Search setSearchResults={setSearchResults} />
                    <div className="search-results">
                        {searchResults !== null && searchResults.length > 0 && (
                        searchResults.map((result) => (
                            <div key={result.user_id} className="search-result-item">
                                <label htmlFor="selected_user_id"></label>
                                <input
                                    class="checkbox"
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
                </div>
                <div id="error" className="alert"></div>
                <button className="button" type="submit">Create group</button>
            </form>
        )}
    </div>
  );
}

export default CreateGroup;
