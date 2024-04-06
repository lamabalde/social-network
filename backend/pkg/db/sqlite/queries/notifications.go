package queries

import (
	"database/sql"
	"fmt"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
)

func (dbm *DBModel) AddNotification(note *models.Notification) (int, error) {
	q := `INSERT OR REPLACE INTO 'notifications' (userID, fromUserID, type, body, groupID, postID , read, dateCreate) VALUES (?,?,?,?,?,?,?,?)`
	res, err := dbm.DB.Exec(q, note.UserID, note.FromUserID, note.Type, note.Body, note.GroupID, note.PostID, note.Read, note.DateCreate)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (dbm *DBModel) DeleteNotification(id int) error {
	q := `DELETE FROM 'notifications' WHERE id=?`
	res, err := dbm.DB.Exec(q, id)
	if err != nil {
		return err
	}

	return dbm.checkUnique(res)
}

func (dbm *DBModel) GetUserUnReadNotification(userID string) ([]*models.Notification, error) {
	var notes []*models.Notification
	q := `SELECT id, userID, fromUserID, type, body, groupID, postID , read, dateCreate
			FROM 'notifications'
			WHERE userID = ? AND read=0
	`
	rows, err := dbm.DB.Query(q, userID)
	if err != nil {
		return notes, err
	}
	defer rows.Close()
	for rows.Next() {
		note := &models.Notification{}
		err := rows.Scan(&note.ID, &note.UserID, &note.FromUserID, &note.Type, &note.Body, &note.GroupID, &note.PostID, &note.Read, &note.DateCreate)
		if err != nil {
			return notes, fmt.Errorf("GetUserUnReadNotification: row scan failed: %w", err)
		}

		notes = append(notes, note)
	}

	return notes, nil
}

func (dbm *DBModel) GetNotificationByID(id int) (*models.Notification, error) {
	q := `SELECT id, userID, fromUserID, type, body, groupID, postID , read, dateCreate
			FROM 'notifications'
			WHERE id = ? 
	`
	note := &models.Notification{}
	err := dbm.DB.QueryRow(q, id).Scan(&note.ID, &note.UserID, &note.FromUserID, &note.Type, &note.Body, &note.GroupID, &note.PostID, &note.Read, &note.DateCreate)
	if err != nil {
		if err == sql.ErrNoRows {
			return note, models.ErrNoRecords
		}
		return note, err
	}

	return note, nil
}

func (dbm *DBModel) MarkNotificationRead(notificationID int) error {
	q := `UPDATE OR REPLACE notifications
	SET read = 1
	WHERE
	id = ?`

	res, err := dbm.DB.Exec(q, notificationID)
	if err != nil {
		return err
	}

	return dbm.checkUnique(res)
}
