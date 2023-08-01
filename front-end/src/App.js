import './styles.css';
import { useState } from 'react';
import { Outlet } from 'react-router-dom';
import Footer from "./components/Footer"
import Header from "./components/Header"

function App() {
  const [token, setToken] = useState("")

  return (
    <div className="app-container">
      <div className="content">
        <Header />
        <Outlet context ={{
          token, setToken, }}
          />
      </div>
      <Footer />
    </div>
  );
}

export default App;
