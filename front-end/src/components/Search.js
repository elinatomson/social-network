import React, { useState, useEffect } from 'react';
import { displayErrorMessage } from "./ErrorMessage";
import { useNavigate } from "react-router-dom"

function Search({ setSearchResults }) {
  const [searchTerm, setSearchTerm] = useState('');

  const navigate = useNavigate();

  const token = document.cookie
  .split("; ")
  .find((row) => row.startsWith("sessionId="))
  ?.split("=")[1];

  useEffect(() => {
    if (!token) {
      navigate("/login");
    }
  }, [token, navigate]);

  const handleSearch = (query) => {
    if (query === '') {
      setSearchResults(null); 
      return;
    }

    fetch(`/search?query=${query}`)
      .then((response) => response.json())
      .then((data) => {
        setSearchResults(data !== null ? data : []);
      })
      .catch(error => {
        displayErrorMessage(`${error.message}`);
        setSearchResults([]);
      });
  };

  const handleChange = (e) => {
    const searchTerm = e.target.value;
    setSearchTerm(searchTerm);
    handleSearch(searchTerm);
  };

  return (
    <div>
      <input className="search" type="text" placeholder="Search users..." value={searchTerm} onChange={handleChange} />
    </div>
  );
}

export default Search;
