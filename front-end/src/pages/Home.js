import { Link } from 'react-router-dom'
import Header from '../components/Header';
import Footer from "../components/Footer";

function Home () {
  return (
    <div className="app-container">
      <div className="content">
        <Header />
          <div>
            <Link className="button" to="/register">Register</Link>
            <Link className="button" to="/login">Log In</Link>
          </div>
      </div>
      <Footer />
    </div>
  )
}

export default Home