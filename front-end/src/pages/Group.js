import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { displayErrorMessage } from "../components/ErrorMessage";
import Header from '../components/Header';
import Footer from "../components/Footer";

function Group() {
  const navigate = useNavigate();
  const [groupData, setGroupData] = useState({});
  const { groupId } = useParams();
  const token = document.cookie
  .split("; ")
  .find((row) => row.startsWith("sessionId="))
  ?.split("=")[1];

  useEffect(() => {
    if (!token) {
      navigate("/login");
    } else {
      fetch(`/group/${groupId}`, {
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
          return response.json();
        }
      })
      .then((data) => {
          setGroupData(data);
      })
      .catch((error) => {
          displayErrorMessage(`An error occured while displaying user: ${error.message}`);
      });
    }
  }, [navigate, groupId, token]);

  return (
      <div className="app-container">
          <Header />
          <div className="home">
          <div>
            <div id="error" className="alert"></div>
              <div className="container">
                <div className="left-container">
                    <div className="users">
                        <div className="following">Group members</div>
                        {groupData.selected_user_id}
                    </div>
                </div>
                <div className="middle-container">
                    <div className="activity">
                        {groupData.title}
                    </div>
                    <div className="nothing">
                        {groupData.description}
                    </div>
                </div>
                <div className="right-container">
                </div>
              </div>
          </div>
        </div>
      <Footer />
    </div>
  );
};  

export default Group;
