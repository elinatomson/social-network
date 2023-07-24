import { Link, useNavigate } from 'react-router-dom'

function Profile () {

    const navigate = useNavigate();

    const logOut = () => {
      navigate("/")
    }
  

  return (
    <div>
      <Link onClick={ logOut } className="button" to="/">Log Out</Link>
    </div>
  )
}

export default Profile