import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { displayErrorMessage } from "../components/ErrorMessage";
import Header from '../components/Header';
import Footer from "../components/Footer";
import RequestToJoinGroup from "../components/RequestToJoinGroup";

function Group() {
  const navigate = useNavigate();
  const [groupData, setGroupData] = useState({});
  const { groupId } = useParams();
  const [isMember, setIsMember] = useState(false);
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
          const currentUserID = data.userID;
          const selectedUserIDs = data.group.selected_user_id.split(",").map(id => parseInt(id));
          const groupCreator = data.group.user_id
          setIsMember(selectedUserIDs.includes(currentUserID) || currentUserID === groupCreator);
      })
      .catch((error) => {
          displayErrorMessage(`An error occured while displaying group: ${error.message}`);
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
                  {isMember && (
                    <div className="users">
                        <div className="following">Group creator</div>
                        {groupData.group.first_name} {groupData.group.last_name}
                        <div className="following">Group members</div>
                        {groupData.group.selected_user_id}
                    </div>
                  )}
                </div>
                <div className="middle-container">
                    <div className="activity">
                      {groupData.group && groupData.group.title}
                    </div>
                    {isMember && (
                    <div className="nothing">
                      {groupData.group && groupData.group.description}
                    </div>
                    )}
                    {!isMember && (
                      <RequestToJoinGroup groupId={parseInt(groupId)} />
                    )}
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
