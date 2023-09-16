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
    fetch("/all-posts")
      .then((response) => response.json())
      .then((data) => {
        setAllPosts(data);
      })
      .catch((error) => {
        displayErrorMessage(`${error.message}`);
      });
  }, [navigate, postContent]);

  const sortedPosts = Array.isArray(allPosts)
  ? allPosts.sort((a, b) => new Date(b.date) - new Date(a.date))
  : [];

  const addNewComment = (postID, newComment) => {
    //finding the post by postID and update its comments array
    setAllPosts((prevAllPosts) =>
      prevAllPosts.map((post) =>
        post.post_id === postID
          ? {
              ...post,
              //If post.comments is null, it means there are no existing comments on this post. So initializing the comments property as a new array containing only the newComment.
              //If post.comments is not null, it means there are already some comments on this post. In this case, using the spread operator (...) again to create a new array. Coping all the existing comments from post.comments and adding the newComment to the end of the array.
              comments: post.comments === null ? [newComment] : [...post.comments, newComment],
            }
          //If the post_id doesn't match, simply returning the original post object without any changes.
          : post
      )
    );
  };
  
  return (
    <div>
      <div id="error" className="alert"></div>
      {sortedPosts.length === 0 ? (
        <p className="nothing">No posts.</p>
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
              {post.image && <img className="post-image" src={`/images/${post.image}`} alt="PostImage" />}
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
                        {comment.image && <img className="post-image" src={`/images/${comment.image}`} alt="CommentImage" />}
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
