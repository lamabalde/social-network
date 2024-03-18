import { useState } from "react";
import AllGroups from "../group/AllGroups";
import AllJoinedGroups from "../group/AllJoinedGroups";
import CreateGroup from "../group/CreateGroup";
import JoinedGroups from "../group/JoinedGroup";

// import classes from './GroupPage.module.css';
import classes from './layout.module.css';

function GroupPage() {
    const [refresh, setRefresh] = useState(false)


    function creategroupupdate() {
        refresh ? setRefresh(false) : setRefresh(true);

    }

    return <div className={classes.container}>
        <div className={classes.mid}>
            {/* <CreateGroup></CreateGroup> */}

            <AllGroups refresh={refresh}></AllGroups>
        </div>
        <div className={classes.right}>
            <CreateGroup onnewgroup={creategroupupdate}  ></CreateGroup>
            {/* <JoinedGroups></JoinedGroups> */}
        </div>

    </div>

}

export default GroupPage;