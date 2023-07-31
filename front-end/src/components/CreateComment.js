import React, { useState, useEffect } from "react";
import { displayErrorMessage } from "./ErrorMessage";
import { useNavigate } from "react-router-dom"

function CreateComment({ postID }) {
  const [commentContent, setCommentContent] = useState("");
  const [imageOrGif, setImageOrGif] = useState(null);
  const [isCommentFocused, setIsCommentFocused] = useState(false);

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
    setCommentContent(e.target.value);
  };

  const handleImageChange = (e) => {
    const file = e.target.files[0];
    setImageOrGif(file);
  };

  const handleBlur = () => {
    // Adding a small delay before hiding the input fields to allow clicking the submit button
    setTimeout(() => {
      setIsCommentFocused(false);
    }, 100);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    
    const commentData = {
        post_id: postID,
        comment: commentContent,
        image: imageOrGif,
    };

    const headers = new Headers();
    headers.append("Content-Type", "application/json");
    headers.append("Authorization", token);

    let requestOptions = {
      body: JSON.stringify(commentData),
      method: "POST",
      headers: headers,
    }

    fetch("/create-comment", requestOptions)
      .then((response) => {
        if (response.ok) {
          setCommentContent("");
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
    <div>
      <form onSubmit={handleSubmit}>
        <input className="content" placeholder="Comment..." value={commentContent} onChange={handleContentChange} onFocus={() => setIsCommentFocused(true)}
          onBlur={handleBlur} />
            {isCommentFocused && (
                <>
                    <label htmlFor="image"></label>
                    <input className="insert" type="file" name="image" accept="image/*, .gif" value={imageOrGif} onChange={handleImageChange}/>
                    <button className="comment-button" type="submit">Add Comment</button>
                </>
            )}
      </form>
    </div>
  );
}

export default CreateComment;
