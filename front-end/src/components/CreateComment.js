import React, { useState } from "react";
import { displayErrorMessage } from "./ErrorMessage";

function CreateComment({ postID, addNewComment }) {
  const [commentContent, setCommentContent] = useState("");
  const [imageOrGif, setImageOrGif] = useState("");
  const [isCommentFocused, setIsCommentFocused] = useState(false);


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

    let requestOptions = {
      body: JSON.stringify(commentData),
      method: "POST",
      headers: headers,
    }

    fetch("/create-comment", requestOptions)
      .then((response) => {
        if (response.ok) {
          setCommentContent("");
          setImageOrGif("");
          response.json().then((createdComment) => {
            // Construct the new comment object using the response from the server
            const newComment = {
              comment_id: createdComment.comment_id,
              first_name: createdComment.first_name,
              last_name: createdComment.last_name,
              comment: createdComment.comment,
              date: createdComment.date,
            };
            addNewComment(postID, newComment); // Call the function to update the state of allPosts
          });
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
          onBlur={handleBlur} required/>
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
