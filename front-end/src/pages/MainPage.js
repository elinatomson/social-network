
import Footer from "../components/Footer"
import Header from "../components/Header"
import React, { useState, useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { displayErrorMessage } from "../components/ErrorMessage";
import CreatePost from "../components/CreatePost";
import AllPosts from "../components/AllPosts";
import Users from "../components/Users";
import Groups from "../components/AllGroups";
import CreateGroup from "../components/CreateGroup";
import Search from "../components/Search";

function MainPage() {
  const navigate = useNavigate();
  const [userData, setUserData] = useState(null);
  const [searchResults, setSearchResults] = useState(null);

  useEffect(() => {
    const token = document.cookie
      .split("; ")
      .find((row) => row.startsWith("sessionId="))
      ?.split("=")[1];

    if (!token) {
      navigate("/login");
    } else {
      fetch("/main", {
        headers: {
          Authorization: `${token}`,
        },
      })
      .then((response) => {        
        if (!response.ok) {
        return response.json().then((data) => {
          throw new Error(data.message);
        });
        } else {
          return response.json();
        }
      })
      .then((data) => {
        setUserData(data);
      })
      .catch((error) => {
        displayErrorMessage(`${error.message}`);
      });
    }
  }, [navigate]);

  return (
    <div className="app-container">
      <Header />
      <div className="home">
        <div>
          {userData ? (
          <div className="container">
            <div className="left-container">
              <Users />
            </div>
            <div className="middle-container">
              <CreateGroup/>
              <CreatePost/>
              <AllPosts />
            </div>
            <div className="right-container2">
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
              <Groups />
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
}

export default MainPage;
