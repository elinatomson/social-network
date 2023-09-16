import React, { useState } from "react";
import { displayErrorMessage } from "./ErrorMessage";
import { useNavigate } from "react-router-dom"
import Search from "../components/Search";

function CreatePost({ groupId }) {
  const [showFields, setShowFields] = useState(false);
  const [postContent, setPostContent] = useState("");
  const [postPrivacy, setPostPrivacy] = useState("public");
  const [imageOrGif, setImageOrGif] = useState(null);
  const [searchResults, setSearchResults] = useState(null);
  const [selectedUsers, setSelectedUsers] = useState([]);
  const [errors, setErrors] =useState([])

  const navigate = useNavigate();

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

    const selectedUserIdString = selectedUsers.join(",");
    
    const postData = new FormData();
    postData.append("content", postContent);
    postData.append("privacy", postPrivacy);
    postData.append("selected_user_id", selectedUserIdString);
    postData.append("image", imageOrGif);
    postData.append("group_id", groupId);

    const headers = new Headers();

    let requestOptions = {
      body: postData,
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
          setShowFields(false);
          //using the navigate function to redirect to the /main page and pass the post data in the state object
          if (groupId == null) {
            navigate(`/main`, { state: { postContent } });
          } else {
            navigate(`/group/${groupId}`, { state: { postContent } });
          }
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
      <h2 className="center" onClick={handleToggleFields}>New post</h2>
      {showFields && (
        <form className="login-form" onSubmit={handleSubmit}>
          <input placeholder="Post something..." value={postContent} onChange={handleContentChange} name="content"></input>
          {errors.includes("content") && (
            <p className="alert">Please fill in the input field.</p>
          )}
          <div className="container">
          {groupId ? null : (
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
          )}
              <div className="right-container1">
                  <label htmlFor="image"></label>
                  <input className="insert" type="file" name="image" accept="image/*, .gif" onChange={handleImageChange}/>
              </div>
          </div>
              <button className="button" type="submit">Create Post</button>
        </form>
      )}
      <div id="error" className="alert"></div>
    </div>
  );
}

export default CreatePost;
