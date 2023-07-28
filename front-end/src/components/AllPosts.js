import React, { useState, useEffect } from "react";

function AllPosts() {
  const [allPosts, setAllPosts] = useState([]);

  useEffect(() => {
    fetch("/all-posts")
      .then((response) => response.json())
      .then((data) => {
        setAllPosts(data);
      })
      .catch((error) => {
        console.error("Error fetching all posts:", error);
      });
  }, []);

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
