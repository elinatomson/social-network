import { useState } from "react";
import { Link, useNavigate } from "react-router-dom"
import { displayErrorMessage } from "../components/ErrorMessage";
import Footer from "../components/Footer";

function Register () {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [firstName, setFirstName] = useState("");
    const [lastName, setLastName] = useState("");
    const [dateOfBirth, setDateOfBirth] = useState("");
    const [avatar, setAvatar] = useState(null);
    const [nickname, setNickname] = useState("");
    const [aboutMe, setAboutMe] = useState("");
    const [errors, setErrors] =useState([])
    const navigate = useNavigate()

    const handleSubmit = (e) => {
      e.preventDefault();

      let errors = []
      let required = [
        { field: email, name: "email"},
        { field: password, name: "password"},
        { field: firstName, name: "first_name"},
        { field: lastName, name: "last_name"},
        { field: dateOfBirth, name: "date_of_birth"},
      ]

      required.forEach(function (obj) {
        if (obj.field === "" ) {
          errors.push(obj.name);
        }
      })

      if (password.length < 5) {
        errors.push("password");
      }

      if (aboutMe.length > 100) {
        errors.push("about_me");
      }

      if (nickname.length > 10) {
        errors.push("nickname");
      }

      setErrors(errors)

      if (errors.length > 0) {
        return;
      }

      const formData = new FormData();
      formData.append("email", email);
      formData.append("password", password);
      formData.append("first_name", firstName);
      formData.append("last_name", lastName);
      formData.append("date_of_birth", dateOfBirth);
      formData.append("avatar", avatar); 
      formData.append("nickname", nickname);
      formData.append("about_me", aboutMe);

      const headers = new Headers()

      let requestOptions = {
        body: formData,
        method: "POST",
        headers: headers,
      }

      fetch("/register", requestOptions)
      .then((response) => {
        if (!response.ok) {
          return response.json().then((data) => {
            throw new Error(data.message);
          });
        } else {
          return response.json();
        }
      })
      .then(() => {
        navigate("/login");
      })
      .catch((error) => {
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
                <h2>Register</h2>
                <h3>Password has to be 5 letters long!</h3>
                <form className="register-form" onSubmit={handleSubmit}>
                    <label htmlFor="email">Email*</label>
                    <input value={email} onChange={(e) => setEmail(e.target.value)} type="email" placeholder="youremail@gmail.com" id="email" name="email" />
                      {errors.includes("email") && (
                        <p className="alert">Please fill in the email.</p>
                      )}
                    <label htmlFor="password">Password*</label>
                    <input value={password} onChange={(e) => setPassword(e.target.value)} type="password" placeholder="********" id="password" name="password" />
                      {errors.includes("password") && (
                        <p className="alert">Please fill in the password (at least 5 letters).</p>
                      )}
                    <label htmlFor="firstName">First Name*</label>
                    <input value={firstName} onChange={(e) => setFirstName(e.target.value)} type="text" placeholder="First name" id="firstName" name="first_name" />
                      {errors.includes("first_name") && (
                        <p className="alert">Please fill in the first name.</p>
                      )}
                    <label htmlFor="lastName">Last Name*</label>
                    <input value={lastName} onChange={(e) => setLastName(e.target.value)} type="text" placeholder="Last name" id="lastName" name="last_name" />
                      {errors.includes("last_name") && (
                        <p className="alert">Please fill in the last name.</p>
                      )}
                    <label htmlFor="dateOfBirth">Date of Birth*</label>
                    <input value={dateOfBirth} onChange={(e) => setDateOfBirth(e.target.value)} type="date" id="dateOfBirth" name="date_of_birth" max={new Date().toISOString().split('T')[0]}/>
                      {errors.includes("date_of_birth") && (
                        <p className="alert">Please select a date of birth.</p>
                      )}
                    <label htmlFor="avatar">Avatar/Image (Optional)</label>
                    <input onChange={(e) => setAvatar(e.target.files[0])} type="file" accept="image/*" id="avatar" name="avatar" />
                    <label htmlFor="nickname">Nickname (Optional)</label>
                    <input value={nickname} onChange={(e) => setNickname(e.target.value)} type="text" placeholder="Nickname" id="nickname" name="nickname"/>
                    {errors.includes("nickname") && (
                        <p className="alert">Too long Nickname (make it less than 10 letters).</p>
                      )}
                    <label htmlFor="aboutMe">About Me (Optional)</label>
                    <input value={aboutMe} onChange={(e) => setAboutMe(e.target.value)} placeholder="Something about yourself..." id="about_me" name="aboutMe"/>
                    {errors.includes("about_me") && (
                        <p className="alert">Too long About Me (make it less than 100 letters).</p>
                      )}
                    <div id="error" className="alert"></div>
                    <button className="button" type="submit">Register</button>
                    <Link className="button" to="/" type="submit">Cancel</Link>
                </form>
                <div>
                    <button className="link-btn" onClick={() => navigate("/login")}> Already have an account? Click here to log in!</button>
                </div>
            </div>
          </div>
          <div className="right-container"></div>
        </div>
      <Footer />
      </div>
    </div>
  );
}

export default Register