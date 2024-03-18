import classes from './AllGroupPosts.module.css'
import GroupPost from "./GroupPost";


function AllGroupPosts(props) {
    // data.data && data.data.sort((a, b) => Date.parse(b.createdat) - Date.parse(a.createdat));
   console.log(props.post, "postprrops")
    return <div className={classes.container}>
        {props.posts && props.posts.map((post) => (
            <GroupPost
                key={post.id}
                id={post.id}
                avatar={post.image}
                author={post.author}
                fname={post.fname}
                lname={post.lname}
                nname={post.nname}
                message={post.message}
                image={post.image}
                createdat={post.createdat}
                authorId={post.author}
            />
        ))}
    </div>
}

export default AllGroupPosts;
