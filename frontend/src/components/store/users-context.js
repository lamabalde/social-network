import React, {useState, useEffect} from "react";

export const UsersContext = React.createContext({
    users: [],
    // onlineUsers: [],
    onNewUserReg: () => {},
    // onPrivacyChange: () => {},
    // publicUsers: [],
    // onOnline: (onlineUser) => {},
    // onOffline: (offlineUser) => {},
    // otherListedChatUsers: [],
    // setOtherListedChatUsers: () => {},
});

export const UsersContextProvider = (props) => {
    const [usersList, setUsersList] = useState([]);
    // const [publicUsersList, setPublicUsersList] = useState([]);
    // const [otherListedChatUsers, setOtherListedChatUsers] = useState([]);

    // get users
    const getUsersHandler = () => {
        const userUrl = "http://localhost:8080/user";
        fetch(userUrl)
        .then(resp => resp.json())
        .then(data => {
            console.log("user (context): ", data);
            let [usersArr] = Object.values(data);
            usersArr.sort((a,b) => a.id - b.id);
            setUsersList(usersArr);
        })
        .catch(
            err => console.log(err)
        );
    };

    // const getInitialUserPrivacy = () => {
    //     setPublicUsersList(usersList.filter(user => user.public === 1))
    // };

    // const privacyChangeHandler = (userid, privacy) => {
    //     // usersList[userid].public = privacy;
    //     let privacyChangedUser;
    //     setUsersList((prevUsersList) => prevUsersList.map((user) => {
    //         if (user.id === userid){
    //             privacyChangedUser = user;
    //             return user.public = privacy;
    //         } else {
    //             return user;
    //         }
    //     }));
    //     // if (publicUsersList) {
    //     //     // public
    //     //     privacy && setPublicUsersList(prevPublicList => [privacyChangedUser, ...new Set(...prevPublicList)]);
    //     //     !privacy && setPublicUsersList(prevPublicList => prevPublicList.filter(publicUser => publicUser.id !== privacyChangedUser.id));
    //     // } else { // first public user
    //     //     setPublicUsersList([privacyChangedUser]);
    //     // }
    // };



    useEffect(getUsersHandler, []);
    // useEffect(getInitialUserPrivacy, []);

    return (
        <UsersContext.Provider value={{
            users: usersList,
            // onlineUsers: onlineUsersList,
            onNewUserReg: getUsersHandler,
            // onPrivacyChange: privacyChangeHandler,
            // publicUsers: publicUsersList,
            // onOnline: userOnlineHandler,
            // onOffline: userOfflineHandler,
            // otherListedChatUsers: otherListedChatUsers,
            // setOtherListedChatUsers: setOtherListedChatUsers,
        }}>
        {props.children}
        </UsersContext.Provider>
    );
};