import React from "react";
import OldMsgItem from "./OldMsgItem";

const AllOldMsgItems = (props) => {

    console.log("msg in AllOldMsgItems", props.oldMsgItems);
    
    // console.log("is array in AllOldMsgItems", Array.isArray(props.oldMsgItems));
    return (
        props.oldMsgItems.map((oldMsg) => {
            return <OldMsgItem
                key={oldMsg.id}
                id={oldMsg.id}
                targetid={oldMsg.targetid}
                sourceid={oldMsg.sourceid}
                // groupid={oldMsg.groupid}
                msg={oldMsg.message}
                createdat={oldMsg.createdat}
            />
        })
    );
}

export default React.memo(AllOldMsgItems);