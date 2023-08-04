import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import Avatar from '../images/avatar.PNG';
import FollowRequests from "../components/FollowRequests";
import Search from "../components/Search";

function Header() {

  const navigate = useNavigate();
  const [userData, setUserData] = useState(null);
  const [searchResults, setSearchResults] = useState(null);

  useEffect(() => {
    const token = document.cookie
      .split("; ")
      .find((row) => row.startsWith("sessionId="))
      ?.split("=")[1];

    if (!token) {
      navigate("/");
    } else {
      fetch("/main", {
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
    <header>
          {userData ? (
            <div className="container">
                <div className="left-container">
                    <div className="request">
                        <FollowRequests />
                    </div>
                    <div className="avatar">
                        <img
                        className="avatar-img"
                        src={userData.avatar ? userData.avatar : Avatar}
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
                    <Link className="log-out-button" to="/main">Main Page</Link>
                    <Link className="log-out-button" to="/logout">Log Out</Link>
                    <Search setSearchResults={setSearchResults}/>
                <div className="search-results">
                  {searchResults !== null && searchResults.length > 0 && (
                    searchResults.map((result) => (
                      <div key={result.user_id} className="search-result-item">
                        <Link className="link-btn" to={`/user/${result.user_id}`} >{result.first_name} {result.last_name}</Link>
                      </div>
                    ))
                  )}
                  {searchResults !== null && searchResults.length === 0 && (
                  <p>No results found for your search query.</p>
                )}
                </div>
            </div>
            </div>
            ) : (
                <p></p>
            )}
    </header>
  );
}

export default Header