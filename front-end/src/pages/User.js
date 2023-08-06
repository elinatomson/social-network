import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { displayErrorMessage } from "../components/ErrorMessage";
import Header from '../components/Header';
import Footer from "../components/Footer";
import Email from './../images/email.PNG';
import DOB from './../images/dob.PNG';
import About from './../images/about.PNG';
import Avatar from './../images/avatar.PNG';
import CreateComment from "../components/CreateComment";
import Follow from "../components/Follow";
import Following from "../components/Following";
import Followers from "../components/Followers";


function User() {
  const navigate = useNavigate();
  const [userData, setUserData] = useState({});
  const { userId } = useParams();
  const token = document.cookie
  .split("; ")
  .find((row) => row.startsWith("sessionId="))
  ?.split("=")[1];

  useEffect(() => {
    if (!token) {
      navigate("/login");
    } else {
      fetch(`/user/${userId}`, {
        headers: {
          Authorization: `${token}`,
        },
      })
      .then((response) => {
        if (!response.ok) {
          return response.json().then((data) => {
            throw new Error(data.message);
          })
        } else {
          document.cookie = "sessionId=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/user;";
          return response.json();
        }
      })
      .then((data) => {
          setUserData(data);
      })
      .catch((error) => {
          displayErrorMessage(`An error occured while displaying user: ${error.message}`);
      });
    }
  }, [navigate, userId, token]);

  const sortedPosts = Array.isArray(userData.posts)
  ? userData.posts.sort((a, b) => new Date(b.date) - new Date(a.date))
  : [];

  return (
      <div className="app-container">
          <Header />
          <div className="home">
          <div>
            <div id="error" className="alert"></div>
            {userData.user_data ? (
              <div className="container">
                <div className="left-container">
                </div>
                <div className="middle-container">
                <div className="user-avatar">
                    <img
                      className="user-avatar-img"
                      src={userData.user_data.avatar ? userData.user_data.avatar : Avatar}
                      alt="avatar"
                    />
                  </div>
                  <div className="user-user1">
                    {userData.user_data.first_name} {userData.user_data.last_name}
                  </div>
                  {userData.user_data.public ? (
                    <div className="user-user2">
                      <p>
                        <img className="user-dob" src={DOB} alt="dob" />
                        {userData.user_data.date_of_birth}
                      </p>
                      <p>
                        <img className="user-email" src={Email} alt="email" />
                        {userData.user_data.email}
                      </p>
                      <p>
                        <img className="user-about" src={About} alt="about" />
                        Nickname: "{userData.user_data.nickname}" About the user: "
                        {userData.user_data.about_me}"
                      </p>
                      <div className="container">
                        <div className="left-container2">
                          <Following />
                        </div>
                        <div className="right-container1">
                          <Followers />
                        </div>
                      </div>
                    </div>
                  ) : null}
                  <div>
                    <Follow userData={!userData.user_data.public} userId={parseInt(userId)} />
                  </div>
                  {userData.user_data.public ? (
                    <>
                      <div className="activity">User activity</div>
                      {sortedPosts.length === 0 ? (
                        <p className="nothing">No posts found.</p>
                      ) : (
                        <div>
                          {sortedPosts.map((post) => (
                            <div className="posts" key={post.post_id}>
                              <div>
                                <span className="poster">
                                  {post.first_name} {post.last_name}
                                </span>
                                <span className="post-date">
                                  {new Date(post.date).toLocaleString()}
                                </span>
                                <div className="post-date">
                                  {post.privacy} 
                                </div>
                              </div>
                              <p className="post">{post.content}</p>
                              {post.image && (
                                <img src={post.image} alt="PostImage" />
                              )}
                              <div className="comments">
                                {post.comments === null ? (
                                  <p className="comment-text">No comments</p>
                                ) : (
                                  <div>
                                    {post.comments.map((comment) => (
                                      <div className="comment" key={comment.comment_id}>
                                        <div>
                                          <span className="poster">
                                            {comment.first_name} {comment.last_name}
                                          </span>
                                          <span className="post-date">
                                            {new Date(comment.date).toLocaleString()}
                                          </span>
                                        </div>
                                        <div className="comment-text">
                                          {comment.comment}
                                        </div>
                                      </div>
                                    ))}
                                  </div>
                                )}
                              </div>
                              <CreateComment postID={post.post_id} />
                            </div>
                          ))}
                        </div>
                      )}
                    </>
                  ) : (
                    <p className="center">This user's profile is private.</p>
                  )}
                </div>
                <div className="right-container">
                </div>
              </div>
            ) : (
              <p>Loading user data...</p>
            )}
          </div>
        </div>
      <Footer />
    </div>
  );
};  

export default User;
