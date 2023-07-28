import React, { useState } from "react";
import { displayErrorMessage } from "./ErrorMessage";
import { useNavigate } from "react-router-dom"

function CreatePost() {
  const [postContent, setPostContent] = useState("");
  const [postPrivacy, setPostPrivacy] = useState("public");
  const [imageOrGif, setImageOrGif] = useState(null);

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

  const handleSubmit = (e) => {
    e.preventDefault();
    // Handle the form submission here (e.g., API call to create the post)
    // You can access the postContent, postPrivacy, and imageOrGif states to send the data to the server
    // Reset the form after successful submission
    const postData = {
    content: postContent,
    privacy: postPrivacy,
    image: imageOrGif,
    }

    const headers = new Headers()
    headers.append("Content-Type", "application/json")

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
          setImageOrGif(null);
          navigate("/social")
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
        displayErrorMessage(`An error occurred while posting: ${error.message}`);
      });
  };

  return (
    <div className="posting">
      <form onSubmit={handleSubmit}>
        <input className="content" placeholder="Post something..." value={postContent} onChange={handleContentChange} required></input>
        <div className="container1">
            <div className="left-container1">
                <label htmlFor="privacy"></label>
                <select className="privacy" name="privacy" value={postPrivacy} onChange={handlePrivacyChange} required>
                    <option value="public">Public</option>
                    <option value="private">Private</option>
                    <option value="almost-private">Almost Private</option>
                </select>
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
    </div>
  );
}

export default CreatePost;
