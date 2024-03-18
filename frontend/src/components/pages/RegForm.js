   import React, { useEffect, useState, useContext } from "react";
import { Link, useNavigate } from "react-router-dom";
import Form from '../UI/Form';
import FormInput from "../UI/FormInput";
import FormLabel from "../UI/FormLabel";
import FormTextarea from "../UI/FormTextarea";
import LgButton from "../UI/LgButton";
import ImgUpload from "../UI/ImgUpload";
import LoadingSpinner from "../UI/LoadingSpinner";
import { AuthContext } from "../store/auth-context";
import styles from "./RegForm.module.css";

const RegForm = () => {
    const imageSrc = "../../images/";
    let defaultImagePath = "user.png";

    const ctx = useContext(AuthContext);

    const [enteredEmail, setEnteredEmail] = useState("");
    const [enteredPw, setEnteredPw] = useState("");
    const [enteredFName, setEnteredFName] = useState("");
    const [enteredLName, setEnteredLName] = useState("");
    const [enteredDob, setEnteredDob] = useState("");
    const [uploadedImg, setUploadedImg] = useState("");
    const [enteredNickname, setEnteredNickname] = useState("");
    const [enteredAbout, setEnteredAbout] = useState("");
    // const [regErrMsg, setRegErrMsg] = useState("");
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState(null);

    const navigate = useNavigate();

    useEffect(() => {
        ctx.regSuccess && navigate("/login", { replace: true })
    }, [ctx.regSuccess]);

    // useEffect(() => {
    //     setRegErrMsg(ctx.errMsg);
    //     // navigate("/reg", { replace: true });
    // }, [ctx.errMsg]);

    useEffect(() => {
        setIsLoading(ctx.regIsLoading);
        setError(ctx.regError);
    }, [ctx.regIsLoading, ctx.regError]);

    const emailChangeHandler = (e) => {
        setEnteredEmail(e.target.value);
        console.log(enteredEmail);
    };
    const pwChangeHandler = (e) => {
        setEnteredPw(e.target.value);
        console.log(enteredPw);
    };
    const fNameChangeHandler = (e) => {
        setEnteredFName(e.target.value);
        console.log(enteredFName);
    };
    const lNameChangeHandler = (e) => {
        setEnteredLName(e.target.value);
        console.log(enteredLName);
    };
    const dobChangeHandler = (e) => {
        setEnteredDob(e.target.value);
        console.log(enteredDob);
    };
    const avatarHandler = (e) => {
        const file = e.target.files[0];
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.addEventListener("load", () => {
            console.log(reader.result);
            setUploadedImg(reader.result);
        })

        setUploadedImg(e.target.value);
        console.log(uploadedImg);
    };
    const nicknameChangeHandler = (e) => {
        setEnteredNickname(e.target.value);
        console.log(enteredNickname);
    };
    const aboutChangeHandler = (e) => {
        setEnteredAbout(e.target.value);
        console.log(enteredAbout);
    };

    const submitHandler = (e) => {
        e.preventDefault();
        const regPayloadObj = {
            email: enteredEmail,
            pw: enteredPw,
            fname: enteredFName,
            lname: enteredLName,
            Dob: enteredDob,
            avatar: uploadedImg,
            nname: enteredNickname,
            about: enteredAbout,
        };
        console.log(regPayloadObj);

        ctx.onReg(regPayloadObj);

        setEnteredEmail("");
        setEnteredPw("");
        setEnteredFName("");
        setEnteredLName("");
        setEnteredDob("");
        setUploadedImg("");
        setEnteredNickname("");
        setEnteredAbout("");
        // navigate("/login", {replace: true});

        ctx.setErrMsg("");
    };

    const resetHandler = () => {
        setError(false);
    };

    return (
        <div className={`${styles.container} ${styles.scrollable}`}>
            <h1 className={styles["title"]}>🆁🅴🅶🅸🆂🆃🅴🆁</h1>
            {console.log(isLoading)}
            {!isLoading && error && 
            <>
            {error && <h2 className={styles["error-msg"]}>{error}</h2>}
            <div  className={styles["try-again"]}><LgButton onClick={resetHandler}>Try Again</LgButton></div>
            </>}
            {!isLoading && !error && <Form className={`${styles["reg-form"]} ${styles["reg-form-scrollable"]}`} onSubmit={submitHandler}>
                <FormLabel className={styles["reg-label"]} htmlFor="email">Email</FormLabel>
                <FormInput className={styles["reg-input"]} type="email" name="email" id="email" placeholder="abc@gmail.com" value={enteredEmail} onChange={emailChangeHandler} />
                <FormLabel className={styles["reg-label"]} htmlFor="password">Password</FormLabel>
                <FormInput className={styles["reg-input"]} type="password" name="password" id="password" placeholder="Password" value={enteredPw} onChange={pwChangeHandler} />
                <FormLabel className={styles["reg-label"]} htmlFor="fname">First Name</FormLabel>
                <FormInput className={styles["reg-input"]} type="text" name="fname" id="fname" placeholder="Aissata" value={enteredFName} onChange={fNameChangeHandler} />
                <FormLabel className={styles["reg-label"]} htmlFor="lname">Last Name</FormLabel>
                <FormInput className={styles["reg-input"]} type="text" name="lname" id="lname" placeholder="Sow" value={enteredLName} onChange={lNameChangeHandler} />
                <FormLabel className={styles["reg-label"]} htmlFor="Dob">Date of Birth</FormLabel>
                <FormInput className={styles["reg-input"]} type="date" name="Dob" id="Dob" value={enteredDob} onChange={dobChangeHandler} />
               
                
                <FormLabel className={styles["reg-label"]} htmlFor="nname">Nickname (Optional)</FormLabel>
                <input className={styles["reg-input"]} type="text" name="nname" id="nname" placeholder="BledBoy" value={enteredNickname} onChange={nicknameChangeHandler} />
                <FormLabel className={styles["reg-label"]} >Avatar (Optional)</FormLabel>
                <figure>
                    {!uploadedImg && <img src={require("../../images/" + defaultImagePath)} alt="Default Image" className={styles["img-preview"]} width={"250px"} />}
                    {uploadedImg && <img src={uploadedImg} className={styles["img-preview"]} width={"250px"} />}
                </figure>
                <ImgUpload className={styles["reg-input"]} name="avatar" id="avatar" accept=".jpg, .jpeg, .png, .gif" text="Upload Image" onChange={avatarHandler} />
                <FormLabel className={styles["reg-label"]} htmlFor="about">About Me (Optional)</FormLabel>
                {/* <EmojiPicker /> */}
                <FormTextarea className={styles["reg-input"]} name="about" id="about" placeholder="About me..." rows={5} value={enteredAbout} onChange={aboutChangeHandler} />
                <LgButton className={styles["sub-btn"]} type="submit">Register</LgButton>
                <p>Already have an account? <Link to="/login">Login</Link></p>
            </Form>}
            {isLoading && <div>
                <div className={styles["loading-spinner-div"]}><LoadingSpinner/></div>
                <h2 className={styles["loading"]}>Registering New User...</h2>
            </div> }
        </div>

    )
};

export default RegForm;