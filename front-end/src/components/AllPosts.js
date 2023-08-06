import React, { useState, useEffect } from "react";
import { useNavigate, useLocation } from 'react-router-dom';
import { displayErrorMessage } from "./ErrorMessage";
import CreateComment from "./CreateComment";

function AllPosts() {
  const [allPosts, setAllPosts] = useState([]);
  const navigate = useNavigate();
  const location = useLocation();
  //to automatically display a new post from CreatePost.js right after submitting
  const { postContent } = location.state || {};

  useEffect(() => {
    const token = document.cookie
    .split("; ")
    .find((row) => row.startsWith("sessionId="))
    ?.split("=")[1];

  if (!token) {
    navigate("/login");
  } else {
    fetch("/all-posts", {
      headers: {
        Authorization: `${token}`,
      },
    })
      .then((response) => response.json())
      .then((data) => {
        setAllPosts(data);
      })
      .catch((error) => {
        displayErrorMessage(`${error.message}`);
      });
    }
  }, [navigate, postContent]);

  const sortedPosts = Array.isArray(allPosts)
  ? allPosts.sort((a, b) => new Date(b.date) - new Date(a.date))
  : [];

  const addNewComment = (postID, newComment) => {
    // Find the post by postID and update its comments array
    setAllPosts((prevAllPosts) =>
      prevAllPosts.map((post) =>
      //"post.post_id === postID ? ... : post" checks if the post_id of the current post matches the postID provided as an argument to the addNewComment function. 
      //If they match, it means we have found the post to which the new comment needs to be added, and we should modify this post.
        post.post_id === postID
          ? {
              ...post,
              //If post.comments is null, it means there are no existing comments on this post. So, we initialize the comments property as a new array containing only the newComment.
              //If post.comments is not null, it means there are already some comments on this post. In this case, we use the spread operator (...) again to create a new array. We copy all the existing comments from post.comments and add the newComment to the end of the array.
              comments: post.comments === null ? [newComment] : [...post.comments, newComment],
            }
          //If the post_id doesn't match (i.e., we are not updating this specific post), we simply return the original post object without any changes.
          : post
      )
    );
  };
  
  return (
    <div>
      {sortedPosts.length === 0 ? (
        <p>No posts found.</p>
      ) : (
        <div>
          {sortedPosts.map((post) => (
            <div className="posts" key={post.post_id}>
              <div>
                <span className="poster">{post.first_name} {post.last_name}</span>
                <span className="post-date">{new Date(post.date).toLocaleString()}</span>
              </div>
              <div className="post-date">
                {post.privacy} 
              </div>
              <p className="post">{post.content}</p>
              {post.image && <img src={post.image} alt="PostImage" />}
              <div className="comments">
                {post.comments === null ? (
                  <p className="comment-text">No comments</p>
                ) : (
                  <div>
                    {post.comments.map((comment, index) => (
                      <div className="comment" key={`${comment.comment_id}-${index}`}>
                        <div>
                          <span className="poster">{comment.first_name} {comment.last_name}</span>
                          <span className="post-date">{new Date(comment.date).toLocaleString()}</span>
                        </div>
                        <div className="comment-text">
                          {comment.comment}
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
              <CreateComment postID={post.post_id} addNewComment={addNewComment}/>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default AllPosts;
