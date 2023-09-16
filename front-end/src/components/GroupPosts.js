import React, { useState, useEffect } from "react";
import { useNavigate, useLocation } from 'react-router-dom';
import { displayErrorMessage } from "./ErrorMessage";
import CreateComment from "./CreateComment";

function GroupPosts({ groupId }) {
  const [allPosts, setAllPosts] = useState([]);
  const navigate = useNavigate();
  const location = useLocation();
  const { postContent } = location.state || {};

  useEffect(() => {
    fetch(`/group-posts?groupId=${groupId}`)
    .then((response) => response.json())
    .then((data) => {
      setAllPosts(data);
    })
    .catch((error) => {
      displayErrorMessage(`${error.message}`);
    });
  }, [groupId, navigate, postContent]);

  const sortedPosts = Array.isArray(allPosts)
  ? allPosts.sort((a, b) => new Date(b.date) - new Date(a.date))
  : [];

  const addNewComment = (postID, newComment) => {
    setAllPosts((prevAllPosts) =>
      prevAllPosts.map((post) =>
        post.post_id === postID
          ? {
              ...post,
              comments: post.comments === null ? [newComment] : [...post.comments, newComment],
            }
          : post
      )
    );
  };
  
  return (
    <div>
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

export default GroupPosts;
