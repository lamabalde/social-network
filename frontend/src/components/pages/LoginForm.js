import React, { useState, useContext, useEffect } from "react";
import { Link, useNavigate } from "react-router-dom";
import Form from '../UI/Form';
import FormInput from "../UI/FormInput";
import FormLabel from "../UI/FormLabel";
import LgButton from "../UI/LgButton";
import { AuthContext } from "../store/auth-context";
import styles from "./LoginForm.module.css";

const LoginForm = () => {
    const [enteredEmail, setEnteredEmail] = useState("");
    const [enteredPw, setEnteredPw] = useState("");
    const [loginErrMsg, setLoginErrMsg] = useState("");
    const navigate = useNavigate();
    const ctx = useContext(AuthContext);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState(null);

    useEffect(() => {
        setLoginErrMsg(ctx.errMsg);
        navigate("/login", { replace: true });
    }, [ctx.errMsg]);

    useEffect(() => {
        setIsLoading(ctx.loginIsLoading);
        setError(ctx.loginError);
    }, [ctx.loginIsLoading, ctx.loginError]);

    const emailChangeHandler = (e) => {
        setEnteredEmail(e.target.value);
        console.log(enteredEmail);
    };
    const pwChangeHandler = (e) => {
        setEnteredPw(e.target.value);
        console.log(enteredPw);
    };

    const submitHandler = (e) => {
        e.preventDefault();
        const loginPayloadObj = {
            email: enteredEmail,
            pw: enteredPw
        };
        console.log(loginPayloadObj);
        ctx.onLogin(loginPayloadObj);
        setEnteredEmail("");
        setEnteredPw("");

        ctx.setErrMsg("");

        navigate("/", { replace: true })
    };

    return (
        <div className={styles.container}>
            <h1 className={styles["title"]}>🅻🅾🅶🅸🅽</h1>
            <h2>{loginErrMsg}</h2>
            {!isLoading && <Form className={styles["login-form"]} onSubmit={submitHandler}>
                <FormLabel htmlFor="email">Email</FormLabel>
                <FormInput className={styles["login-input"]} name="email" id="email" placeholder="abc@gmail.com" value={enteredEmail} onChange={emailChangeHandler} />
                <FormLabel htmlFor="password">Password</FormLabel>
                <FormInput className={styles["login-input"]} type="password" name="password" id="password" placeholder="Password" value={enteredPw} onChange={pwChangeHandler} />
                <LgButton className={styles["sub-btn"]} type="submit">Login</LgButton>
                <p>Don't have an account? <Link to={"/reg"}>Register</Link></p>
            </Form>}
            {!isLoading && error && <h2>{error}</h2>}
            {isLoading && <h2>Logging in...</h2>}
        </div>
    );
};

export default LoginForm;