import Group from "./Group";
import classes from './AllGroups.module.css';
import { useEffect, useState } from "react";

function AllGroups(props) {

  const currUserId = localStorage.getItem("user_id");
const [groupData, setGroupData] = useState([])
    useEffect(() => {
      async function fetchData() {
        try {
          const response = await fetch(`http://localhost:8080/group`);
          const data = await response.json();
          console.log("fklgkdlf",data.data)
          setGroupData(data.data);
        } catch (error) {
          console.log(error);
        }
      }
      fetchData();
    }, [props.refresh]);



    return <div className={classes.container}>
        {groupData && groupData.map((group) => (
         <Group
        key={group.id}
        grpid={group.id}
        title={group.title}
        creator={group.creator}
        description={group.description}
        // img={group.img}
        />
        ))}
            </div>
}

export default AllGroups;
