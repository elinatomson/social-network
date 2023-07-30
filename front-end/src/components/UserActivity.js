import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

function UserActivity() {
  const [userPosts, setUserPosts] = useState([]);

  const navigate = useNavigate();

  useEffect(() => {
    const token = document.cookie
      .split("; ")
      .find((row) => row.startsWith("sessionId="))
      ?.split("=")[1];

    if (!token) {
      navigate("/login");
    } else {
    fetch("/activity")
      .then((response) => response.json())
      .then((data) => {
        setUserPosts(data);
      })
      .catch((error) => {
        console.error("Failed to fetch user's posts:", error);
      });
    }
  }, [navigate]);

  const sortedPosts = Array.isArray(userPosts)
  ? userPosts.sort((a, b) => new Date(b.date) - new Date(a.date))
  : [];
  

  return (
    <div>
        {sortedPosts.length === 0 ? (
            <p className="nothing">No posts found.</p>
        ) : (
          <ul>
            {userPosts.map((post) => (
            <div className="post" key={post.post_id}>
            <p className="poster">
                {post.first_name} {post.last_name}
            </p>
            <p className="post-date">
                {new Date(post.date).toLocaleString()}
            </p>
            <p>{post.content}</p>
            {post.image && <img src={post.image} alt="PostImage" />}
        </div>
            ))}
            </ul>
            )}
    </div>
  );
}

export default UserActivity;
