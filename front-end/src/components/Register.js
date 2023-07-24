import { useState } from "react";
import { Link, useNavigate } from "react-router-dom"
import { displayErrorMessage } from "./ErrorMessage";

function Register () {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [firstName, setFirstName] = useState("");
    const [lastName, setLastName] = useState("");
    const [dateOfBirth, setDateOfBirth] = useState("");
    const [avatar, setAvatar] = useState("");
    const [nickname, setNickname] = useState("");
    const [aboutMe, setAboutMe] = useState("");
    const navigate = useNavigate()

    const handleSubmit = async (e) => {
      e.preventDefault();

      var passwordLength = password.length >= 5 && password.length <= 50;
      if (!passwordLength) {
        displayErrorMessage('Password has to be 5 letters long!');
        return;
      }

      const userData = {
          email: email,
          password: password,
          first_name: firstName,
          last_name: lastName,
          date_of_birth: dateOfBirth,
          avatar: avatar,
          nickname: nickname,
          about_me: aboutMe,
      };

      const headers = new Headers()
      headers.append("Content-Type", "application/json")

      let requestOptions = {
        body: JSON.stringify(userData),
        method: "POST",
        headers: headers,
      }

      fetch("http://localhost:8080/register", requestOptions)
      .then((response) => {
        if (response.ok) {
          navigate("/login");
        } else {
          return response.json(); 
        }
      })
      .then((errorMessage) => {
        if (errorMessage) {
          displayErrorMessage(errorMessage.error); 
        }
      })
      .catch((error) => {
        displayErrorMessage(`An error occurred while registering: ${error.message}`);
      });
  };

  return (
    <div className="auth-form-container">
        <h2>Register</h2>
        <h3>Password has to be 5 letters long!</h3>
        <form className="register-form" onSubmit={handleSubmit}>
            <label htmlFor="email">Email*</label>
            <input value={email} onChange={(e) => setEmail(e.target.value)} type="email" placeholder="youremail@gmail.com" id="email" name="email" required/>
            <label htmlFor="password">Password*</label>
            <input value={password} onChange={(e) => setPassword(e.target.value)} type="password" placeholder="********" id="password" name="password" required/>
            <label htmlFor="firstName">First Name*</label>
            <input value={firstName} onChange={(e) => setFirstName(e.target.value)} type="text" placeholder="John" id="firstName" name="first_name" required/>
            <label htmlFor="lastName">Last Name*</label>
            <input value={lastName} onChange={(e) => setLastName(e.target.value)} type="text" placeholder="Doe" id="lastName" name="last_name" required/>
            <label htmlFor="dateOfBirth">Date of Birth*</label>
            <input value={dateOfBirth} onChange={(e) => setDateOfBirth(e.target.value)} type="date" id="dateOfBirth" name="date_of_birth" required/>
            <label htmlFor="avatar">Avatar/Image (Optional)</label>
            <input value={avatar} onChange={(e) => setAvatar(e.target.value)} type="url" placeholder="https://example.com/avatar.jpg" id="avatar" name="avatar" />
            <label htmlFor="nickname">Nickname (Optional)</label>
            <input value={nickname} onChange={(e) => setNickname(e.target.value)} type="text" placeholder="Nickname" id="nickname" name="nickname"/>
            <label htmlFor="aboutMe">About Me (Optional)</label>
            <input value={aboutMe} onChange={(e) => setAboutMe(e.target.value)} placeholder="Something about yourself..." id="about_me" name="aboutMe"/>
            <div id="error" class="alert"></div>
            <button className="button" type="submit">Register</button>
            <Link className="button" to="/" type="submit">Cancel</Link>
        </form>
        <div>
            <button className="link-btn" onClick={() => navigate("/login")}> Already have an account? Click here to log in!</button>
        </div>
    </div>
  );
}

export default Register