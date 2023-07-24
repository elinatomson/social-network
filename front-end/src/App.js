import './styles.css';
import { Link, Outlet } from 'react-router-dom';

function App() {

  return (
    <div>
      <div>
        <Link class="link" to="/">
          <span className="heading1">Welcome to </span>
          <span className="heading2">Social Network page</span>
        </Link>
      </div>
      <Outlet />
    </div>
  );
}

export default App;
