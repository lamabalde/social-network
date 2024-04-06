package queries

func (dbm *DBModel) AddCloseFriendToDB(friendID, currentUser string) error {
	q := `INSERT INTO friends (mainUser, friendUser)
	VALUES (?, ?)
	`
	_, err := dbm.DB.Exec(q, currentUser, friendID)
	if err != nil {
		return err
	}
	return nil
}
