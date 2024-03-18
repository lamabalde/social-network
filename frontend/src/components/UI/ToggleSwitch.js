import { useEffect, useState } from "react";
import ProfilePosts from "../profile/profilePost";
import classes from "./ToggleSwitch.module.css";
  
const ToggleSwitch = ({ label, value, onChange, onClick }) => {

  const [isChecked, setIsChecked] = useState(localStorage.getItem('isChecked') === 'true');

  useEffect(() => {
    localStorage.setItem('isChecked', isChecked);
  }, [isChecked]);




  return (
    <div className={classes.container}>
      <div className={classes.labelWrapper}>
        {label}{""}
    
      <div className={classes.toggleSwitch}>
        <input type="checkbox" className={classes.checkbox} name={label} id={value} value={value}  checked={isChecked}
        onChange={() => setIsChecked(!isChecked)} onClick={onClick}/>
        <label className={classes.value} htmlFor={value}>
          <span className={classes.inner} />
          <span className={classes.switch} />
        </label>
        </div>
      </div>
    </div>
  );
};
  
export default ToggleSwitch;