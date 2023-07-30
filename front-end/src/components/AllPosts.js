import React, { useState, useEffect } from "react";
import { useNavigate } from 'react-router-dom';
import { displayErrorMessage } from "./ErrorMessage";

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

export default AllPosts;
