import './styles.css';
import { useState } from 'react';
import { Outlet } from 'react-router-dom';

function App() {
  const [token, setToken] = useState("")

  return (
    <div>
      <div>
          <span className="heading1">Welcome to </span>
          <span className="heading2">Social Network page</span>
      </div>
      <Outlet context ={{
        token, setToken, }}
        />
    </div>
  );
}

export default App;
