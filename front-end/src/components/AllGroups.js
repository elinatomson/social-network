import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { displayErrorMessage } from "./ErrorMessage";
import { Link } from 'react-router-dom';
import WebSocketComponentForGroup from '../components/WebsocketForGroup'; 

function Groups() {
    const [showGroups, setShowGroups] = useState(false);
    const [groups, setGroups] = useState([]);
    const [groupName, setGroupName] = useState(null);
    const [firstNameFrom, setFirstNameFrom] = useState(null);
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
            const currentUser = data.current_user;
            setFirstNameFrom(currentUser.first_name);
            setGroups(data.groups);
        }
    })
    .catch((error) => {
        displayErrorMessage(`${error.message}`);
    });
    }, [navigate, groupContent]);

    const handleGroupClick = (group) => {
        if (groupName === group) {
          setGroupName(null); 
        } else {
          setGroupName(group);
        }
    };

    const handleChatClose = () => {
        setGroupName(null); 
    };


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
                                <Link className="link-btn" onClick={() => handleGroupClick(group)}>
                                    {group.title}
                                </Link>
                            </div>
                        ))}
                    </div>
                )}
            </div>
            )}
            <div className="chat">
            {groupName && <WebSocketComponentForGroup groupName={groupName} firstNameFrom={firstNameFrom} closeChat={handleChatClose}/>}
      </div>
        </div>
    );
}
export default Groups;
