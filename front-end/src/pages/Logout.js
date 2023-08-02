import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { displayErrorMessage } from "../components/ErrorMessage";

function Logout() {
  const navigate = useNavigate();

  useEffect(() => {
    fetch("/logout")
      .catch(error => {
        displayErrorMessage(`${error.message}`);
      })
      .finally(() => {
        // Remove the "session" cookie by setting its expiration to the past
        document.cookie = "sessionId=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
        navigate("/");
      });
  }, [navigate]);

  return null
}

export default Logout;
