import React, {useState, useEffect} from "react";

export const GroupsContext = React.createContext({
    groups: [],
    onNewGroupCreated: () => {},
});

export const GroupsContextProvider = (props) => {
    const [groupsList, setGroupsList] = useState([]);

    // get groups
    const getGroupsHandler = () => {
        const groupUrl = "http://localhost:8080/group";
        fetch(groupUrl)
        .then(resp => resp.json())
        .then(data => {
            console.log("group (context): ", data);
            let [groupsArr] = Object.values(data); 
            setGroupsList(groupsArr);
        })
        .catch(
            err => console.log(err)
        );
    };
    
    useEffect(getGroupsHandler, []);

    return (
        <GroupsContext.Provider value={{
            groups: groupsList,
            // onlineGroups: onlineGroupsList,
            onNewGroupCreated: getGroupsHandler,
            // onOnline: groupOnlineHandler,
            // onOffline: groupOfflineHandler,
        }}>
        {props.children}
        </GroupsContext.Provider>
    );
};