package tests

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/helpers"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/queries"
)

type postInDB struct {
	id               string
	theme            string
	content          string
	images           sql.NullString
	category         sql.NullString
	userID           string
	dateCreate       time.Time
	commentsQuantity int
	groupID          sql.NullString
}

func TestInsertAndDeletePost(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	f := queries.DBModel{DB: db}

	loc1, _ := time.LoadLocation("Europe/Kyiv")
	content := models.Content{
		UserID:     "1",
		Text:       "it's content1",
		Images:     []string{},
		DateCreate: time.Date(2023, time.March, 1, 11, 11, 11, 0, loc1).UTC(),
	}
	post := &models.Post{
		Content:    content,
		Theme:      "theme1",
		Categories: []string{"cat1", "cat2"},
	}

	var nullStr sql.NullString

	exp := postInDB{
		theme:            post.Theme,
		content:          content.Text,
		images:           nullStr,
		category:         sql.NullString{String: "cat1" + helpers.SEPARATOR + "cat2", Valid: true},
		userID:           content.UserID,
		dateCreate:       content.DateCreate,
		commentsQuantity: 0,
		groupID:          nullStr,
	}
	post, err = f.AddPost(post)
	if err != nil {
		t.Fatal(err)
	}
	exp.id = post.ID

	var res postInDB
	row := db.QueryRow("SELECT * FROM posts WHERE id=?", post.ID)
	err = row.Scan(&res.id, &res.theme, &res.content, &res.images, &res.category, &res.userID, &res.dateCreate, &res.commentsQuantity, &res.groupID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(exp, res) {
		t.Fatalf("expected \n%#v, \ngot \n%#v", exp, res)
	}
	fmt.Printf("---post=-------\n%s\n", post)
	fmt.Printf("---res=-------\n%v\n", res)

	f.DeletePost(post.ID)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("---delete err=%v-------\n", err)
	fmt.Printf("try to get deleted post (id: %s)\n", post.ID)
	p, err := f.GetPostByID(post.ID, "1")

	if err != nil {
		if errors.Is(err, models.ErrNoRecords) {
			fmt.Println("no post found")
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatalf("the post was not deleted successfully. Got the post from DB %v", p)
	}

	fmt.Println("-----add post and 2 comments to it-----------")
	post, err = f.AddPost(post)
	if err != nil {
		t.Fatal(err)
	}
	comm1, err := f.AddComment(post.ID, "comm1 from 2", []string{}, "2")
	if err != nil {
		t.Fatal(err)
	}

	comm2, err := f.AddComment(post.ID, "comm2 from 1", []string{}, "1")
	if err != nil {
		t.Fatal(err)
	}

	exp.id = post.ID
	exp.commentsQuantity = 2
	row = db.QueryRow("SELECT * FROM posts WHERE id=?", post.ID)
	err = row.Scan(&res.id, &res.theme, &res.content, &res.images, &res.category, &res.userID, &res.dateCreate, &res.commentsQuantity, &res.groupID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(exp, res) {
		t.Fatalf("\nexpected %v\n, \ngot %v\n", exp, res)
	}

	post, err = f.GetPostByID(post.ID, "1")
	fmt.Printf("---post=-------\n%s", post)

	f.DeletePost(post.ID)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("---delete err=%v-------\n", err)
	fmt.Printf("try to get deleted post (id: %s)\n", post.ID)
	p, err = f.GetPostByID(post.ID, "1")
	if err != nil {
		if errors.Is(err, models.ErrNoRecords) {
			fmt.Println("no post found")
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatalf("the post was not deleted successfully. Got the post from DB %v", p)
	}

	// check if comments are deleted
	row = db.QueryRow("SELECT * FROM comments WHERE id=?", comm1.ID)

	err = row.Scan(&res.id, &res.content, &res.images, &res.userID, &res.groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("no comments found")
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatalf("the comment was not deleted successfully. Got the post from DB %v", res)
	}

	row = db.QueryRow("SELECT * FROM comments WHERE id=?", comm2.ID)
	err = row.Scan(&res.id, &res.content, &res.images, &res.userID, &res.groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("no comments found")
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatalf("the comment was not deleted successfully. Got the post from DB %v", res)
	}

	// add post to group
	exp.commentsQuantity = 0
	exp.groupID = sql.NullString{String: "1", Valid: true}
	post.GroupID = "1"

	post, err = f.AddPost(post)
	if err != nil {
		t.Fatal(err)
	}
	exp.id = post.ID

	row = db.QueryRow("SELECT * FROM posts WHERE id=?", post.ID)
	err = row.Scan(&res.id, &res.theme, &res.content, &res.images, &res.category, &res.userID, &res.dateCreate, &res.commentsQuantity, &res.groupID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(exp, res) {
		t.Fatalf("expected %v, got %v", exp, res)
	}
	fmt.Printf("---post=-------\n%s\n", post)
	fmt.Printf("---res=-------\n%v\n", res)
}

func TestGetPostsByCondition(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	fmt.Println("------------check reaction of user id 3 --------------------")
	posts, err := f.GetPostsByCondition("", nil, "3", 30, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.String())
	}
	for _, post := range posts {
		if post.Content.UserReaction != -1 {
			t.Fatalf("for the user with id=3: it is expected to be no reaction, ther is a reaction to the post with id=%s", post.ID)
		}
	}

	fmt.Println("------------check reaction of user id 2 --------------------")
	posts, err = f.GetPostsByCondition("", nil, "2", 20, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.String())
	}
	for _, post := range posts {
		if !(post.ID == "2" || post.ID == "3" || post.ID == "4") && post.Content.UserReaction != -1 {
			t.Fatalf("for the user with id=2: it is expected to be no reaction for the all posts except for id=2, 3 or 4 , ther is a reaction to the post with id=%s", post.ID)
		}
		if post.ID == "2" && post.Content.UserReaction != 1 {
			t.Fatalf("for the user with id=2: it is expected like (reaction = 1 ) for the posts  id=2, ther is a reaction %d to the post", post.Content.UserReaction)
		}
		if (post.ID == "3" || post.ID == "4") && post.Content.UserReaction != 0 {
			t.Fatalf("for the user with id=2: it is expected dislike (reaction = 0 ) for the posts  id=3 or 4, ther is a reaction %d to the post id %s", post.Content.UserReaction, post.ID)
		}
	}
}

func BenchmarkGetPostsByCondition(b *testing.B) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	posts, err := f.GetPostsByCondition("", nil, "3", 20, 0)
	if err != nil {
		b.Fatal(err)
	}

	for _, post := range posts {
		if post.Content.UserReaction != -1 {
			b.Fatalf("for the user with id=3: it is expected to be no reaction, ther is a reaction to the post with id=%s", post.ID)
		}
	}

	posts, err = f.GetPostsByCondition("", nil, "2", 20, 0)
	if err != nil {
		b.Fatal(err)
	}

	for _, post := range posts {
		if !(post.ID == "2" || post.ID == "3" || post.ID == "4") && post.Content.UserReaction != -1 {
			b.Fatalf("for the user with id=2: it is expected to be no reaction for the all posts except for id=2, 3 or 4 , ther is a reaction to the post with id=%s", post.ID)
		}
		if post.ID == "2" && post.Content.UserReaction != 0 {
			b.Fatalf("for the user with id=2: it is expected like (reaction = 0 ) for the posts  id=2, ther is a reaction %d to the post", post.Content.UserReaction)
		}
		if (post.ID == "3" || post.ID == "4") && post.Content.UserReaction != 1 {
			b.Fatalf("for the user with id=2: it is expected dislike (reaction = 1 ) for the posts  id=3 or 4, ther is a reaction %d to the post id %s", post.Content.UserReaction, post.ID)
		}
	}
}

func TestGetPostsFilters(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	fmt.Println("--get posts--")

	filter := models.Filter{
		AuthorID:      "",
		GroupID:       "",
		LikedByUserID: "",
	}
	posts, err := f.GetFiltredPosts(filter, "0", 20, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.String())
	}

	fmt.Println("--get posts author 2 (2,4,8)--")
	filter = models.Filter{
		AuthorID:      "2",
		GroupID:       "",
		LikedByUserID: "",
	}
	posts, err = f.GetFiltredPosts(filter, "0", 20, 0)
	if err != nil {
		t.Fatal(err)
	}
	expRes := []string{"8", "4", "2"}
	for i, post := range posts {
		fmt.Printf("%s\n", post.String())
		if post.ID != expRes[i] {
			t.Fatalf("expected post ID %s, got %s", expRes[i], post.ID)
		}
	}

	fmt.Println("--get posts liked by user 1--")
	filter = models.Filter{
		AuthorID:      "",
		GroupID:       "",
		LikedByUserID: "1",
	}
	posts, err = f.GetFiltredPosts(filter, "0", 20, 0)
	if err != nil {
		t.Fatal(err)
	}

	expRes = []string{"3", "2", "1"}
	for i, post := range posts {
		fmt.Printf("%s\n", post.String())
		if post.ID != expRes[i] {
			t.Fatalf("expected post ID %s, got %s", expRes[i], post.ID)
		}
	}

	fmt.Println("--get posts by author 1--")
	filter = models.Filter{
		AuthorID:      "1",
		GroupID:       "",
		LikedByUserID: "",
	}
	posts, err = f.GetFiltredPosts(filter, "0", 20, 0)
	if err != nil {
		t.Fatal(err)
	}

	expRes = []string{"6", "5", "1"}
	for i, post := range posts {
		fmt.Printf("%s\n", post.String())
		if post.ID != expRes[i] {
			t.Fatalf("expected post ID %s, got %s", expRes[i], post.ID)
		}
	}
}

func TestGetPostsNumberPosts(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	filter := models.Filter{
		AuthorID:      "",
		GroupID:       "",
		LikedByUserID: "",
	}

	fmt.Println("--get posts offset 10 5posts--")
	expRes := 0
	posts, err := f.GetFiltredPosts(filter, "0", 5, 10)
	if err != nil {
		t.Fatal(err)
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.String())
	}
	if len(posts) != 0 {
		t.Fatalf("should be %d posts, got %d", expRes, len(posts))
	}

	fmt.Println("--get posts offset 0 10posts expect 8--")
	expRes = 8
	posts, err = f.GetFiltredPosts(filter, "0", 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.String())
	}
	if len(posts) != expRes {
		t.Fatalf("should be %d posts, got %d", expRes, len(posts))
	}

	fmt.Println("--get posts offset=6 5posts: must be 2 posts--")
	expRes = 2
	posts, err = f.GetFiltredPosts(filter, "0", 5, 6)
	if err != nil {
		t.Fatal(err)
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.String())
	}
	if len(posts) != expRes {
		t.Fatalf("should be %d posts, got %d", expRes, len(posts))
	}

	fmt.Println("--get posts offset 0 (from the last one) 5 posts--")
	expRes = 5
	posts, err = f.GetFiltredPosts(filter, "0", 5, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.String())
	}
	if len(posts) != expRes {
		t.Fatalf("should be %d posts, got %d", expRes, len(posts))
	}
}

func TestGetPostsLikedByUser(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	fmt.Println("--get liked by user 1 --")
	expRes := []string{"3", "2", "1"}
	posts, err := f.GetPostsLikedByUser("1", 20, 0)
	if err != nil {
		t.Fatal(err)
	}
	for i, post := range posts {
		fmt.Printf("%s\n", post.String())
		if post.ID != expRes[i] {
			t.Fatalf("expected post ID %s, got %s", expRes[i], post.ID)
		}
	}

	fmt.Println("--get posts liked by user 2 --")
	expRes = []string{"2"}
	posts, err = f.GetPostsLikedByUser("2", 20, 0)
	if err != nil {
		t.Fatal(err)
	}
	for i, post := range posts {
		fmt.Printf("%s\n", post.String())
		if post.ID != expRes[i] {
			t.Fatalf("expected post ID %s, got %s", expRes[i], post.ID)
		}
	}

	fmt.Println("--get liked by user 0 --")
	expRes = []string{}
	posts, err = f.GetPostsLikedByUser("0", 20, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.String())
	}
	if len(posts) > 0 {
		t.Fatalf("expected post %d posts, got %d", len(expRes), len(posts))
	}
}

func TestGetPostByDI(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	fmt.Println("--get post 1--")
	expRes := "1"
	post, err := f.GetPostByID("1", "0")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", post.String())
	if post.ID != expRes {
		t.Fatalf("expected post ID %s, got %s", expRes, post.ID)
	}

	fmt.Println("--get post 3--")
	expRes = "3"
	post, err = f.GetPostByID("3", "0")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", post.String())
	if post.ID != expRes {
		t.Fatalf("expected post ID %s, got %s", expRes, post.ID)
	}
}

func TestGetPosts(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	testCases := []struct {
		groupID, currUserID,authorID string
		exp             []*models.Post
	}{
		{
			groupID: "",
			authorID: "",
			currUserID:  "1",
			exp: []*models.Post{
				{ID: "8", Theme: "Post with comments", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 0},
				{ID: "7", Theme: "Wise Kaa", Content: models.Content{UserID: "3"}, GroupID: "", Privacy: 0},
				{ID: "6", Theme: "Seamus", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 0},
				{ID: "5", Theme: "My parrot", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 0},
				{ID: "4", Theme: "My dog", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 0},
				{ID: "3", Theme: "My cat", Content: models.Content{UserID: "3"}, GroupID: "", Privacy: 0},
				{ID: "fr2", Theme: "postForFriend1#2", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 2},
				{ID: "fr3", Theme: "postForFriend2#1", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 2},
				{ID: "fr4", Theme: "postForFriend2#2", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 2},
				{ID: "pr1", Theme: "postPrivate2#1", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 1},
				{ID: "pr2", Theme: "postPrivate1#1", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 1},
				{ID: "pr3", Theme: "postPrivate3#1", Content: models.Content{UserID: "3"}, GroupID: "", Privacy: 1},
				{ID: "2", Theme: "dogs", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 0},
				{ID: "fr1", Theme: "postForFriend1#1", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 2},
				{ID: "1", Theme: "cats", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 0},
			},
		},
		{
			groupID: "",
			authorID: "",
			currUserID:  "2",
			exp: []*models.Post{
				{ID: "8", Theme: "Post with comments", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 0},
				{ID: "7", Theme: "Wise Kaa", Content: models.Content{UserID: "3"}, GroupID: "", Privacy: 0},
				{ID: "6", Theme: "Seamus", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 0},
				{ID: "5", Theme: "My parrot", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 0},
				{ID: "4", Theme: "My dog", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 0},
				{ID: "3", Theme: "My cat", Content: models.Content{UserID: "3"}, GroupID: "", Privacy: 0},
				{ID: "fr2", Theme: "postForFriend1#2", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 2},
				{ID: "fr3", Theme: "postForFriend2#1", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 2},
				{ID: "fr4", Theme: "postForFriend2#2", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 2},
				{ID: "pr1", Theme: "postPrivate2#1", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 1},
				{ID: "pr2", Theme: "postPrivate1#1", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 1},
				{ID: "pr3", Theme: "postPrivate3#1", Content: models.Content{UserID: "3"}, GroupID: "", Privacy: 1},
				{ID: "2", Theme: "dogs", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 0},
				{ID: "fr1", Theme: "postForFriend1#1", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 2},
				{ID: "1", Theme: "cats", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 0},
			},
		},
		{
			groupID: "",
			authorID: "",
			currUserID:  "3",
			exp: []*models.Post{
				{ID: "8", Theme: "Post with comments", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 0},
				{ID: "7", Theme: "Wise Kaa", Content: models.Content{UserID: "3"}, GroupID: "", Privacy: 0},
				{ID: "6", Theme: "Seamus", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 0},
				{ID: "5", Theme: "My parrot", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 0},
				{ID: "4", Theme: "My dog", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 0},
				{ID: "3", Theme: "My cat", Content: models.Content{UserID: "3"}, GroupID: "", Privacy: 0},
				{ID: "fr2", Theme: "postForFriend1#2", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 2},
				{ID: "fr5", Theme: "postForFriend3#1", Content: models.Content{UserID: "3"}, GroupID: "", Privacy: 2},
				{ID: "pr1", Theme: "postPrivate2#1", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 1},
				{ID: "pr2", Theme: "postPrivate1#1", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 1},
				{ID: "pr3", Theme: "postPrivate3#1", Content: models.Content{UserID: "3"}, GroupID: "", Privacy: 1},
				{ID: "2", Theme: "dogs", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 0},
				{ID: "fr1", Theme: "postForFriend1#1", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 2},
				{ID: "1", Theme: "cats", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 0},
			},
		},
		{
			groupID: "",
			authorID: "",
			currUserID:  "4",
			exp: []*models.Post{
				{ID: "8", Theme: "Post with comments", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 0},
				{ID: "7", Theme: "Wise Kaa", Content: models.Content{UserID: "3"}, GroupID: "", Privacy: 0},
				{ID: "6", Theme: "Seamus", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 0},
				{ID: "5", Theme: "My parrot", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 0},
				{ID: "4", Theme: "My dog", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 0},
				{ID: "3", Theme: "My cat", Content: models.Content{UserID: "3"}, GroupID: "", Privacy: 0},
				{ID: "fr2", Theme: "postForFriend1#2", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 2},
				{ID: "2", Theme: "dogs", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 0},
				{ID: "fr1", Theme: "postForFriend1#1", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 2},
				{ID: "1", Theme: "cats", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 0},
			},
		},
		{
			groupID: "1",
			authorID: "",
			currUserID:  "2",
			exp: []*models.Post{
				{ID: "gr2", Theme: "postingr1#2", Content: models.Content{UserID: "3"}, GroupID: "1", Privacy: 0},
				{ID: "gr1", Theme: "postingr1#1", Content: models.Content{UserID: "2"}, GroupID: "1", Privacy: 0},
			},
		},
		{
			groupID: "1",
			authorID: "",
			currUserID:  "1",
			exp: []*models.Post{
				{ID: "gr2", Theme: "postingr1#2", Content: models.Content{UserID: "3"}, GroupID: "1", Privacy: 0},
				{ID: "gr1", Theme: "postingr1#1", Content: models.Content{UserID: "2"}, GroupID: "1", Privacy: 0},
			},
		},
		{
			groupID: "",
			authorID: "1",
			currUserID:  "1",
			exp: []*models.Post{
				{ID: "6", Theme: "Seamus", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 0},
				{ID: "5", Theme: "My parrot", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 0},
				{ID: "fr2", Theme: "postForFriend1#2", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 2},
				{ID: "pr2", Theme: "postPrivate1#1", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 1},
				{ID: "fr1", Theme: "postForFriend1#1", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 2},
				{ID: "1", Theme: "cats", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 0},
			},
		},
		{
			groupID: "",
			authorID: "2",
			currUserID:  "2",
			exp: []*models.Post{
				{ID: "8", Theme: "Post with comments", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 0},
				{ID: "4", Theme: "My dog", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 0},
				{ID: "fr3", Theme: "postForFriend2#1", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 2},
				{ID: "fr4", Theme: "postForFriend2#2", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 2},
				{ID: "pr1", Theme: "postPrivate2#1", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 1},
				{ID: "2", Theme: "dogs", Content: models.Content{UserID: "2"}, GroupID: "", Privacy: 0},
				{ID: "gr1", Theme: "postingr1#1", Content: models.Content{UserID: "2"}, GroupID: "1", Privacy: 0},
			},
		},
		{
			groupID: "",
			authorID: "1",
			currUserID:  "3",
			exp: []*models.Post{
				{ID: "6", Theme: "Seamus", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 0},
				{ID: "5", Theme: "My parrot", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 0},
				{ID: "fr2", Theme: "postForFriend1#2", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 2},
				{ID: "pr2", Theme: "postPrivate1#1", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 1},
				{ID: "fr1", Theme: "postForFriend1#1", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 2},
				{ID: "1", Theme: "cats", Content: models.Content{UserID: "1"}, GroupID: "", Privacy: 0},
			},
		},
		{
			groupID: "",
			authorID: "3",
			currUserID:  "4",
			exp: []*models.Post{
				{ID: "7", Theme: "Wise Kaa", Content: models.Content{UserID: "3"}, GroupID: "", Privacy: 0},
				{ID: "3", Theme: "My cat", Content: models.Content{UserID: "3"}, GroupID: "", Privacy: 0},
				{ID: "gr2", Theme: "postingr1#2", Content: models.Content{UserID: "3"}, GroupID: "1", Privacy: 0},
			},
		},
	}

	for i, test := range testCases {
		var posts []*models.Post
		if test.groupID == "" && test.authorID == ""  {
			posts, err = f.GetPostsNoGroup(test.currUserID, 20, 0)
			if err != nil {
				t.Fatal(err)
			}
		} 
		if test.groupID != "" && test.authorID == ""  {
			posts, err = f.GetPostsInGroup(test.groupID, test.currUserID, 20, 0)
			if err != nil {
				t.Fatal(err)
			}
		}
		if test.groupID == "" && test.authorID != ""  {
			posts, err = f.GetUserPosts(test.authorID, test.currUserID, 20, 0)
			if err != nil {
				t.Fatal(err)
			}
		}

		if !comparePosts(posts, test.exp) {
			t.Fatalf("Test# %d (userID %s, groupID %s): \nExpected-----\n %s, \n\ngot----------\n %s", i, test.currUserID, test.groupID, stringPosts(test.exp), stringPosts(posts))
		}

	}
}

func stringPost(p *models.Post) string {
	return fmt.Sprintf("ID: %s\nTheme: %s\n AuthorID: %s, Group %s, Privacy: %d", p.ID, p.Theme, p.Content.UserID, p.GroupID, p.Privacy)
}

func stringPosts(posts []*models.Post) string {
	str := ""
	for _, post := range posts {
		str += fmt.Sprintf("%s\n", stringPost(post))
	}
	return str
}

func comparePost(a, b *models.Post) bool {
	return a.ID == b.ID && a.Theme == b.Theme && a.Content.UserID == b.Content.UserID && a.GroupID == b.GroupID && a.Privacy == b.Privacy
}

func comparePosts(a, b []*models.Post) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if !comparePost(a[i], b[i]) {
			return false
		}
	}
	return true
}
