import React, { useState, useEffect } from 'react';
import { useNavigate, useParams, Link } from 'react-router-dom';
import { displayErrorMessage } from "../components/ErrorMessage";
import Header from '../components/Header';
import Footer from "../components/Footer";
import Email from './../images/email.PNG';
import DOB from './../images/dob.PNG';
import About from './../images/about.PNG';
import Avatar from './../images/avatar.PNG';
import CreateComment from "../components/CreateComment";
import Follow from "../components/Follow";

function User() {
  const navigate = useNavigate();
  const [userData, setUserData] = useState({});
  const { userId } = useParams();
  const [isFollowing, setIsFollowing] = useState(false); 
  const [followers, setFollowers] = useState([]);
  const [isPublic, setIsPublic] = useState([]);


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
        if (data.user_data.avatar) {
        const avatarPath = `/images/${data.user_data.avatar}`;
        data.user_data.avatar = avatarPath;
        } else{
          data.user_data.avatar = Avatar
        }
        setUserData(data);
        const currentUserID = data.current_user;
        const allFollowers = data.followers || [];
        setFollowers(allFollowers)
        const user = data.user_data.user_id;
        const following = followers.some(follower => follower.user_id === currentUserID) || currentUserID === user;
        setIsFollowing(following);
        const isPublicAccount = (data.public || following);
        setIsPublic(isPublicAccount)
      })
      .catch((error) => {
        displayErrorMessage(`${error.message}`);
      });
    }
  }, [navigate, userId, token, followers]);

  const sortedPosts = Array.isArray(userData.posts)
  ? userData.posts.sort((a, b) => new Date(b.date) - new Date(a.date))
  : [];

  return (
      <div className="app-container">
          <Header />
          <div className="home">
          <div>
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
                  {userData.user_data.public || isFollowing ? (
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
                        <div className="following">Following</div>
                          <div id="error" className="alert"></div>
                          {userData.following === null ? (
                            <p className="user">No users.</p>
                          ) : (
                            <div className="user">
                              {userData.following.map((user) => (
                                <div key={user.user_id}>
                                  <Link className="link" to={`/user/${user.user_id}`}>
                                    {user.first_name} {user.last_name}
                                  </Link>
                                </div>
                              ))}
                            </div>
                          )}
                        </div>
                        <div className="right-container1">
                        <div className="following">Followers</div>
                          <div id="error" className="alert"></div>
                          {userData.followers === null ? (
                            <p className="user">No users.</p>
                          ) : (
                            <div className="user">
                              {userData.followers.map((user) => (
                                <div key={user.user_id}>
                                  <Link className="link" to={`/user/${user.user_id}`}>
                                    {user.first_name} {user.last_name}
                                  </Link>
                                </div>
                              ))}
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                  ) : null}
                  <div>
                    <Follow userData={!userData.user_data.public} userId={parseInt(userId)} updateFollowers={setFollowers} isPublic={isPublic}/>
                  </div>
                  {userData.user_data.public || isFollowing ?  (
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
                              {post.image && <img className="post-image" src={`/images/${post.image}`} alt="PostImage" />}
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
                                        {comment.image && <img className="post-image" src={`/images/${comment.image}`} alt="CommentImage" />}
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
              <div id="error" className="alert"></div>
            )}
          </div>
        </div>
      <Footer />
    </div>
  );
};  

export default User;
