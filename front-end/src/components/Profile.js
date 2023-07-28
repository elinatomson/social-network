import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import Email from './../images/email.PNG';
import DOB from './../images/dob.PNG';
import About from './../images/about.PNG';
import Avatar from './../images/avatar.PNG';

function Profile() {
  const navigate = useNavigate();
  const [userData, setUserData] = useState(null);

  useEffect(() => {
    const token = document.cookie
      .split("; ")
      .find((row) => row.startsWith("sessionId="))
      ?.split("=")[1];

    if (!token) {
      navigate("/login");
    } else {
      // Fetch user data from the backend API.
      fetch("/profile", {
        headers: {
          Authorization: `${token}`,
        },
      })
        .then((response) => response.json())
        .then((data) => {
          setUserData(data);
        })
        .catch((error) => {
          console.error("Failed to fetch user data:", error);
        });
    }
  }, [navigate]);

  return (
    <div>
      {userData ? (
        <div className="container">
          <div className="left-container">
            <div className="avatar">
              <img
                className="avatar-img"
                src={userData.avatar ? userData.avatar : Avatar}
                alt="avatar"
              />
            </div>
            <div className="user1">{userData.first_name} {userData.last_name}</div>
            <div className="user2">
              <p><img className="dob" src={DOB} alt="dob"></img> {userData.date_of_birth}</p>
              <p><img className="email" src={Email} alt="email"></img> {userData.email}</p>
              <p>
                <img className="about" src={About} alt="about"></img>
                Nickname: "{userData.nickname}" About me: "{userData.about_me}"
              </p>
            </div>
          </div>
          <div className="middle-container">
            <div>
              (In the middle comes posts)
            </div>
          </div>
          <div className="right-container">
            <Link className="log-out-button" to="/social">Main Page</Link>
            <Link className="log-out-button" to="/logout">Log Out</Link>
          </div>
        </div>
      ) : (
        <p>Loading user data...</p>
      )}
    </div>
  );
}

export default Profile;
