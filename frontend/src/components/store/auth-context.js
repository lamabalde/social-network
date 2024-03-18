import React, { useState, useEffect, useContext } from "react";
import { UsersContext } from "./users-context";

export const AuthContext = React.createContext({
  isLoggedIn: false,
  onReg: (regPayloadObj) => {},
  onLogin: (loginPayloadObj) => {},
  onLogout: () => {},
  regSuccess: false,
  notif: [],
  errMsg: "",
  setErrMsg: () => {},
  regIsLoading: false,
  regError: null,
  loginIsLoading: false,
  loginError:null,
});

export const AuthContextProvider = (props) => {
  const [loggedIn, setLoggedIn] = useState(false);
  const [regSuccess, setRegSuccess] = useState(false);
  const [notif, setNotif] = useState([])
  const [errMsg, setErrMsg] = useState("");
  const loginURL = "http://localhost:8080/login";
  const regURL = "http://localhost:8080/reg";
  const logoutURL = "http://localhost:8080/logout";

  const usersCtx = useContext(UsersContext);

  const [regIsLoading, setRegIsLoading] = useState(false);
  const [regError, setRegError] = useState(null);

  const regHandler = (regPayloadObj) => {
    setRegIsLoading(true);
    setRegError(null);

    console.log("app.js", regPayloadObj);

    const reqOptions = {
      method: "POST",
      credentials: "include",
      mode: "cors",
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(regPayloadObj)
    };
    
    // fetch(regURL, reqOptions)
    //   .then(resp => {
    //     if (!resp.ok) throw new Error ("Failed to Register");
    //     return resp.json();
    //   })
    //   .then(data => {
    //     console.log(data);
    //     // redirect to login
    //     if (data.success) {
    //       console.log(data.success);
    //       setRegSuccess(true);
    //       usersCtx.onNewUserReg();
    //     } else {
    //       setRegSuccess(false);
    //       throw new Error ("Email Already Registered");
    //     }
    //     setRegIsLoading(false);
    //   })
    //   .catch(err => {
    //     // console.log(err);
    //     setRegError(err.message);
    //     setRegIsLoading(false);
    //   })
    setTimeout(() => {
      fetch(regURL, reqOptions)
      .then(resp => {
        if (!resp.ok) throw new Error ("Failed to Register");
        return resp.json();
      })
      .then(data => {
        console.log(data);
        // redirect to login
        if (data.success) {
          console.log(data.success);
          setRegSuccess(true);
          usersCtx.onNewUserReg();
        } else {
          setRegSuccess(false);
          throw new Error ("Email Already Registered");
        }
        setRegIsLoading(false);
      })
      .catch(err => {
        // console.log(err);
        setRegError(err.message);
        setRegIsLoading(false);
      })
    }, 2000);
  };

  const [loginIsLoading, setLoginIsLoading] = useState(false);
  const [loginError, setLoginError] = useState(null);

  const loginHandler = (loginPayloadObj) => {
    setLoginIsLoading(true);
    setLoginError(null);

    const reqOptions = {
      method: "POST",
      credentials: "include",
      mode: "cors",
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(loginPayloadObj)
    };
    setTimeout(() => {
      fetch(loginURL, reqOptions)
      .then(resp => {
        if (!resp.ok) throw new Error("Failed to Login");
        return resp.json();
      })
      .then(data => {
        console.log("data reached", data);
        if (data.resp.success) {
          setLoggedIn(true);
          console.log("login resp success", data);
          localStorage.setItem("user_id", data.resp.user_id);
          localStorage.setItem("fname", data.resp.fname);
          localStorage.setItem("lname", data.resp.lname);
          localStorage.setItem("dob", data.resp.dob);
          data.resp.nname && localStorage.setItem("nname", data.resp.nname);
          data.resp.avatar && localStorage.setItem("avatar", data.resp.avatar);
          data.resp.about && localStorage.setItem("about", data.resp.about);
          localStorage.setItem("public", data.resp.public);
          localStorage.setItem("new_notif", "[]");

          data.notif && localStorage.setItem("new_notif", JSON.stringify(Object.values(data.notif)));
          data.notif && setNotif(Object.values(data.notif));

        } else {
          setLoggedIn(false)
          console.log("login resp failed", data);
          // setErrMsg("ERROR - Please check your credentials")
          throw new Error("Login Failed: Please check your credentials and login again");
        }
        setLoginIsLoading(false);
      })
      .catch(err => {
        // console.log(err);
        console.log("lerrmsg", err.message);
        setLoginError(err.message);
        setLoginIsLoading(false);
      })
    }, 2000);
    // fetch(loginURL, reqOptions)
    //   .then(resp => {
    //     if (!resp.ok) throw new Error("Failed to Login");
    //     return resp.json();
    //   })
    //   .then(data => {
    //     console.log("login", data);
    //     if (data.resp.success) {
    //       setLoggedIn(true);
    //       localStorage.setItem("user_id", data.resp.user_id);
    //       localStorage.setItem("fname", data.resp.fname);
    //       localStorage.setItem("lname", data.resp.lname);
    //       localStorage.setItem("dob", data.resp.dob);
    //       data.resp.nname && localStorage.setItem("nname", data.resp.nname);
    //       data.resp.avatar && localStorage.setItem("avatar", data.resp.avatar);
    //       data.resp.about && localStorage.setItem("about", data.resp.about);
    //       localStorage.setItem("public", data.resp.public);
    //       localStorage.setItem("new_notif", "[]");

    //       localStorage.setItem("new_notif", JSON.stringify(Object.values(data.notif)));
    //       setNotif(Object.values(data.notif));

    //     } else {
    //       setLoggedIn(false)
    //       // setErrMsg("ERROR - Please check your credentials")
    //       throw new Error("Login Failed: Please check your credentials and login again");
    //     }
    //     setLoginIsLoading(false);
    //   })
    //   .catch(err => {
    //     // console.log(err);
    //     setLoginError(err.message);
    //     setLoginIsLoading(false);
    //   })
  };

  const logoutHandler = () => {
    localStorage.clear();
    localStorage.removeItem("following");

    const reqOptions = {
      method: "GET",
      credentials: "include",
      mode: "cors",
      headers: {
        'Content-Type': 'application/json'
      }
    };
    fetch(logoutURL, reqOptions)
      .then(resp => resp.json())
      .then(data => console.log(data))
      .catch(err => console.log(err))

    setLoggedIn(false);

    // in case the next user wants to reg
    setRegSuccess(false);
  };

  useEffect(() => { localStorage.getItem("user_id") && setLoggedIn(true) }, []);

  return (
    <AuthContext.Provider
      value={{
        isLoggedIn: loggedIn,
        onReg: regHandler,
        onLogin: loginHandler,
        onLogout: logoutHandler,
        regSuccess: regSuccess,
        notif: notif,
        errMsg: errMsg,
        setErrMsg: setErrMsg,
        regIsLoading: regIsLoading,
        regError: regError,
        loginIsLoading: loginIsLoading,
        loginError: loginError,
      }}
    >
      {props.children}
    </AuthContext.Provider>
  );
};