import { ref } from "vue";
import ws from "../ws";
import { formToObj } from "../helpers/helpers";
import { reactionModel } from "../models/reactionModel";
import { submitPostImg } from "../requests";

export const postsComponent = ref({
  posts: [],
  temporaryPostImg: "",
  //awaitPayload: any, // perhaps have it update this in the request and if the req is sucessful, then use it in response/rendering
  renderPostsPortion(payload) {
    if (payload.data != null) {
      payload.data.forEach((post) => {
        this.posts.push(post);
      });
    }
  },
  // new post
  async submitNewPost(event) {
    const form = await formToObj(event); 
    form["postType"] = parseInt(form["postType"]); // because form value doesn't support int, so need to convert
    const image = form["image"]
    if (image!== undefined) {
      this.temporaryPostImg = image
      delete form["image"];
      delete form["imageType"];
      // handle image submission with ajax
    }
    ws.request("newPost", form);
  },
  async renderNewPost(payload) {
    if (this.temporaryPostImg != "") { 
      const image = await submitPostImg(this.temporaryPostImg, payload.data.id); // submit the post image with the post ID, wait for response before rendering post
      const imageURL = `img/post/${image}`;
      payload.data.image = imageURL // add the image url to the ws response payload before rendering
    }
    this.temporaryPostImg=""
    if (!this.posts.length) {
      // to prevent errors
      this.posts.push(payload.data);
    } else {
      this.posts.unshift(payload.data);
    }
  },
  // delete post
  deletePostReq(postID) {
    ws.request("deletePost", postID); // Handle post deletion logic with the backend
  },
  deletePost(payload) {
    // If the request is successful, remove the post locally
    if (payload.result == "success") {
      let i = this.posts.map((post) => post.id).indexOf(payload.data);
      this.posts.splice(i, 1);
    } else {
      alert("Error deleting post");
    }
  },
  // request different types of posts
  getMainPosts(offset) {
    if (offset == 0) {
      this.posts = [];
    }
    ws.request("postsPortion", offset);
  },
  getUserPosts(userID, offset) {
    if (offset == 0) {
      this.posts = [];
    }
    const payload = {
      userID: userID,
      offset: offset,
    };
    ws.request("userPostsPortion", payload);
  },
  getGroupPosts(groupID, offset) {
    if (offset == 0) {
      this.posts = [];
    }
    const payload = {
      groupID: groupID,
      offset: offset,
    };
    ws.request("postsPortionInGroup", payload);
  },
  // likes
  likePost(postID) {
    const init = new reactionModel().postLike(postID)
    ws.request("likePost", init);
    this.awaitPostLikeID = postID
  },
  awaitPostLikeID: "",
  dislikePost(postID) {
    ws.request("dislikePost", postID);
  },
  renderPostLike(payload) {
    let i = this.posts.map((post) => post.id).indexOf(this.awaitPostLikeID);
    this.posts[i].content.likes[1] = payload.data.likes
  }
});
