import './styles.css';
import { useState } from 'react';
import { Link, Outlet } from 'react-router-dom';

function App() {
  const [token, setToken] = useState("")

  return (
    <div>
      <div>
        <Link class="link" to="/">
          <span className="heading1">Welcome to </span>
          <span className="heading2">Social Network page</span>
        </Link>
      </div>
      <Outlet context ={{
        token, setToken, }}
        />
    </div>
  );
}

export default App;
