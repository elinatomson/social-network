import React, { useState, useEffect } from "react";
import { useNavigate } from 'react-router-dom';
import { displayErrorMessage } from "./ErrorMessage";
import CreateComment from "./CreateComment";

function AllPosts() {
  const [allPosts, setAllPosts] = useState([]);
  const navigate = useNavigate();

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
  }, [navigate]);

  const sortedPosts = Array.isArray(allPosts)
  ? allPosts.sort((a, b) => new Date(b.date) - new Date(a.date))
  : [];
  
  return (
    <div>
      {sortedPosts.length === 0 ? (
        <p>No posts found.</p>
      ) : (
        <ul>
          {sortedPosts.map((post) => (
            <div className="posts" key={post.post_id}>
              <div>
                <span className="poster">{post.first_name} {post.last_name}</span>
                <span className="post-date">{new Date(post.date).toLocaleString()}</span>
              </div>
              <p className="post">{post.content}</p>
              {post.image && <img src={post.image} alt="PostImage" />}
              <div className="comments">
                {post.comments.length === 0 ? (
                  <p>No comments</p>
                ) : (
                  <ul>
                    {post.comments.map((comment) => (
                      <div class="comment" key={comment.comment_id}>
                        <div>
                          <span className="poster">{comment.first_name} {comment.last_name}</span>
                          <span className="post-date">{new Date(comment.date).toLocaleString()}</span>
                        </div>
                        {comment.comment}
                      </div>
                    ))}
                  </ul>
                )}
              </div>
              <CreateComment postID={post.post_id}/>
            </div>
          ))}
        </ul>
      )}
    </div>
  );
}

export default AllPosts;
