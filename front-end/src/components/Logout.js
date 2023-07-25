import { useEffect } from 'react';
import { useNavigate, useOutletContext } from 'react-router-dom';

function Logout() {
  const navigate = useNavigate();
  const { setToken } = useOutletContext();

  useEffect(() => {
    const requestOptions = {
      method: "GET",
      credentials: "include",
    };

    fetch("/logout", requestOptions)
      .catch(error => {
        console.log("error logging out", error);
      })
      .finally(() => {
        // Remove the "session" cookie by setting its expiration to the past
        document.cookie = "sessionId=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
        navigate("/login");
      });
  }, [navigate, setToken]);

  return null; 
}

export default Logout;
