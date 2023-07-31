import React, { useState, useEffect } from "react";
import { displayErrorMessage } from "./ErrorMessage";
import { useNavigate } from "react-router-dom"

function CreatePost() {
  const [postContent, setPostContent] = useState("");
  const [postPrivacy, setPostPrivacy] = useState("public");
  const [imageOrGif, setImageOrGif] = useState(null);

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

  const handleSubmit = (e) => {
    e.preventDefault();
    
    const postData = {
    content: postContent,
    privacy: postPrivacy,
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
          setImageOrGif(null);
          navigate("/main")
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
      <div id="error"></div>
    </div>
  );
}

export default CreatePost;
