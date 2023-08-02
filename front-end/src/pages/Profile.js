import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { displayErrorMessage } from "../components/ErrorMessage";
import Header from '../components/Header';
import Footer from "../components/Footer";
import Email from './../images/email.PNG';
import DOB from './../images/dob.PNG';
import About from './../images/about.PNG';
import Avatar from './../images/avatar.PNG';
import CreateComment from "../components/CreateComment";
import Following from "../components/Following";
import Followers from "../components/Followers";
import FollowRequests from "../components/FollowRequests";

function Profile() {
  const navigate = useNavigate();
  const [userData, setUserData] = useState({});
  const [profileType, setProfileType] = useState('');

  const token = document.cookie
  .split("; ")
  .find((row) => row.startsWith("sessionId="))
  ?.split("=")[1];

  useEffect(() => {
    if (!token) {
      navigate("/login");
    } else {
      fetch("/profile", {
        headers: {
          Authorization: `${token}`,
        },
      })
        .then((response) => response.json())
        .then((data) => {
          setUserData(data);
          // Read the current profile type from localStorage if available for the current user
          const storedProfileType = localStorage.getItem(`profileType_${data.user_data.email}`); // Use user email as the key
          // Set the initial profile type based on localStorage or user data
          const initialProfileType = storedProfileType || (data.user_data.public ? 'Set your profile as public' : 'Set your profile as private');
          setProfileType(initialProfileType);
          // If the profile type was not in localStorage, then store the initial value in localStorage for the current user
          if (!storedProfileType) {
            localStorage.setItem(`profileType_${data.user_data.email}`, initialProfileType); // Use user email as the key
          }
        })
        .catch((error) => {
          displayErrorMessage(`${error.message}`);
        });
    }
  }, [navigate, token]);

  const sortedPosts = Array.isArray(userData.posts)
  ? userData.posts.sort((a, b) => new Date(b.date) - new Date(a.date))
  : [];

  const handleProfileTypeToggle = () => {
    if (!token) {
      navigate("/login");
    } else {
      fetch("/profile-type", {
        method: "POST",
        headers: {
          Authorization: `${token}`,
        },
      })
      .then((response) => {
        if (response.ok) {
          // Update the profileType only when the button is clicked
          const newProfileType = profileType === "Set your profile as private" ? "Set your profile as public" : "Set your profile as private";
          setProfileType(newProfileType);
          // Update the profileType in localStorage for the current user
          localStorage.setItem(`profileType_${userData.user_data.email}`, newProfileType); // Use user email as the key
        } else {
          throw new Error("Failed to update profile type");
        }
      })
      .catch((error) => {
        displayErrorMessage(`${error.message}`);
      });
    };
  }

  return (
    <div className="app-container">
      <div className="content">
        <Header />
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
                <div className="user1">{userData.user_data.first_name} {userData.user_data.last_name}</div>
                <div className="user2">
                  <p><img className="dob" src={DOB} alt="dob"></img> {userData.user_data.date_of_birth}</p>
                  <p><img className="email" src={Email} alt="email"></img> {userData.user_data.email}</p>
                  <p>
                    <img className="about" src={About} alt="about"></img>
                    Nickname: "{userData.user_data.nickname}" About me: "{userData.user_data.about_me}"
                  </p>
                  <button className="profile-type-button" onClick={handleProfileTypeToggle}>{profileType}</button>
                  <FollowRequests />
                </div>
                <div className="container">
                  <div className="left-container1">
                    <Following />
                  </div>
                  <div className="right-container1">
                    <Followers />
                  </div>
                </div>
              </div>
              <div className="middle-container">
                <div className="activity">My activity</div>
                {sortedPosts.length === 0 ? (
                <p className="nothing">No activity.</p>
            ) : (
              <div>
                {sortedPosts.map((post) => (
                <div className="posts" key={post.post_id}>
                  <div>
                    <span className="poster">{post.first_name} {post.last_name}</span>
                    <span className="post-date">{new Date(post.date).toLocaleString()}</span>
                  </div>
                  <p className="post">{post.content}</p>
                  {post.image && <img src={post.image} alt="PostImage" />}
                  <div className="comments">
                    {post.comments === null ? (
                      <p className="comment-text">No comments</p>
                    ) : (
                      <div>
                        {post.comments.map((comment) => (
                          <div className="comment" key={comment.comment_id}>
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
                  <CreateComment postID={post.post_id}/>
                </div>
                ))}
                </div>
                  )}
              </div>
              <div className="right-container">
                <Link className="log-out-button" to="/main">Main Page</Link>
                <Link className="log-out-button" to="/logout">Log Out</Link>
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
}

export default Profile;
