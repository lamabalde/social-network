
import classes from './AllComments.module.css'
import Comment from './Comment';

function AllComments(props) {

    

    return (
        props.comments.map((comment) => {
            return <Comment
                key={comment.id}
                id={comment.id}
                // postId={comment.postid}
                authorId={comment.userid}
                fname={comment.fname}
                lname={comment.lname}
                avatar={comment.avatar}
                nname={comment.nname}
                message={comment.message}
                createdAt={comment.createdat} 
                image={comment.image}
            />
        })
    );
    
  

    
    
    ;
}

export default AllComments;