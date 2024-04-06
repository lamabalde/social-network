import { computed, ref, toRaw } from "vue";
import ws from "../ws";
import { formToObj } from "../helpers/helpers";
import { submitCommentImg } from "../requests";

export const commentsComponent = ref({
  comments: {}, // key: postID, value: array of comments
  temporaryCommentImg: "",
  getAllComments(postID) {
    ws.request("fullPostAndComments", postID);
  },
  renderAllComments(payload) {
    if (payload.data.comments != undefined) {
      payload.data.comments.sort(function (a, b) {
        return a.content.dateCreate > b.content.dateCreate;
      });
    }
    
    this.comments[payload.data.id]= payload.data.comments;
  },
  async submitComment(event) {
    const form = await formToObj(event);
    const image = form["image"];
    if (image !== undefined) {
      this.temporaryCommentImg = image;
      delete form["image"];
      delete form["imageType"];
      // handle image submission with ajax
    }
    ws.request("newComment", form);
  },
  async renderNewComment(payload) {
    if (this.temporaryCommentImg != "") { 
      const image = await submitCommentImg(this.temporaryCommentImg, payload.data.comments[0].id); // submit the post image with the post ID, wait for response before rendering post
      const imageURL = `img/comment/${image}`;
      payload.data.comments[0].content.image = imageURL; // add the image url to the ws response payload before rendering
    }
    this.temporaryCommentImg = ""

    if (Object.values(this.comments).includes(payload.data.id)) { // if comments already exist for this post, add more, else create new 
      this.comments[payload.data.id].push(payload.data.comments);
    } else {
      payload.data.comments.sort(function (a, b) {
        return a.content.dateCreate > b.content.dateCreate;
      });
      this.comments[payload.data.id] = payload.data.comments;
    }
  },
  amountOfComments(postID) {
    const commentsAmount = toRaw(commentsComponent.value.comments[postID]);
    return commentsAmount == undefined ? 0 : commentsAmount.length;
  }
});
