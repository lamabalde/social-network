package queries

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/helpers"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
)

func (dbm *DBModel) CreateGroup(title, description, creatorID string, dateCreate time.Time) (string, error) {
	groupID, err := helpers.GenerateNewUUID()
	if err != nil {
		return "", fmt.Errorf("CreateGroup: generate new UUID failed failed: %w", err)
	}

	q := `INSERT INTO groups (id, title, description, creatorID, dateCreate) VALUES (?,?,?,?,?);
	INSERT INTO group_members (groupID, userID, isMember) VALUES (?,?, true);
	`
	err = dbm.runTransations(q, groupID, title, description, creatorID, dateCreate, groupID, creatorID)
	if err != nil {
		return "", fmt.Errorf("CreateGroup: run transaction failed: %w", err)
	}

	return groupID, nil
}

func (dbm *DBModel) AddUserToGroup(groupID, userID string, isMember bool) error {
	q := `INSERT INTO group_members (groupID, userID, isMember) VALUES (?,?,?);`
	_, err := dbm.DB.Exec(q, groupID, userID, isMember)
	if err != nil {
		return fmt.Errorf("AddUserToGroup failed: %w", err)
	}
	return nil
}

func (dbm *DBModel) SetGroupMemberStatus(groupID, userID string, isMember bool) error {
	q := `UPDATE group_members SET isMember=? WHERE  groupID=? AND userID=?`

	res, err := dbm.DB.Exec(q, isMember, groupID, userID)
	if err != nil {
		return fmt.Errorf("SetMemberStatus failed: %w", err)
	}

	return dbm.checkUnique(res)
}

func (dbm *DBModel) DeletUserFromGroup(groupID, userID string) error {
	q := `DELETE FROM group_members WHERE groupID=? AND userID=?`
	_, err := dbm.DB.Exec(q, groupID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (dbm *DBModel) GetGroupNameByID(id string) (string, error) {
	var groupName string
	err := dbm.DB.QueryRow(`SELECT title FROM 'groups' WHERE id=?`, id).Scan(&groupName)
	if err != nil {
		if err == sql.ErrNoRows {
			err = models.ErrNoRecords
		}
		return "", fmt.Errorf("GetGroupNameByID: row scan failed: %w", err)
	}

	return groupName, nil
}

/*
returns group with its members
*/
func (dbm *DBModel) GetGroupByID(id string) (models.Group, error) {
	var group models.Group

	q := `SELECT groups.title, groups.description, groups.creatorID, groups.dateCreate
	FROM 'groups' 
	WHERE groups.id = ?
	`

	err := dbm.DB.QueryRow(q, id).Scan(&group.Name, &group.Description, &group.CreatorID, &group.DateCreate)
	if err != nil {
		if err == sql.ErrNoRows {
			err = models.ErrNoRecords
		}
		return group, fmt.Errorf("GetGroupByID: row scan failed: %w", err)
	}

	group.ID = id

	group.Members, err = dbm.GetListOfGroupMembers(id)
	if err != nil {
		return group, fmt.Errorf("GetGroupByID: getting list of members failed: %w", err)
	}

	return group, nil
}

func (dbm *DBModel) GetListOfGroupMembers(groupID string) ([]models.UserBase, error) {
	var members []models.UserBase
	q := `SELECT userID, users.userName FROM group_members
	LEFT JOIN users ON group_members.userID=users.id 
		 WHERE groupID=? AND isMember`

	rows, err := dbm.DB.Query(q, groupID)
	if err != nil {
		return members, fmt.Errorf("GetListOfGroupMembers: query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.UserBase
		err := rows.Scan(&user.ID, &user.UserName)
		if err != nil {
			return members, fmt.Errorf("GetListOfGroupMembers: row scan failed: %w", err)
		}

		members = append(members, user)
	}

	if err := rows.Err(); err != nil {
		return members, fmt.Errorf("GetListOfGroupMembers: row scan failed: %w", err)
	}

	return members, nil
}

func (dbm *DBModel) GetGroupMembersIDs(groupID string) ([]string, error) {
	var membersIDs []string
	q := `SELECT userID FROM group_members
		 WHERE groupID=? AND isMember`

	rows, err := dbm.DB.Query(q, groupID)
	if err != nil {
		return membersIDs, fmt.Errorf("GetGroupMembersIDs: query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var userID string
		err := rows.Scan(&userID)
		if err != nil {
			return membersIDs, fmt.Errorf("GetGroupMembersIDs: row scan failed: %w", err)
		}

		membersIDs = append(membersIDs, userID)
	}

	if err := rows.Err(); err != nil {
		return membersIDs, fmt.Errorf("GetGroupMembersIDs: row scan failed: %w", err)
	}

	return membersIDs, nil
}

/*
returns list of groups which the user belongs to
*/
func (dbm *DBModel) GetListOfUserGroups(userID string) ([]models.Group, error) {
	groups := []models.Group{}
	q := `SELECT groupID, groups.title
			FROM group_members
			LEFT JOIN  groups ON group_members.groupID = groups.id
			WHERE userID = ? AND isMember
	`
	rows, err := dbm.DB.Query(q, userID)
	if err != nil {
		return groups, err
	}
	defer rows.Close()
	for rows.Next() {
		var group models.Group
		err := rows.Scan(&group.ID, &group.Name)
		if err != nil {
			return groups, fmt.Errorf("GetListOfUserGroups: row scan failed: %w", err)
		}

		groups = append(groups, group)
	}

	return groups, nil
}

/*
adds a new group message into DB, returns an ID for the message
*/
func (dbm *DBModel) AddGroupChatMessage(groupID string, authorID string, content string, images []string, dateCreate time.Time) (string, error) {
	id, err := helpers.GenerateNewUUID()
	if err != nil {
		return "", fmt.Errorf("AddGroupMessage: generate new UUID failed: %w", err)
	}

	groupMembersID, err := dbm.GetGroupMembersID(groupID, authorID)
	if err == models.ErrNoRecords {
		return "", models.ErrAddConstaints
	}
	if err != nil {
		return "", fmt.Errorf("AddGroupMessage: getting a groupMembersID failed: %w", err)
	}

	strOfImages := helpers.JoinToNullString(images)

	q := `INSERT INTO group_chat (id, content, images, group_membersID, dateCreate) VALUES (?,?,?,?,?)`
	_, err = dbm.DB.Exec(q, id, content, strOfImages, groupMembersID, dateCreate)
	if err != nil {
		return "", fmt.Errorf("AddGroupMessage: exec query failed: %w", err)
	}

	return id, nil
}

func (dbm *DBModel) DeleteGroupChatMessage(id string) error {
	q := `DELETE FROM group_chat WHERE id=?`
	_, err := dbm.DB.Exec(q, id)
	if err != nil {
		return err
	}
	return nil
}

/*
returns 'limit' messages of the group chat with the given 'id', skips 'offset' last messages
*/
func (dbm *DBModel) GetGroupChatMessagesByGroupId(id string, limit, offset int) ([]*models.ChatMessage, error) {
	condition, arguments := createConditionSelectMessagesByGroupId(id, limit, offset)
	query := createQuerySelectGroupChatMessages(
		" group_chat.id, group_chat.content, group_chat.images, group_chat.dateCreate, group_members.userID, users.userName ",
		condition,
	)
	groupChatMessages, err := dbm.execQuerySelectChatMessages(query, arguments, scanRowToChatMessage)
	if err != nil {
		return nil, err
	}
	return groupChatMessages, nil
}

/*
creates a condition for the query to select group messages.
Used in the GetGroupMessagesByGroupId function.
*/
func createConditionSelectMessagesByGroupId(groupID string, limit, offset int) (string, []any) {
	condition := ` WHERE groups.id = ? `
	arguments := []any{groupID}

	arguments = append(arguments, limit, offset)

	return condition, arguments
}

/*
creates the query to select group messages.
Used in the GetGroupMessagesByGroupId and GetGroupMessagesByUsersId function.
*/
func createQuerySelectGroupChatMessages(fields, condition string) (query string) {
	query = `
		SELECT  ` + fields + ` 
	    FROM group_chat
	    INNER JOIN group_members ON group_chat.group_membersID=group_members.id  AND isMember
		LEFT JOIN groups ON group_members.groupID=groups.id 
		LEFT JOIN users ON group_members.userID=users.id 
		` + condition + ` 
		ORDER BY group_chat.dateCreate DESC
		LIMIT ? OFFSET ?;
		`
	return
}

/*
returns id of the row in group_members table for a given group and user
used in AddGroupMessage
*/
func (dbm *DBModel) GetGroupMembersID(groupID string, authorID string) (int, error) {
	var id int
	q := `SELECT id FROM group_members WHERE groupID=? AND userID=? AND isMember`
	row := dbm.DB.QueryRow(q, groupID, authorID)

	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, err // models.ErrAddConstaints ?
		}
		return 0, err
	}

	return id, nil
}

func (dbm *DBModel) CheckGroupMemberStatus(groupID string, authorID string) (int, error) {
	var memberStatus int
	q := `SELECT isMember FROM group_members WHERE groupID=? AND userID=?`
	row := dbm.DB.QueryRow(q, groupID, authorID)

	err := row.Scan(&memberStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, nil // models.ErrAddConstaints ?
		}
		return -1, err
	}

	return memberStatus, nil
}
