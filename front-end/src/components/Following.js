import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { displayErrorMessage } from "../components/ErrorMessage";
import { useNavigate, useLocation } from 'react-router-dom';


function Following() {
  const [followingUsers, setFollowingUsers] = useState([]);
  const navigate = useNavigate();
  const location = useLocation();
  const { followContent } = location.state || {};

  useEffect(() => {
    fetch('/following')
      .then((response) => response.json())
      .then((data) => {
        if (data) {
            setFollowingUsers(data);
          }
      })
      .catch((error) => {
        displayErrorMessage(`${error.message}`);
      });
    }, [navigate, followContent]);


  return (
    <div>
      <div className="following">Following</div>
      <div id="error" className="alert"></div>
      {followingUsers.length === 0 ? (
        <p className="user">No users.</p>
      ) : (
        <div className="user">
          {followingUsers.map((user) => (
            <div key={user.user_id}>
              <Link className="link" to={`/user/${user.user_id}`}>
                {user.first_name} {user.last_name}
              </Link>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
export default Following;
