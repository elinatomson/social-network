import React, { useState } from "react";
import { displayErrorMessage } from "./ErrorMessage";
import { useNavigate } from "react-router-dom"

function CreateEvent({ groupId }) {
    const [showFields, setShowFields] = useState(false);
    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [time, setTime] = useState("");
    const [errors, setErrors] =useState([])
    const [eventContent, setEventContent] = useState("");

    const navigate = useNavigate();

    const handleTitleChange = (e) => {
        setTitle(e.target.value);
    };
    
    const handleDescriptionChange = (e) => {
        setDescription(e.target.value);
    };

    const handleTimeChange = (e) => {
      setTime(e.target.value);
    };

    const handleToggleFields = () => {
        setShowFields(!showFields);
    };

    const handleSubmit = (e) => {
        e.preventDefault();

    let errors = []
    let required = [
        { field: title, name: "title"},
        { field: description, name: "description"},
        { field: time, name: "time"},
    ]

    required.forEach(function (obj) {
      if (obj.field === "") {
        errors.push(obj.name);
      }
    })

    setErrors(errors)

    if (errors.length > 0) {
      return;
    }
    
    const eventData = {
    title: title,
    description: description,
    time: time,
    group_id: groupId, 
    };

    const headers = new Headers();
    headers.append("Content-Type", "application/json");

    let requestOptions = {
      body: JSON.stringify(eventData),
      method: "POST",
      headers: headers,
    }

    fetch("/create-event", requestOptions)
      .then((response) => {
        if (response.ok) {
          setEventContent("");
          setTitle("");
          setDescription("");
          setTime("");
          setShowFields(false);
          navigate(`/group/${groupId}`, { state: { eventContent } });
        } else {
          return response.json(); 
        }
      })
      .then((errorMessage) => {
        if (errorMessage) {
          displayErrorMessage(errorMessage.error); 
        }
      })
      .catch((error) => {
        displayErrorMessage(`${error.message}`);
      });
  };

  return (
    <div className="posting">
        <h2 className="center" onClick={handleToggleFields}>New event</h2>
        {showFields && (
            <form className="login-form" onSubmit={handleSubmit}>
                <input value={title} onChange={handleTitleChange} placeholder="Title" name="title"/>
                {errors.includes("title") && (
                    <p className="alert">Please fill in the title.</p>
                )}
                <input value={description} onChange={handleDescriptionChange} placeholder="Description" name="description"/>
                {errors.includes("description") && (
                    <p className="alert">Please fill in the description.</p>
                )}
                <input value={time} onChange={handleTimeChange} placeholder="Time" name="time" type="date" min={new Date().toISOString().split('T')[0]}/>
                {errors.includes("time") && (
                    <p className="alert">Please fill in the event time.</p>
                )}
                <div id="error" className="alert"></div>
                <div id="error" className="alert"></div>
                <button className="button" type="submit">Create event</button>
            </form>
        )}
    </div>
  );
}

export default CreateEvent;
