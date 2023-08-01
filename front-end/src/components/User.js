import React, { useState, useEffect } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { displayErrorMessage } from "./ErrorMessage";
import Email from './../images/email.PNG';
import DOB from './../images/dob.PNG';
import About from './../images/about.PNG';
import Avatar from './../images/avatar.PNG';
import CreateComment from "./CreateComment";


function User() {
  const navigate = useNavigate();
  const [userData, setUserData] = useState({});

  const { userId } = useParams();

  useEffect(() => {
    const token = document.cookie
      .split("; ")
      .find((row) => row.startsWith("sessionId="))
      ?.split("=")[1];

    if (!token) {
      navigate("/login");
    } else {
      fetch(`/user/${userId}`, {
        headers: {
          Authorization: `${token}`,
        },
      })
        .then((response) => response.json())
        .then((data) => {
          setUserData(data);
        })
        .catch((error) => {
          displayErrorMessage(`${error.message}`);
        });
    }
  }, [navigate, userId]);

  const sortedPosts = Array.isArray(userData.posts)
  ? userData.posts.sort((a, b) => new Date(b.date) - new Date(a.date))
  : [];

  return (
    <div>
      {userData.user_data ? (
        <div className="container">
          <div className="left-container">
            <div className="avatar">
              <img
                className="avatar-img"
                src={userData.user_data.avatar ? userData.user_data.avatar : Avatar}
                alt="avatar"
              />
            </div>
            <div className="user1">
              {userData.user_data.first_name} {userData.user_data.last_name}
            </div>
            {userData.user_data.public ? (
              <div className="user2">
                <p>
                  <img className="dob" src={DOB} alt="dob"></img>{" "}
                  {userData.user_data.date_of_birth}
                </p>
                <p>
                  <img className="email" src={Email} alt="email"></img>{" "}
                  {userData.user_data.email}
                </p>
                <p>
                  <img className="about" src={About} alt="about"></img>
                  Nickname: "{userData.user_data.nickname}" About me: "
                  {userData.user_data.about_me}"
                </p>
              </div>
            ) : null}
            <div className="container">
              <div className="left-container1">
                <Link className="link-btn" to={`/following`}>
                  {" "}
                  <div>(nr)</div>
                  <div>Following</div>
                </Link>
              </div>
              <div className="right-container1">
                <Link className="link-btn" to={`/followers`}>
                  {" "}
                  <div>(nr)</div>
                  <div>Followers</div>
                </Link>
              </div>
            </div>
          </div>
          {userData.user_data.public ? (
            <div className="middle-container">
              <Link className="profile-type-button">Follow/unfollow</Link>
              <div className="activity">User activity</div>
              {sortedPosts.length === 0 ? (
                <p className="nothing">No posts found.</p>
              ) : (
                <ul>
                  {sortedPosts.map((post) => (
                    <div className="posts" key={post.post_id}>
                      <div>
                        <span className="poster">
                          {post.first_name} {post.last_name}
                        </span>
                        <span className="post-date">
                          {new Date(post.date).toLocaleString()}
                        </span>
                      </div>
                      <p className="post">{post.content}</p>
                      {post.image && (
                        <img src={post.image} alt="PostImage" />
                      )}
                      <div className="comments">
                        {post.comments === null ? (
                          <p>No comments</p>
                        ) : (
                          <ul>
                            {post.comments.map((comment) => (
                              <div class="comment" key={comment.comment_id}>
                                <div>
                                  <span className="poster">
                                    {comment.first_name} {comment.last_name}
                                  </span>
                                  <span className="post-date">
                                    {new Date(comment.date).toLocaleString()}
                                  </span>
                                </div>
                                {comment.comment}
                              </div>
                            ))}
                          </ul>
                        )}
                      </div>
                      <CreateComment postID={post.post_id} />
                    </div>
                  ))}
                </ul>
              )}
            </div>
          ) : (
            <div className="middle-container">
              <Link className="profile-type-button">Follow/unfollow</Link>
              <p>This user's profile is private.</p>
            </div>
          )}
          <div className="right-container">
            <Link className="log-out-button" to="/main">
              Main Page
            </Link>
            <Link className="log-out-button" to="/logout">
              Log Out
            </Link>
          </div>
        </div>
      ) : (
        <p>Loading user data...</p>
      )}
    </div>
  );
}

export default User;
