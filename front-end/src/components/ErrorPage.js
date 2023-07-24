import { useRouteError } from "react-router-dom";
import Error from './../images/error.jpg';

export default function ErrorPage() {
    const error = useRouteError();

    return (
        <div className="App">
            <div className="error">
            <img className="image" src={Error} alt="error"></img><br></br>
                Sorry, an unexpected error has occured: <br></br>
                {error.statusText || error.message}
            </div>
        </div>
    )
}