
export class chatMessage {
    constructor(content) {
        this.messageID = null;
        this.userID = null; // server uses current user to identify the message's author and chat which the user is writing to
        this.groupID = null
        this.userName = null;
        this.content = content;
        this.dateCreate = new Date()
        this.images = null; 
    }
}