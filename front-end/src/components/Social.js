import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import Avatar from './../images/avatar.PNG';

function Social() {
  const navigate = useNavigate();
  const [userData, setUserData] = useState(null);
  const [searchTerm, setSearchTerm] = useState("");
  const [searchResults, setSearchResults] = useState(null);

  useEffect(() => {
    const token = document.cookie
      .split("; ")
      .find((row) => row.startsWith("sessionId="))
      ?.split("=")[1];

    if (!token) {
      navigate("/login");
    } else {
      // Fetch user data from the backend API.
      fetch("/social", {
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
  
  const handleSearch = (e) => {
    const searchTerm = e.target.value;
    setSearchTerm(searchTerm);
    if (searchTerm === '') {
      setSearchResults(null); 
      return;
    }

    fetch(`/search?query=${searchTerm}`)
      .then((response) => response.json())
      .then((data) => {
        setSearchResults(data.length > 0 ? data : []);
      })
      .catch((error) => {
        console.error("Failed to perform search:", error);
        setSearchResults([]);
      });
  };

    return (
        <div>
          {userData ? (
            <>
              <Link className="log-out-button" to="/logout">Log Out</Link>
              <div className="search">
                <input type="text" placeholder="Search users..." value={searchTerm} onChange={handleSearch} />
                <div className="search-results">
                  {searchResults !== null && searchResults.length > 0 && (
                    searchResults.map((result) => (
                      <div key={result.user_id} className="search-result-item">
                        <img className="avatar-img" src={result.avatar ? result.avatar : Avatar} alt="avatar" />
                        <div>{result.first_name} {result.last_name}</div>
                        <Link to={`/user/${result.user_id}`}>View Profile</Link>
                      </div>
                    ))
                  )}
                {searchResults !== null && searchResults.length === 0 && (
                  <p>No results found for your search query.</p>
                )}
            </div>
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
            </>
          ) : (
            <p>Loading data...</p>
          )}
        </div>
      );
}

export default Social