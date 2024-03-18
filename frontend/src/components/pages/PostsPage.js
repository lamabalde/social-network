import React, { useParams, useEffect, useState, useCallback } from "react";
import FormLabel from "../UI/FormLabel";
import FormInput from "../UI/FormInput";
import FormTextarea from "../UI/FormTextarea";
// import styles from "./PostsPage.module.css";
import styles from "./layout.module.css";
import CreatePost from "../posts/CreatePost";
import AllPosts from "../posts/AllPosts";
import AllEvents from "../group/AllEvents";
import FollowRequest from "../requests/FollowRequest";
import Card from "../UI/Card";
import GroupRequest from "../requests/GroupRequests";
import useGet from "../fetch/useGet";
import { Navigate } from "react-router-dom";
import JoinedGroups from "../group/JoinedGroup";
import UserEvents from "../posts/UserEvents";
import AllUsers from "../posts/allUsers";
import refreshImg from "../assets/refresh.svg";

const PostsPage = () => {
	const sessionUrl = "http://localhost:8080/session";
	// const postUrl = "http://localhost:8080/post";
	const postCommentUrl = "http://localhost:8080/post-comment";
	const [postData, setPostData] = useState([]);
	// const [noPost, setNoPost] = useState(true);
	const [isLoadingPost, setIsLoadingPost] = useState(false);
	const [commentData, setCommentData] = useState([]);
	const [refreshState, setRefreshState] = useState(false);

	// get posts
	let userId = localStorage.getItem("user_id");
	console.log(userId);
	useEffect(() => {
		setIsLoadingPost(true);
		fetch(`http://localhost:8080/post?id=${userId}`)
			.then((resp) => resp.json())
			.then((data) => {
				console.log("all post data: ", data);
				// data &&
				if (data.length > 0) {
					data.sort(
						(a, b) => Date.parse(b.createdat) - Date.parse(a.createdat)
					);
					// setNoPost(false);
					setPostData(data);
					setIsLoadingPost(false);
				} else {
					// setNoPost(true);
					setIsLoadingPost(false);
				}
			})
			.catch((err) => {
				console.log(err);
				setIsLoadingPost(false);
			});
	}, [refreshState]);

	// get comments
	useEffect(() => {
		fetch(postCommentUrl)
			.then((resp) => resp.json())
			.then((data) => {
				data.sort((a, b) => Date.parse(a.createdat) - Date.parse(b.createdat)); // ascending order
				setCommentData(data);
			})
			.catch((err) => console.log(err));
	}, [refreshState]);

	// create post
	const createPostHandler = useCallback((createPostPayloadObj) => {
		console.log("postpage create post", createPostPayloadObj);
		const reqOptions = {
			method: "POST",
			body: JSON.stringify(createPostPayloadObj),
		};
		fetch(`http://localhost:8080/post`, reqOptions)
			.then((resp) => resp.json())
			.then((data) => {
				console.log("post success", data.success);
				if (data.success) {
					// render all posts
					fetch(`http://localhost:8080/post?id=${userId}`)
						.then((resp) => {
							return resp.json();
						})
						.then((data) => {
							// console.log("post data: ", data)
							// setPostData(data) // no need to sort
							// console.log("parsed 0", Date.parse(data[0].createdat));
							// console.log("parsed 1", Date.parse(data[1].createdat));
							data.sort(
								(a, b) => Date.parse(b.createdat) - Date.parse(a.createdat)
							);
							console.log("sorted post data: ", data);
							setPostData(data);
						})
						.catch((err) => console.log(err));
				}
			})
			.catch((err) => {
				console.log(err);
			});
	}, []);

	const createCommentSuccessHandler = useCallback((createCommentSuccessful) => {
		// fetch comment
		if (createCommentSuccessful) {
			fetch(postCommentUrl)
				.then((resp) => resp.json())
				.then((data) => {
					data.sort(
						(a, b) => Date.parse(a.createdat) - Date.parse(b.createdat)
					); // ascending order
					console.log("post page sorted comment data: ", data);
					setCommentData(data);
				})
				.catch((err) => console.log(err));
		}
	}, []);

	function refresh() {
		refreshState ? setRefreshState(false) : setRefreshState(true);
	}

	return (
		<div className={styles.container}>
			<div className={styles.mid}>
				<CreatePost onCreatePost={createPostHandler} />
				<div className={styles.refreshContainer}>
					<button className={styles.refresh} onClick={refresh}>
						<img src={refreshImg} alt=""></img>
					</button>
				</div>
				{/* {noPost && <AllPosts noPost={noPost} />} */}
				{/* {!noPost && ( */}
				<AllPosts
					posts={postData}
					comments={commentData}
					onCreateCommentSuccessful={createCommentSuccessHandler}
					isLoadingPost={isLoadingPost}
				/>
				{/* )} */}
			</div>

			<div className={styles.right}>
				<UserEvents userId={userId} refresh={refreshState}></UserEvents>
				<JoinedGroups refresh={refreshState}></JoinedGroups>
				<AllUsers></AllUsers>
			</div>
		</div>
	);
};

export default PostsPage;
