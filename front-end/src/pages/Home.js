import { Link } from 'react-router-dom'
import Footer from "../components/Footer";

function Home () {
  return (
    <div className="app-container">
      <div className="home">
        <div className="container">
          <div className="left-container"></div>
          <div className="middle-container">
              <span className="heading1">Welcome to </span><br/>
              <span className="heading2">Social Network page</span>
            <div>
              <Link className="button" to="/register">Register</Link>
              <Link className="button" to="/login">Log In</Link>
            </div>
          </div>
          <div className="right-container"></div>
        </div>
      </div>
      <Footer />
    </div>
  )
}

export default Home