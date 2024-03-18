import { useState } from "react";
import { Outlet } from "react-router-dom";
import TopNav from "../navigation/TopNav";
import ChatSidebar from "../navigation/ChatSidebar";
import { UsersContextProvider } from "../store/users-context";
import { FollowingContextProvider } from "../store/following-context";
import { WebSocketContextProvider } from "../store/websocket-context";
import { GroupsContextProvider } from "../store/groups-context";
import { JoinedGroupContextProvider } from "../store/joined-group-context";

const Root = () => {
    // console.log("Root");
    const [showChat, setShowChat] = useState(false);

    const chatIconClickHandler = () =>  {
        console.log("Root chat icon clicked");
        !showChat ? setShowChat(true) : setShowChat(false);
    };

    const chatbarClickHandler = () =>  {
        console.log("Root chatsidebar clicked");
        !showChat ? setShowChat(true) : setShowChat(false);
    };

    return (
        <>
        <UsersContextProvider>
            <WebSocketContextProvider>
                <FollowingContextProvider>
                    <GroupsContextProvider>
                        <JoinedGroupContextProvider>
                            <TopNav onClickChatIcon={chatIconClickHandler}/>
                            <ChatSidebar showChat={showChat} toggleChatSidebar={chatbarClickHandler}/>
                            <Outlet/>
                        </JoinedGroupContextProvider>
                    </GroupsContextProvider>
                </FollowingContextProvider>
            </WebSocketContextProvider>
        </UsersContextProvider>
        </>
    );
};

export default Root;