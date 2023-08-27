import { useState } from "react";
import { Link, useNavigate } from "react-router-dom"
import { displayErrorMessage } from "../components/ErrorMessage";
import Footer from "../components/Footer";

function Login () {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [errors, setErrors] =useState([])

    const navigate = useNavigate();

    const handleSubmit = (e) => {
      e.preventDefault();

      let errors = []
      let required = [
        { field: email, name: "email"},
        { field: password, name: "password"},
      ]

      required.forEach(function (obj) {
        if (obj.field === "") {
          errors.push(obj.name);
        }
      })

      setErrors(errors)

      if (errors.length > 0) {
        return;
      }

      const userData = {
        email: email,
        password: password,
      };

    const headers = new Headers();
    headers.append("Content-Type", "application/json");

    let requestOptions = {
      body: JSON.stringify(userData),
      method: "POST",
      headers: headers,
    };

    fetch("/login", requestOptions)
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
    <div className="app-container">
      <div className="home">
        <div className="container">
          <div className="left-container"></div>
          <div className="middle-container">
            <span className="heading1">Welcome to </span><br/>
            <span className="heading2">Social Network page</span>
            <div className="auth-form-container">
                <h2>Login</h2>
                <form className="login-form" onSubmit={handleSubmit}>
                    <label htmlFor="email">email</label>
                    <input value={email} onChange={(e) => setEmail(e.target.value)}type="email" placeholder="youremail@gmail.com" id="email" name="email"/>
                      {errors.includes("email") && (
                        <p className="alert">Please fill in the email.</p>
                      )}
                    <label htmlFor="password">password</label>
                    <input value={password} onChange={(e) => setPassword(e.target.value)} type="password" placeholder="********" id="password" name="password"/>
                      {errors.includes("password") && (
                        <p className="alert">Please fill in the password.</p>
                      )}
                    <div id="error" className="alert"></div>
                    <button className="button" type="submit">Log In</button>
                    <Link className="button" to="/" type="submit">Cancel</Link>
                </form>
                <div>
                    <button className="link-btn" onClick={() => navigate("/register")}> Dont have an account? Click here to register!</button>
                </div>
            </div>
          </div>
          <div className="right-container"></div>
        </div>
      </div>
    <Footer />
  </div>
  )
}

export default Login