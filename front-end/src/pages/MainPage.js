
import Footer from "../components/Footer"
import Header from "../components/Header"
import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import Avatar from '../images/avatar.PNG';
import Search from "../components/Search";
import CreatePost from "../components/CreatePost";
import AllPosts from "../components/AllPosts";
import FollowRequests from "../components/FollowRequests";

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
    <div className="app-container">
      <Header />
      <div className="home">
        <div>
          {userData ? (
            <div className="container">
              <div className="left-container">
              </div>
              <div className="middle-container">
                <CreatePost/>
                <AllPosts />
              </div>
              <div className="right-container">
            </div>
          </div>
          ) : (
            <p>Loading data...</p>
          )}
        </div>
      </div>
      <Footer />
    </div>
  );
}

export default MainPage;
