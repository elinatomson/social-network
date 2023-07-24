import { Link } from 'react-router-dom'

function Home () {
  return (
    <div>
      <Link className="button" to="/register">Register</Link>
      <Link className="button" to="/login">Log In</Link>
    </div>
  )
}

export default Home