import './styles.css';
import React from 'react'
import ReactDOM from 'react-dom/client'
import { RouterProvider, createBrowserRouter } from 'react-router-dom'
import ErrorPage from './components/ErrorPage'
import MainPage from './pages/MainPage'
import Home from './pages/Home'
import Register  from "./pages/Register"
import Login from "./pages/Login"
import Logout from "./pages/Logout"
import Profile from "./pages/Profile"
import User from "./pages/User"
import Group from "./pages/Group"
import Event from "./pages/Event"

const router = createBrowserRouter([
  {
    path: "/",
    errorElement: <ErrorPage />,
    children: [
      {index: true, element: <Home />},
      {
        path: "/register", 
        element: <Register />,
      },
      {
        path: "/login",
        element: <Login />,
      },
    ]
  },
  {
    path: "/main",
    element: <MainPage />,
  },
  {
    path: "/user/:userId",
    element: <User />,
  },
  {
    path: "/profile",
    element: <Profile />,
  },
  {
    path: "/group/:groupId",
    element: <Group />,
  },
  {
    path: "/group-event/:eventId",
    element: <Event />,
  },
  {
    path: "/logout",
    element: <Logout />,
  },
])

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
)
