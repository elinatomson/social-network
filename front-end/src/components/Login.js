import { useState, useEffect } from "react";
import { Link, useNavigate } from "react-router-dom"
import { displayErrorMessage } from "./ErrorMessage";

function Login () {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

    const navigate = useNavigate();

    const token = document.cookie
    .split("; ")
    .find((row) => row.startsWith("sessionId="))
    ?.split("=")[1];

    useEffect(() => {
      if (token) {
        navigate("/main");
      }
    }, [token, navigate]);


    const handleSubmit = (e) => {
      e.preventDefault();

      const userData = {
        email: email,
        password: password,
      };

    const headers = new Headers();
    headers.append("Content-Type", "application/json");
    headers.append("Authorization", token);

    let requestOptions = {
      body: JSON.stringify(userData),
      method: "POST",
      headers: headers,
    };

    fetch("http://localhost:8080/login", requestOptions)
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
        document.cookie = `sessionId=${data.session}`
        navigate("/main");
    })
    .catch(error => {
        displayErrorMessage(`${error.message}`);
    });
  };

    return (
        <div className="auth-form-container">
            <h2>Login</h2>
            <form className="login-form" onSubmit={handleSubmit}>
                <label htmlFor="email">email</label>
                <input value={email} onChange={(e) => setEmail(e.target.value)}type="email" placeholder="youremail@gmail.com" id="email" name="email" required/>
                <label htmlFor="password">password</label>
                <input value={password} onChange={(e) => setPassword(e.target.value)} type="password" placeholder="********" id="password" name="password" required/>
                <div id="error" class="alert"></div>
                <button className="button" type="submit">Log In</button>
                <Link className="button" to="/" type="submit">Cancel</Link>
            </form>
            <div>
                <button className="link-btn" onClick={() => navigate("/register")}> Dont have an account? Click here to register!</button>
            </div>
        </div>
    )
}

export default Login