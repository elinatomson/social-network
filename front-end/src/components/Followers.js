import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { displayErrorMessage } from "../components/ErrorMessage";
import { useNavigate, useLocation } from 'react-router-dom';

function Followers() {
  const [followersUsers, setFollowersUsers] = useState([]);
  const navigate = useNavigate();
  const location = useLocation();
  const { followContent } = location.state || {};

  useEffect(() => {
    fetch('/followers')
      .then((response) => response.json())
      .then((data) => {
        if (data && data.followers_users) {
            setFollowersUsers(data.followers_users);
          }
      })
      .catch((error) => {
        displayErrorMessage(`${error.message}`);
      });
    }, [navigate, followContent]);

  return (
    <div>
      <div className="following">Followers</div>
      <div id="error" className="alert"></div>
      {followersUsers.length === 0 ? (
        <p className="user">No users.</p>
      ) : (
        <div className="user">
          {followersUsers.map((user) => (
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
export default Followers;
