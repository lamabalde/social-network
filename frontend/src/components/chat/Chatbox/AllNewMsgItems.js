import NewMsgItem from "./NewMsgItem";
const AllNewMsgItems = (props) => {
    console.log("msg in AllNewMsgItems", props.newMsgItems);
    return (
        props.newMsgItems.map((newMsg) => {
            return <NewMsgItem
                key={newMsg.id}
                id={newMsg.id}
                targetid={newMsg.targetid}
                sourceid={newMsg.sourceid}
                msg={newMsg.message}
                createdat={newMsg.createdat}
            />
        })
    );
};

export default AllNewMsgItems;