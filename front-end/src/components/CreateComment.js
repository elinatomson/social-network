import React, { useState } from "react";
import { displayErrorMessage } from "./ErrorMessage";

function CreateComment({ postID, addNewComment }) {
  const [commentContent, setCommentContent] = useState("");
  const [imageOrGif, setImageOrGif] = useState(null);
  const [isCommentFocused, setIsCommentFocused] = useState(false);


  const handleContentChange = (e) => {
    setCommentContent(e.target.value);
  };

  const handleImageChange = (e) => {
    const file = e.target.files[0];
    setImageOrGif(file);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    
    const commentData = new FormData();
    commentData.append("post_id", postID);
    commentData.append("comment", commentContent);
    commentData.append("image", imageOrGif);

    const headers = new Headers();

    let requestOptions = {
      body: commentData,
      method: "POST",
      headers: headers,
    }

    fetch("/create-comment", requestOptions)
      .then((response) => {
        if (response.ok) {
          setCommentContent("");
          setImageOrGif("");
          response.json().then((createdComment) => {
            const newComment = {
              comment_id: createdComment.comment_id,
              first_name: createdComment.first_name,
              last_name: createdComment.last_name,
              comment: createdComment.comment,
              image: createdComment.image,
              date: createdComment.date,
            };
            addNewComment(postID, newComment); 
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
      setIsCommentFocused(false)
  };

  return (
    <div>
      <form onSubmit={handleSubmit}>
        <input className="content" placeholder="Comment..." value={commentContent} onChange={handleContentChange} onFocus={() => setIsCommentFocused(true)} required/>
            {isCommentFocused && (
                <>
                  <label htmlFor="image"></label>
                  <input className="insert" type="file" name="image" accept="image/*, .gif" onChange={handleImageChange}/>
                  <button className="comment-button" type="submit">Add Comment</button>
                </>
            )}
      </form>
      <div id="error" className="alert"></div>
    </div>
  );
}

export default CreateComment;
