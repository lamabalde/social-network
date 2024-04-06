export class reactionModel {
    constructor(type, ID, reaction) {
        this.messageType = type,
        this.messageID = ID,
        this.reaction = reaction;
    }
    postLike(postID) {
        return new reactionModel("post", postID, true); // TODO test with it onlt changing the values of the obj
    }
    /* commentLike(commentID) {
        this.messageType = "comment";
        this.postID = postID;
        this.reaction = true;
        return this;
    } */
}