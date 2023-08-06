import React, { useState, useEffect } from "react";
import { displayErrorMessage } from "./ErrorMessage";
import { useNavigate } from "react-router-dom"
import Search from "../components/Search";

function CreatePost() {
  const [postContent, setPostContent] = useState("");
  const [postPrivacy, setPostPrivacy] = useState("public");
  const [imageOrGif, setImageOrGif] = useState("");
  const [searchResults, setSearchResults] = useState(null);
  const [selectedUsers, setSelectedUsers] = useState([]);
  const [errors, setErrors] =useState([])

  const navigate = useNavigate();

  const token = document.cookie
  .split("; ")
  .find((row) => row.startsWith("sessionId="))
  ?.split("=")[1];

  useEffect(() => {
    if (!token) {
      navigate("/login");
    }
  }, [token, navigate]);

  const handleContentChange = (e) => {
    setPostContent(e.target.value);
  };

  const handlePrivacyChange = (e) => {
    setPostPrivacy(e.target.value);
  };

  const handleImageChange = (e) => {
    const file = e.target.files[0];
    setImageOrGif(file);
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

    let errors = []
      let required = [
        { field: postContent, name: "content"},
    ]

    required.forEach(function (obj) {
      if (obj.field === "") {
        errors.push(obj.name);
      }
    })

    setErrors(errors)

    if (errors.length > 0) {
      return;
    }

    // Convert the selectedUsers array to a comma-separated string
    const selectedUserIdString = selectedUsers.join(",");
    
    const postData = {
    content: postContent,
    privacy: postPrivacy,
    selected_user_id: selectedUserIdString,
    image: imageOrGif,
    };

    const headers = new Headers();
    headers.append("Content-Type", "application/json");
    headers.append("Authorization", token);

    let requestOptions = {
      body: JSON.stringify(postData),
      method: "POST",
      headers: headers,
    }

    fetch("/create-post", requestOptions)
      .then((response) => {
        if (response.ok) {
          setPostContent("");
          setPostPrivacy("public");
          setImageOrGif("");
          setSearchResults(null);
          setSelectedUsers([]);
          // Use the navigate function to redirect to the /main page and pass the post data in the state object
          navigate("/main", { state: { postContent } });
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
      <form onSubmit={handleSubmit}>
        <input className="content" placeholder="Post something..." value={postContent} onChange={handleContentChange} name="content"></input>
        {errors.includes("content") && (
          <p className="alert">Please fill in the input field.</p>
        )}
        <div className="container1">
            <div className="left-container2">
                <label htmlFor="privacy"></label>
                <select className="privacy" name="privacy" value={postPrivacy} onChange={handlePrivacyChange} required>
                    <option value="public">Public</option>
                    <option value="private">Private</option>
                    <option value="for-selected-users">For selected users</option>
                </select>
                <div className="post-search">
                {postPrivacy === "for-selected-users" && (
                  <Search setSearchResults={setSearchResults} />
                )}
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
                </div>
            </div>
            <div className="right-container1">
                <label htmlFor="image"></label>
                <input className="insert" type="file" name="image" accept="image/*, .gif" value={imageOrGif} onChange={handleImageChange}/>
            </div>
        </div>
        <div>
            <button className="button" type="submit">Create Post</button>
        </div>
      </form>
      <div id="error"></div>
    </div>
  );
}

export default CreatePost;
