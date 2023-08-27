import React, { useState } from 'react';
import { displayErrorMessage } from "./ErrorMessage";

function Search({ setSearchResults }) {
  const [searchTerm, setSearchTerm] = useState('');

  const performSearch = (query) => {
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

  const handleChange = (event) => {
    const value = event.target.value;
    setSearchTerm(value);
    performSearch(value);
  };

  return (
    <div>
      <input className="search" type="text" placeholder="Search users..." value={searchTerm} onChange={handleChange} />
      <div id="error" className="alert"></div>
    </div>
  );
}

export default Search;
