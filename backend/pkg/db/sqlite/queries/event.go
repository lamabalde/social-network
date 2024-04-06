package queries

import (
	"database/sql"
	"errors"
	"fmt"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/helpers"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
)

func (dbm *DBModel) CreateEvent(event models.Event) (string, error) {
	id, err := helpers.GenerateNewUUID()
	if err != nil {
		return "", fmt.Errorf("CreateEvent: generate new UUID failed failed: %w", err)
	}

	createByID, err := dbm.GetGroupMembersID(event.GroupID, event.CreatorID)
	if err == models.ErrNoRecords {
		return "", models.ErrAddConstaints
	}
	if err != nil {
		return "", fmt.Errorf("CreateEvent: getting a groupMembersID failed: %w", err)
	}

	q := `INSERT INTO 'group_event' (id, title, description, dateCreate, dateEvent, createByID) VALUES (?,?,?,?,?,?);`

	err = dbm.runTransations(q, id, event.Title, event.Description, event.DateCreate, event.DateEvent, createByID)
	if err != nil {
		return "", fmt.Errorf("CreateEvent: run transaction failed: %w", err)
	}

	return id, nil
}

func (dbm *DBModel) AddEventMember(eventID, userID string, mark int) (int, error) {
	groupMembersID, err := dbm.getGroupMembersIDForEvent(eventID, userID)
	if err == models.ErrNoRecords {
		return 0, models.ErrAddConstaints
	}
	if err != nil {
		return 0, fmt.Errorf("AddEventMember: getting a groupMembersID failed: %w", err)
	}

	q := `INSERT OR REPLACE INTO group_event_members (group_eventID, group_membersID, mark) VALUES (?,?,?)` // maybe it works?
	res, err := dbm.DB.Exec(q, eventID, groupMembersID, mark)
	if err != nil {
		return 0, fmt.Errorf("AddEventMember: exec query failed: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddEventMember: getting the new id failed: %w", err)
	}

	return int(id), nil
}

func (dbm *DBModel) ChangeEventMemberOption(eventID, userID string, mark int) error {
	groupMembersID, err := dbm.getGroupMembersIDForEvent(eventID, userID)
	if err == models.ErrNoRecords {
		return models.ErrAddConstaints
	}
	if err != nil {
		return fmt.Errorf("ChangeEventMemberOption: getting a groupMembersID failed: %w", err)
	}

	q := `UPDATE group_event_members SET mark=? WHERE group_eventID=? AND group_membersID=?`
	res, err := dbm.DB.Exec(q, mark, eventID, groupMembersID)
	if err != nil {
		return fmt.Errorf("ChangeEventMemberOption: exec query failed: %w", err)
	}

	return dbm.checkUnique(res)
}

func (dbm *DBModel) GetEventByID(id string) (models.Event, error) {
	var event models.Event

	q := `SELECT group_event.id, group_event.title, group_event.description, group_event.dateCreate, group_event.dateEvent,
		users.id, users.userName, group_members.groupID
		FROM 'group_event' 
		LEFT JOIN group_members ON group_event.createByID=group_members.id 
		LEFT JOIN users ON group_members.userID=users.id

		WHERE  group_event.id = ?
	`

	err := dbm.DB.QueryRow(q, id).
		Scan(&event.ID, &event.Title, &event.Description,
			&event.DateCreate, &event.DateEvent,
			&event.CreatorID, &event.CreatorName,
			&event.GroupID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = models.ErrNoRecords
		}
		return event, fmt.Errorf("GetEventByID: row scan failed: %w", err)
	}

	event.UserOptions, err = dbm.GetListOfUserOptionsForEvent(id)
	if err != nil {
		return event, fmt.Errorf("GetGroupByID: getting list of members failed: %w", err)
	}

	return event, nil
}

func (dbm *DBModel) GetListOfGroupEvents(groupID string) ([]models.Event, error) {
	var events []models.Event
	q := `SELECT group_event.id, group_event.title, group_event.description, group_event.dateCreate, group_event.dateEvent
		FROM 'group_event'
		LEFT JOIN group_members ON group_event.createByID=group_members.id
		WHERE group_members.groupID=? `

	rows, err := dbm.DB.Query(q, groupID)
	if err != nil {
		return events, fmt.Errorf("GetListOfGroupRvents: query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.Title, &event.Description,
			&event.DateCreate, &event.DateEvent)
		if err != nil {
			return events, fmt.Errorf("GetListOfGroupRvents: row scan failed: %w", err)
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return events, fmt.Errorf("GetListOfGroupRvents: row scan failed: %w", err)
	}

	return events, nil
}

func (dbm *DBModel) getGroupMembersIDForEvent(eventID, userID string) (int, error) {
	var id int
	q := `SELECT id FROM group_members gm 
		WHERE gm.userID=? AND 
			gm.groupID = (SELECT groupID FROM group_event ge LEFT JOIN group_members gm2 ON ge.createByID=gm2.id  WHERE ge.id =?)
	 `

	row := dbm.DB.QueryRow(q, userID, eventID)

	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrNoRecords
		}
		return 0, err
	}

	return id, nil
}

func (dbm *DBModel) GetListOfUserOptionsForEvent(eventID string) ([]models.UserOptionForEvent, error) {
	var userOptions []models.UserOptionForEvent
	q := `SELECT users.id, users.userName, group_event_members.mark
		FROM group_event_members  
		LEFT JOIN group_members ON group_event_members.group_membersID=group_members.id 
		LEFT JOIN users ON group_members.userID=users.id
		WHERE group_event_members.group_eventID=?  
	 `

	rows, err := dbm.DB.Query(q, eventID)
	if err != nil {
		return nil, fmt.Errorf("GetListOfUserOptionsForEvent: query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var userOption models.UserOptionForEvent
		err := rows.Scan(&userOption.UserID, &userOption.UserName, &userOption.Option)
		if err != nil {
			return nil, fmt.Errorf("GetListOfUserOptionsForEvent: scan failed: %w", err)
		}

		userOptions = append(userOptions, userOption)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetListOfUserOptionsForEvent: scan failed: %w", err)
	}

	return userOptions, nil
}
