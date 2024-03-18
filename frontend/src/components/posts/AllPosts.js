import Post from "./Post";
import classes from "./AllPosts.module.css";
import useGet from "../fetch/useGet";
import { useCallback, useEffect } from "react";

const AllPosts = ({
	posts,
	comments,
	onCreateCommentSuccessful,
	isLoadingPost,
}) => {
	//userId
	// const { data } = useGet(`/posts`)
	// console.log("out", comments);
	let eachPostCommentsArr = [];

	let noPostMsg = <h3>Loading...</h3>;
	if (!isLoadingPost) noPostMsg = <h3>No Posts Yet...</h3>;
	// useEffect(() => setPosts(posts), [posts]);

	if (posts)
		for (let i = 0; i < posts.length; i++) {
			let thisPostComments = [];
			for (let j = 0; j < comments.length; j++) {
				comments[j] &&
					comments[j].postid === posts[i].id &&
					thisPostComments.push(comments[j]);
			}
			eachPostCommentsArr.push(thisPostComments);
		}
	// console.log("eachPostComments", eachPostCommentsArr);

	const createCommentSuccessHandler = useCallback(
		(createCommentSuccessful) => {
			// lift it up to PostPage
			onCreateCommentSuccessful(createCommentSuccessful);
		},
		[onCreateCommentSuccessful]
	);

	return (
		<>
			{!posts ? (
				noPostMsg
			) : (
				<div className={classes.container}>
					{posts.map((post, p) => (
						<Post
							key={post.id}
							id={post.id}
							avatar={post.avatar}
							fname={post.fname}
							lname={post.lname}
							nname={post.nname}
							message={post.message}
							image={post.image}
							createdat={post.createdat}
							authorId={post.author}
							privacy={post.privacy}
							// totalNumPost={posts.length}
							postNum={p}
							commentsForThisPost={eachPostCommentsArr[p]}
							onCreateCommentSuccessful={createCommentSuccessHandler}
						/>
					))}
				</div>
			)}
		</>
	);
};

export default AllPosts;
