import { useState } from "react";
import { Link, useNavigate, useOutletContext } from "react-router-dom"
import { displayErrorMessage } from "./ErrorMessage";

function Login () {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

    const { setToken } = useOutletContext()

    const navigate = useNavigate();

    const handleSubmit = (e) => {
        e.preventDefault();

        const userData = {
            email: email,
            password: password
        };

        const headers = new Headers()
        headers.append("Content-Type", "application/json")
  
        let requestOptions = {
          body: JSON.stringify(userData),
          method: "POST",
          headers: headers,
          credentials: "include",
        };

        fetch("http://localhost:8080/login", requestOptions)
        .then((response) => response.json())
        .then((data) => {
          if (data.error) {
            displayErrorMessage(data.message); 
          } else {
            setToken(data.access_token)
            navigate("/");
          }
        })
        .catch(error => {
          displayErrorMessage(`An error occurred while logging in: ${error.message}`);
        })
    }
  

    return (
        <div className="auth-form-container">
            <h2>Login</h2>
            <form className="login-form" onSubmit={handleSubmit}>
                <label htmlFor="email">email</label>
                <input value={email} onChange={(e) => setEmail(e.target.value)}type="email" placeholder="youremail@gmail.com" id="email" name="email" />
                <label htmlFor="password">password</label>
                <input value={password} onChange={(e) => setPassword(e.target.value)} type="password" placeholder="********" id="password" name="password" />
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