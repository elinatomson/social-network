import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.js'
import { RouterProvider, createBrowserRouter } from 'react-router-dom'
import ErrorPage from './components/ErrorPage'
import Home from './components/Home'
import Register  from "./components/Register"
import Login from "./components/Login"
import Logout from "./components/Logout"
import Profile from "./components/Profile"
import Social from "./components/Social"
import Search from "./components/Search"
import CreatePost from "./components/CreatePost"
import AllPosts from "./components/AllPosts"
import UserActivity from "./components/UserActivity.js"

const router = createBrowserRouter([
  {
    path: "/",
    element:<App />,
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
      {
        path: "/logout",
        element: <Logout />,
      },
      {
        path: "/profile",
        element: <Profile />,
      },
      {
        path: "/activity",
        element: <UserActivity />,
      },
      {
        path: "/social",
        element: <Social />,
      },
      {
        path: "/search",
        element: <Search />,
      },
      {
        path: "/create-post",
        element: <CreatePost />,
      },
      {
        path: "/all-posts",
        element: <AllPosts />,
      },
    ]
  }
])

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
)
