import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { displayErrorMessage } from "./ErrorMessage";
import { Link } from 'react-router-dom';

function Groups() {
    const [showGroups, setShowGroups] = useState(false);
    const [groups, setGroups] = useState([]);
    const navigate = useNavigate();
    const location = useLocation();
    const { groupContent } = location.state || {};

    const handleToggleGroups = () => {
        setShowGroups(!showGroups);
    };

    useEffect(() => {
        fetch('/all-groups')
        .then((response) => response.json())
        .then((data) => {
            if (data) {
                setGroups(data);
            }
        })
        .catch((error) => {
            displayErrorMessage(`${error.message}`);
        });
        }, [navigate, groupContent]);


    return (
        <div className="users">
            <div className="chat_users" onClick={handleToggleGroups} >All Social Network groups</div>
            {showGroups && (
            <div>
                <div id="error" className="alert"></div>
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
            )}
        </div>
    );
}
export default Groups;
