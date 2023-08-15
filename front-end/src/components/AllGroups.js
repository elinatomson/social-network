import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';

function Groups() {
  const [groups, setGroups] = useState([]);

  useEffect(() => {
    fetch('/all-groups')
    .then((response) => response.json())
    .then((data) => {
        if (data) {
            setGroups(data);
        }
    })
    .catch((error) => {
        console.error('Error fetching groups:', error);
    });
    }, []);


    return (
        <div className="users">
            <div className="following">All Social Network groups</div>
            {groups.length === 0 ? (
                <p className="user">No groups.</p>
            ) : (
                <div className="user">
                {groups.map((group) => (
                    <div key={group.group_id}>
                        <Link className="link-btn" to={`/group/${group.group_id}`}>
                        {group.title}
                        </Link>
                    </div>
                ))}
                </div>
            )}
        </div>
    );
}
export default Groups;
