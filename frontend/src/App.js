// import { useEffect, useState } from "react";
import { useContext } from "react";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { AuthContext } from "./components/store/auth-context";
import Root from "./components/pages/Root";
import Landingpage from './components/pages/Landingpage';
import LoginForm from './components/pages/LoginForm';
import RegForm from './components/pages/RegForm';
import PostsPage from './components/pages/PostsPage';
import GroupPage from "./components/pages/GroupPage";
import GroupProfilePage from "./components/pages/GroupProfilePage";
import ProfilePage from "./components/pages/ProfilePage";
import LoadingTestPage from "./components/pages/LoadingTestPage";

function App() {
  const authCtx = useContext(AuthContext);

  let router = createBrowserRouter([
    {path: "/", element: <Landingpage />},
    {path: "/login", element: <LoginForm/>},
    {path: "/reg", element: <RegForm/>},
    {path: "/loading", element: <LoadingTestPage/>}
  ]);

 if (authCtx.isLoggedIn) router = createBrowserRouter([
    {
      path: "/",
      element: <Root />,
      children: [
          {path: "/", element: <PostsPage />},
          {path: "/:userId", element: <PostsPage/>},
          {path: "/profile", element: <ProfilePage />},
          {path: "/group", element: <GroupPage />},
          {path: "/groupprofile", element: <GroupProfilePage />},
          {path: "/groups", element: <GroupPage />},
          {path: "/profile/:userId", element: <ProfilePage/>}
          // {path: "/user/:userId", element <UserProfilePage />},
      ],
    }
  ]);
  return (
      <RouterProvider router={router}/>
  );
}

export default App;
