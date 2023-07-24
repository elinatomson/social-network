import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.js'
import { RouterProvider, createBrowserRouter } from 'react-router-dom'
import ErrorPage from './components/ErrorPage'
import Home from './components/Home'
import Register  from "./components/Register"
import Login from "./components/Login"
import Profile from "./components/Profile"

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
        path: "/profile",
        element: <Profile />,
      },
    ]
  }
])

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
)
