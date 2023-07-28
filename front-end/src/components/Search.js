import React, { useState } from 'react';

function Search({ setSearchResults }) {
  const [searchTerm, setSearchTerm] = useState('');

  const handleSearch = (query) => {
    if (query === '') {
      setSearchResults(null); 
      return;
    }

    fetch(`/search?query=${query}`)
      .then((response) => response.json())
      .then((data) => {
        setSearchResults(data.length > 0 ? data : []);
      })
      .catch((error) => {
        console.error("Failed to perform search:", error);
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
      <input type="text" placeholder="Search users..." value={searchTerm} onChange={handleChange} />
    </div>
  );
}

export default Search;
