WS messages on backend
**package webmodel** (backend/pkg/webmodel)
declares 
- constants for WS messages type 
- structs which  WS *messages* requests are parsed to, and some of structs for WS replies. 
A *message* is a JSON object {"type": \<string\>, payload: \<data\>} 
payload could be:
	- string (for IDs)
	- number (for offset)
	- object

**package controllers** (backend/pkg/controllers)
contains functions to create data to  reply to WS request. Those funcs named as Reply+nameOfMessage, the other functions there are helpers and they are not exported. (ex: func ReplyNewPost() creates reply data for the request with type "newPost", func savePostToDB() is a helper for ReplyNewPost())
"Repliers" run function to work with DB and create data for WS reply
- if data returned by DB func is enough for the reply, return it (DB funcs return either primitive  or an object from package models)
- if not, creates a reply using struct from webmodels.

These "repliers" must be tied to a message type in **package wsrouter** (backend/pkg/routes/wsrouter/wsrouter.go), func CreateAPIwsrouters(),
like 
`webmodel.MessageType:   sendReplyForLoggedUser(wsconnection.FuncReplyCreator(controllers.ReplyMessageType)),` 
 where
- webmodel.MessageType - a const from **package webmodel**
- sendReplyForLoggedUser - midlware func to check if the user is logged in
- FuncReplyCreator - function, which creates and sends a message with 
	- the according type 
	- and payload like {"result":"success", "data":\<data from replier\>} (any errors payloads, with "result":"error" are created in repliers if an error occurs)

func ReadPump() from package wsconnection reads WS message and runs the "replier" "tied" to the type of the message.
So, creating a message to reply to WS request requires:  
- in **package queries** (backend/pkg/db/sqlite/queries)
	- add a func that will exec a DB query, if there is no appropriate one 
-  in **package webmodel**  (backend/pkg/webmodel)
	- define the constants for WS messages type
	- if needed define a structs which  WS *messages* requests are parsed to as well as Valide() method for the struct and func PayloadTo.... in **package parse** (backend/pkg/webmodel/parse)
- in **package controllers** (backend/pkg/controllers)
	- create a func Reply...
-  in **package wsrouter** (backend/pkg/routes/wsrouter)
	- add an instruction like 
	`webmodel.MessageType: sendReplyForLoggedUser(wsconnection.FuncReplyCreator(controllers.ReplyMessageType))` 
	to func CreateAPIwsrouters()


WS messages
if there is a **server error**, the reply will be like:
{
	type: "ERROR", 
	data:{
		result: "serverError", 
		data: \<errmessage\>
	}
}

if there is a **request error**, the reply will be like:
{
	type: \<typeOfRequest\>, 
	data:{
		result: "error", 
		data: \<errmessage\>
	}
}

**"successful"** replies (payload.result="success"):

*Type*: **"postsPortion"**
	*request payload*: offset: number
	*reply payload.data:* array of posts without comments
		[] {
		    id               string     
		    theme            string  
		    content  {
			    userID           string    
			    userName         string       
			    text             string       
			    dateCreate       time.Time
			    likes            []int       // index 0 keeps number of dislikes, index 1 keeps number of like
			    userReaction     int //-1 => no reaction, 0=>current user disliked, 1=>current user liked
			}    
		    categories          [] string   
		    commentsQuantity    int 
			groupID 			string		    
		}

*Type*: **""userPostsPortion""**
	*request payload*: webmodel.UserPosts  
		{
			UserID string `json:"userID"`
			Offset int    `json:"offset"`
		}
	*reply payload.data:* array of posts without comments
		[] {
		    id               string     
		    theme            string  
		    content  {
			    userID           string    
			    userName         string       
			    text             string       
			    dateCreate       time.Time
			    likes            []int       // index 0 keeps number of dislikes, index 1 keeps number of like
			    userReaction     int //-1 => no reaction, 0=>current user disliked, 1=>current user liked
			}    
		    categories          [] string   
		    commentsQuantity    int 
			groupID 			string		    
		}

*Type*: **"fullPostAndComments"**
	*request payload*:  postId: string
	*reply payload.data:* post with comments
	{
		id        string     
		theme     string  
		content  {
			userID        string    
			userName      string       
			text          string       
			dateCreate    time.Time
			likes         []int       // index 0 keeps number of dislikes, index 1 keeps number of like
			userReaction  int //-1 => no reaction, 0=>current user disliked, 1=>current user liked
		}    
		categories [] string   
		comments   []{
			id      string  
			postID  string  
			content  {
				userID         string    
				userName       string       
				text           string       
				dateCreate     time.Time
				likes          []int       // index 0 keeps number of dislikes, index 1 keeps number of like
				userReaction   int  //-1 => no reaction, 0=>current user disliked, 1=>current user liked
			}
		}
		commentsQuantity    int  
		groupID 			string   
	}

*Type*: **"newPost"**
	*request payload*: webmodel.Post
		{
			title      string   
		    body       string   
		    categories []string
		}
	*reply payload.data:* post without comments
		{
		    id               string     
		    theme            string  
		    content  {
			    userID           string    
			    userName         string       
			    text             string       
			    dateCreate       time.Time
			    likes            []int       // index 0 keeps number of dislikes, index 1 keeps number of like
			    userReaction     int //-1 => no reaction, 0=>current user disliked, 1=>current user liked
			}    
		    categories          [] string   
		    commentsQuantity    int 
			groupID 			string		    
		}

*Type*:  **"newComment"**
	*request payload*:  webmodel.Comment
		{
			postID     string 
			content    string 
		}
	*reply payload.data:* post with comments
		{
			id                  string     
			theme               string  
			content  {
				userID       string    
				userName     string       
				text         string       
				dateCreate   time.Time
				likes        []int       // index 0 keeps number of dislikes, index 1 keeps number of like
				userReaction int //-1 => no reaction, 0=>current user disliked, 1=>current user liked
			}    
			categories [] string   
			comments   []{
				id      string  
				postID  string  
				content  {
					userID       string    
					userName     string       
					text         string       
					dateCreate   time.Time
					likes        []int       // index 0 keeps number of dislikes, index 1 keeps number of like
					userReaction int //-1 => no reaction, 0=>current user disliked, 1=>current user liked
				}
			}
			commentsQuantity    int 
			groupID 			string	    
		}	

*Type*:  **"onlineUsers"**
	*request payload*:  - 
	*reply payload.data:* 
		{
			id              string
			userName        string
			lastMessageDate string
		}

*Type*:  **"openChat"**
	*request payload*:  recepientID: string
	*reply payload.data:* privateChat
		{
			id                 string               
			name               string               
			recipientUserID    string               
			recipientUserName  string               
			messages       []{
				id          string   
				userID      string   
				userName    string   
				content     string   
				dateCreate  time.Time
			}
		}

*Type*:  **"sendMessageToOpendChat"**
	*request payload*:  webmodel.ChatMessage
		{
			content       string   
			dateCreate    time.Time
		}
	*reply payload.data:* 
		*to creator payload.data*: "delivered"
		*to recepient *webmodel.ChatMessage
	
			*Type*:  **"inputChatMessage"**
	 			payload.data
					{
						id         string   
						userID      string   //author
						userName    string   
						content     string   
						dateCreate  time.Time
					}


*Type*:  **"closeChat"**
	*request payload*:  null
	*reply payload.data:*  - // server will delete info about "current opened chat" 

*Type*:  **"chatPortion"**  
	*request payload*: offset: number
	*reply payload.data:* privateChat // may by change to - send just array of messages
		{
			id                  string               
			name                string               
			recipientUserID     string               
			recipientUserName 	string               
			messages                []{ 
				id              string   
				userID          string   
				userName        string   
				content         string   
				dateCreate      time.Time
			}
		}
*Type*:  **"newOnlineUser"**
	*request payload*: -
	*reply payload.data:*  *to other users*: user
		{
		   id                     string  
		   userName        string
		}

*Type*:  **"offlineUser"** 
	*request payload*: -
	*reply payload.data:*  *to other users*: user
		{
		   id                     string  
		   userName        string
		}

*Type*:  **"createGroup"**
	*request payload*:  webmodel.Group
		{
			title       string
			description string
		}
	*reply payload.data:* group
		{
			id          string    
			name        string    
			description string    
			creatorID   string    
			dateCreate  time.Time 
			members     []{
				id       string 
				userName string 
			}
		}

*Type*:  **"addUserToGroup"**
	*request payload*:  webmodel.UserInGroup
		{
			groupID string
			userID  string
			userName string
		}
	*reply payload.data:* webmodel.UserInGroup
		{
			groupID string
			userID  string
			userName string
		}

*Type*:  **"groupMembers"**
	*request payload*: groupID  string
	*reply payload.data:* user
		{
		   id              string  
		   userName        string
		}

*Type*:  **"userGroups"**
	*request payload*: userID  string
	*reply payload.data:* []group
		{
			id          string    
			name        string    
		}


*Type*:  **"getUserProfile"**
	*request payload*: userID  string
	*reply payload.data:* user
		{
			id              string   
			userName      	string   
			email           string   
			dateCreate      time.Time
			dateBirth       time.Time
			gender          string   
			firstName       string   
			lastName        string 
			profileType     int // 0 - public, 1 - private  
		}

*Type*: **"deletePost"**
	*request payload*:  postId: string
	*reply payload.data:* postId 

*Type*:  **"postsPortionInGroup"**
	*request payload*: webmodel.PostsInGroupPortion
		{
			groupID string 	
			offset  int    
		}
	*reply payload.data:* usarray of posts without commentser
		[] {
	    	id               string     
	    	theme            string  
	    	content  {
			    userID           string    
			    userName         string       
			    text             string       
			    dateCreate       time.Time
			    likes            []int       // index 0 keeps number of dislikes, index 1 keeps number of like
			    userReaction     int //-1 => no reaction, 0=>current user disliked, 1=>current user liked
			}    
	    	categories          [] string   
	    	commentsQuantity    int 
			groupID 			string		    
		}

*Type*:  **"userFollowers"**
	*request payload*: userID  string
	*reply payload.data:* user
		{
		  id              string  
		   userName       string
		}

*Type*:  **"userFollowing"**
	*request payload*: userID  string
	*reply payload.data:* user
		{
		  id              string  
		   userName       string
		}

*Type*:  **"followUser"**
	*request payload*: userID  string
	*reply payload.data:* followStatus // "requested" or "following"

*Type*:  **"getFollowStatus"**
	*request payload*: userID  string
	*reply payload.data:* followStatus // "requested" or "following"

*Type*:  **"searchGroupsUsers"**
	*request payload*: searchQuery  string
	*reply payload.data: 
		{
			entityType string  // groups or profile
			id         string 
			name       string 
		}
	
*Type*:  **"setProfileType"**
	*request payload*: profileType  string
	*reply payload.data:* profileType // "requested" or "following"


	NewNotification = "newNotification"


	SetProfileType = "setProfileType"
WS messages on backend
**package webmodel** (backend/pkg/webmodel)
declares 
- constants for WS messages type 
- structs which  WS *messages* requests are parsed to, and some of structs for WS replies. 
A *message* is a JSON object {"type": \<string\>, payload: \<data\>} 
payload could be:
	- string (for IDs)
	- number (for offset)
	- object

**package controllers** (backend/pkg/controllers)
contains functions to create data to  reply to WS request. Those funcs named as Reply+nameOfMessage, the other functions there are helpers and they are not exported. (ex: func ReplyNewPost() creates reply data for the request with type "newPost", func savePostToDB() is a helper for ReplyNewPost())
"Repliers" run function to work with DB and create data for WS reply
- if data returned by DB func is enough for the reply, return it (DB funcs return either primitive  or an object from package models)
- if not, creates a reply using struct from webmodels.

These "repliers" must be tied to a message type in **package wsrouter** (backend/pkg/routes/wsrouter/wsrouter.go), func CreateAPIwsrouters(),
like 
`webmodel.MessageType:   sendReplyForLoggedUser(wsconnection.FuncReplyCreator(controllers.ReplyMessageType)),` 
 where
- webmodel.MessageType - a const from **package webmodel**
- sendReplyForLoggedUser - midlware func to check if the user is logged in
- FuncReplyCreator - function, which creates and sends a message with 
	- the according type 
	- and payload like {"result":"success", "data":\<data from replier\>} (any errors payloads, with "result":"error" are created in repliers if an error occurs)

func ReadPump() from package wsconnection reads WS message and runs the "replier" "tied" to the type of the message.
So, creating a message to reply to WS request requires:  
- in **package queries** (backend/pkg/db/sqlite/queries)
	- add a func that will exec a DB query, if there is no appropriate one 
-  in **package webmodel**  (backend/pkg/webmodel)
	- define the constants for WS messages type
	- if needed define a structs which  WS *messages* requests are parsed to as well as Valide() method for the struct and func PayloadTo.... in **package parse** (backend/pkg/webmodel/parse)
- in **package controllers** (backend/pkg/controllers)
	- create a func Reply...
-  in **package wsrouter** (backend/pkg/routes/wsrouter)
	- add an instruction like 
	`webmodel.MessageType: sendReplyForLoggedUser(wsconnection.FuncReplyCreator(controllers.ReplyMessageType))` 
	to func CreateAPIwsrouters()


WS messages
if there is a **server error**, the reply will be like:
{
	type: "ERROR", 
	data:{
		result: "serverError", 
		data: \<errmessage\>
	}
}

if there is a **request error**, the reply will be like:
{
	type: \<typeOfRequest\>, 
	data:{
		result: "error", 
		data: \<errmessage\>
	}
}

**"successful"** replies (payload.result="success"):

*Type*: **"postsPortion"**
	*request payload*: offset: number
	*reply payload.data:* array of posts without comments
		[] {
		    id               string     
		    theme            string  
		    content  {
			    userID           string    
			    userName         string       
			    text             string       
			    dateCreate       time.Time
			    likes            []int       // index 0 keeps number of dislikes, index 1 keeps number of like
			    userReaction     int //-1 => no reaction, 0=>current user disliked, 1=>current user liked
			}    
		    categories          [] string   
		    commentsQuantity    int 
			groupID 			string		    
		}

*Type*: **""userPostsPortion""**
	*request payload*: webmodel.UserPosts  
		{
			UserID string `json:"userID"`
			Offset int    `json:"offset"`
		}
	*reply payload.data:* array of posts without comments
		[] {
		    id               string     
		    theme            string  
		    content  {
			    userID           string    
			    userName         string       
			    text             string       
			    dateCreate       time.Time
			    likes            []int       // index 0 keeps number of dislikes, index 1 keeps number of like
			    userReaction     int //-1 => no reaction, 0=>current user disliked, 1=>current user liked
			}    
		    categories          [] string   
		    commentsQuantity    int 
			groupID 			string		    
		}

*Type*: **"fullPostAndComments"**
	*request payload*:  postId: string
	*reply payload.data:* post with comments
	{
		id        string     
		theme     string  
		content  {
			userID        string    
			userName      string       
			text          string       
			dateCreate    time.Time
			likes         []int       // index 0 keeps number of dislikes, index 1 keeps number of like
			userReaction  int //-1 => no reaction, 0=>current user disliked, 1=>current user liked
		}    
		categories [] string   
		comments   []{
			id      string  
			postID  string  
			content  {
				userID         string    
				userName       string       
				text           string       
				dateCreate     time.Time
				likes          []int       // index 0 keeps number of dislikes, index 1 keeps number of like
				userReaction   int  //-1 => no reaction, 0=>current user disliked, 1=>current user liked
			}
		}
		commentsQuantity    int  
		groupID 			string   
	}

*Type*: **"newPost"**
	*request payload*: webmodel.Post
		{
			title      string   
		    body       string   
		    categories []string
		}
	*reply payload.data:* post without comments
		{
		    id               string     
		    theme            string  
		    content  {
			    userID           string    
			    userName         string       
			    text             string       
			    dateCreate       time.Time
			    likes            []int       // index 0 keeps number of dislikes, index 1 keeps number of like
			    userReaction     int //-1 => no reaction, 0=>current user disliked, 1=>current user liked
			}    
		    categories          [] string   
		    commentsQuantity    int 
			groupID 			string		    
		}

*Type*:  **"newComment"**
	*request payload*:  webmodel.Comment
		{
			postID     string 
			content    string 
		}
	*reply payload.data:* post with comments
		{
			id                  string     
			theme               string  
			content  {
				userID       string    
				userName     string       
				text         string       
				dateCreate   time.Time
				likes        []int       // index 0 keeps number of dislikes, index 1 keeps number of like
				userReaction int //-1 => no reaction, 0=>current user disliked, 1=>current user liked
			}    
			categories [] string   
			comments   []{
				id      string  
				postID  string  
				content  {
					userID       string    
					userName     string       
					text         string       
					dateCreate   time.Time
					likes        []int       // index 0 keeps number of dislikes, index 1 keeps number of like
					userReaction int //-1 => no reaction, 0=>current user disliked, 1=>current user liked
				}
			}
			commentsQuantity    int 
			groupID 			string	    
		}	

*Type*:  **"onlineUsers"**
	*request payload*:  - 
	*reply payload.data:* 
		{
			id              string
			userName        string
			lastMessageDate string
		}

*Type*:  **"openChat"**
	*request payload*:  recepientID: string
	*reply payload.data:* privateChat
		{
			id                 string               
			name               string               
			recipientUserID    string               
			recipientUserName  string               
			messages       []{
				id          string   
				userID      string   
				userName    string   
				content     string   
				dateCreate  time.Time
			}
		}

*Type*:  **"sendMessageToOpendChat"**
	*request payload*:  webmodel.ChatMessage
		{
			content       string   
			dateCreate    time.Time
		}
	*reply payload.data:* 
		*to creator payload.data*: "delivered"
		*to recepient *webmodel.ChatMessage
	
			*Type*:  **"inputChatMessage"**
	 			payload.data
					{
						id         string   
						userID      string   //author
						userName    string   
						content     string   
						dateCreate  time.Time
					}


*Type*:  **"closeChat"**
	*request payload*:  null
	*reply payload.data:*  - // server will delete info about "current opened chat" 

*Type*:  **"chatPortion"**  
	*request payload*: offset: number
	*reply payload.data:* privateChat // may by change to - send just array of messages
		{
			id                  string               
			name                string               
			recipientUserID     string               
			recipientUserName 	string               
			messages                []{ 
				id              string   
				userID          string   
				userName        string   
				content         string   
				dateCreate      time.Time
			}
		}
*Type*:  **"newOnlineUser"**
	*request payload*: -
	*reply payload.data:*  *to other users*: user
		{
		   id                     string  
		   userName        string
		}

*Type*:  **"offlineUser"** 
	*request payload*: -
	*reply payload.data:*  *to other users*: user
		{
		   id                     string  
		   userName        string
		}

*Type*:  **"createGroup"**
	*request payload*:  webmodel.Group
		{
			title       string
			description string
		}
	*reply payload.data:* group
		{
			id          string    
			name        string    
			description string    
			creatorID   string    
			dateCreate  time.Time 
			members     []{
				id       string 
				userName string 
			}
		}

*Type*:  **"addUserToGroup"**
	*request payload*:  webmodel.UserInGroup
		{
			groupID string
			userID  string
			userName string
		}
	*reply payload.data:* webmodel.UserInGroup
		{
			groupID string
			userID  string
			userName string
		}

*Type*:  **"groupMembers"**
	*request payload*: groupID  string
	*reply payload.data:* user
		{
		   id              string  
		   userName        string
		}

*Type*:  **"userGroups"**
	*request payload*: userID  string
	*reply payload.data:* []group
		{
			id          string    
			name        string    
		}


*Type*:  **"getUserProfile"**
	*request payload*: userID  string
	*reply payload.data:* user
		{
			id              string   
			userName      	string   
			email           string   
			dateCreate      time.Time
			dateBirth       time.Time
			gender          string   
			firstName       string   
			lastName        string 
			profileType     int // 0 - public, 1 - private  
		}

*Type*: **"deletePost"**
	*request payload*:  postId: string
	*reply payload.data:* postId 

*Type*:  **"postsPortionInGroup"**
	*request payload*: webmodel.PostsInGroupPortion
		{
			groupID string 	
			offset  int    
		}
	*reply payload.data:* usarray of posts without commentser
		[] {
	    	id               string     
	    	theme            string  
	    	content  {
			    userID           string    
			    userName         string       
			    text             string       
			    dateCreate       time.Time
			    likes            []int       // index 0 keeps number of dislikes, index 1 keeps number of like
			    userReaction     int //-1 => no reaction, 0=>current user disliked, 1=>current user liked
			}    
	    	categories          [] string   
	    	commentsQuantity    int 
			groupID 			string		    
		}

*Type*:  **"userFollowers"**
	*request payload*: userID  string
	*reply payload.data:* user
		{
		  id              string  
		   userName       string
		}

*Type*:  **"userFollowing"**
	*request payload*: userID  string
	*reply payload.data:* user
		{
		  id              string  
		   userName       string
		}

*Type*:  **"followUser"**
	*request payload*: userID  string
	*reply payload.data:* followStatus // "requested" or "following"

*Type*:  **"getFollowStatus"**
	*request payload*: userID  string
	*reply payload.data:* followStatus // "requested" or "following"

*Type*:  **"searchGroupsUsers"**
	*request payload*: searchQuery  string
	*reply payload.data: 
		{
			entityType string  // groups or profile
			id         string 
			name       string 
		}
	
*Type*:  **"setProfileType"**
	*request payload*: profileType  string
	*reply payload.data:* profileType  // 0 - public, 1 - private

*Type*:  **"newNotification"**
	*request payload*:  - 
	*reply payload.data:* Notification
		{
			id 	          int           
			userId 	      string        
			fromUserId 	  sql.NullString
			type 	      string        
			body 	      string        
			groupId 	  sql.NullString
			postId 	  	  sql.NullString
			rea 	  	  int           
			dateCreate 	  time.Time     
		}

*Type*:  **"CreateGroupEvent"**
	*request payload*:  webmodel.Event
		{
			title       string    
			Description string    
			GroupID     string    
			DateCreate  time.Time  // omitempty"`
			DateEvent   time.Time 
		}
	*reply payload.data:* Event
		{
			id          string    
			title       string    
			description string    
			groupID     string    
			creatorID   string    
			dateEvent   time.Time 
			dateCreate  time.Time 
		}
		