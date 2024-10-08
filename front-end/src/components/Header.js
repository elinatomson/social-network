import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import Avatar from '../images/avatar.PNG';
import FollowRequests from "../components/FollowRequests";
import GroupRequests from "../components/GroupRequests";
import GroupInvitations from "../components/GroupInvitations";
import { displayErrorMessage } from "../components/ErrorMessage";
import GroupEventNotifications from './GroupEventNotifications';

function Header() {
  const [userData, setUserData] = useState(null);

  useEffect(() => {
      fetch("/main")
      .then((response) => response.json())
      .then((data) => {
        if (data.avatar) {
          const avatarPath = `/images/${data.avatar}`;
          data.avatar = avatarPath;
        } else {
          data.avatar = Avatar;
        }
        setUserData(data);
      })
      .catch((error) => {
        displayErrorMessage(`${error.message}`);
      });
  }, []);

    const isOnMainPage = window.location.pathname === '/main';
  
    const handleClick = () => {
      if (isOnMainPage) {
        window.location.reload();
      }
    }


  return (
    <header>
        {userData ? (
          <div className="container">
            <div className="left-container">
              <div className="request">
                  <FollowRequests />
                  <GroupRequests />
                  <GroupInvitations />
                  <GroupEventNotifications />
              </div>
              <div className="avatar">
                  <img
                  className="avatar-img"
                  src={userData.avatar}
                  alt="avatar"
                  />
              </div>
              <div className="user1">{userData.first_name} {userData.last_name}</div>
              <Link className="profile-btn" to="/profile">Profile</Link>
            </div>
            <div className="middle-container">
              <span className="heading1">Welcome to </span><br/>
              <span className="heading2">Social Network page</span>
            </div>
            <div className="right-container">
              <Link className="log-out-button" to="/main" onClick={handleClick}>Main Page</Link>
              <Link className="log-out-button" to="/logout">Log Out</Link>
            </div>
          </div>
          ) : (
          <div id="error" className="alert"></div>
        )}
    </header>
  );
}

export default Header