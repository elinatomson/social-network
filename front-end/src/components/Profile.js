import React, { useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';

function Profile() {
  const navigate = useNavigate();

  useEffect(() => {
    const token = document.cookie
      .split("; ")
      .find((row) => row.startsWith("sessionId="))
      ?.split("=")[1];

    if (!token) {
      navigate("/login");
    }
  }, [navigate]);

  return (
    <div>
      <h2>Profile Page</h2>
      {/* Add profile content here */}
      <Link className="button" to="/logout">
        Log Out
      </Link>
    </div>
  );
}

export default Profile;
