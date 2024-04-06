package queries

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/helpers"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
)

const (
	ConstFields    = ` u.id, u.userName, u.email, u.dateCreate, u.dateBirth, u.gender, u.firstName, u.lastName, u.profileType, u.aboutMe `
	UserBaseFields = ` users.id, users.userName `
)

/*
returns list of all users in DB
*/
func (dbm *DBModel) GetAllUsers() ([]*models.UserBase, error) {
	return dbm.getUsersByCondition("")
}

/*
returns list of users who filtered by userIDs.CheckID(). userIDs has to implement CheckID method.
*/
func (dbm *DBModel) GetFilteredUsers(userIDs models.IdChecker) ([]*models.UserBase, error) {
	q := `SELECT ` + UserBaseFields + ` FROM users `
	rows, err := dbm.DB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.UserBase
	for rows.Next() {
		user := &models.UserBase{}
		err := rows.Scan(&user.ID, &user.UserName)
		if err != nil {
			return nil, err
		}
		if userIDs.CheckID(user.ID) {
			users = append(users, user)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (dbm *DBModel) getUsersByCondition(condition string, arguments ...any) ([]*models.UserBase, error) {
	q := `SELECT ` + UserBaseFields + ` FROM users ` + condition

	rows, err := dbm.DB.Query(q, arguments...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.UserBase
	for rows.Next() {
		user := &models.UserBase{}
		err := rows.Scan(&user.ID, &user.UserName)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (dbm *DBModel) GetUsersByPartialNameNotInGroup(searchQuery, groupID string) ([]*models.UserBase, error) {
	searchQuery = "%" + searchQuery + "%"
	condition := `WHERE userName LIKE ? AND 0 = (SELECT count(*) FROM group_members WHERE userID = users.id AND groupID = ?)`

	return dbm.getUsersByCondition(condition, searchQuery, groupID)
}

/*
returns a user from DB by ID
*/
func (dbm *DBModel) GetUserByID(id string) (*models.User, error) {
	q := `SELECT ` + ConstFields + ` 
	      FROM users u WHERE u.id=?`

	user := &models.User{}
	row := dbm.DB.QueryRow(q, id)
	err := row.Scan(&user.ID, &user.UserName, &user.Email, &user.DateCreate, &user.DateBirth, &user.Gender, &user.FirstName, &user.LastName, &user.ProfileType, &user.AboutMe)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecords
		}
		return nil, err
	}

	user.Followers, err = dbm.GetFollowers(user.ID)
	if err != nil {
		return nil, fmt.Errorf("getting followers for user %s(id:%s) failed: %v", user.UserName, user.ID, err)
	}

	user.Followings, err = dbm.GetFollowing(user.ID)
	if err != nil {
		return nil, fmt.Errorf("getting followings for user %s(id:%s) failed: %v", user.UserName, user.ID, err)
	}

	return user, nil
}

/*
returns a user from DB by the name
*/
func (dbm *DBModel) GetUserByName(userName string) (*models.User, error) {
	q := `SELECT ` + ConstFields + `  
	      FROM users u 
		  WHERE u.userName=?`

	user := &models.User{}
	row := dbm.DB.QueryRow(q, userName)
	err := row.Scan(&user.ID, &user.UserName, &user.Email, &user.DateCreate, &user.DateBirth, &user.Gender, &user.FirstName, &user.LastName, &user.ProfileType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecords
		}

		return nil, err
	}

	return user, nil
}

/*
returns a user from DB by the email
*/
func (dbm *DBModel) GetUserByEmail(email string) (*models.User, error) {
	q := `SELECT ` + ConstFields + `
	      FROM users u LEFT JOIN sessions s ON u.id=s.userID 
		  WHERE u.email=?`

	user := &models.User{}
	row := dbm.DB.QueryRow(q, email)
	err := row.Scan(&user.ID, &user.UserName, &user.Email, &user.DateCreate, &user.DateBirth, &user.Gender, &user.FirstName, &user.LastName, &user.ProfileType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecords
		}
		return nil, err
	}

	return user, nil
}

/*
returns a user from DB by the userName or email
*/
func (dbm *DBModel) GetUserByNameOrEmail(str string) (*models.User, error) {
	q := `SELECT ` + ConstFields + `, u.password_hash 
	      FROM users u 
		  WHERE u.userName = ? OR u.email = ?`

	user := &models.User{}
	row := dbm.DB.QueryRow(q, str, str)
	err := row.Scan(&user.ID, &user.UserName, &user.Email, &user.DateCreate, &user.DateBirth, &user.Gender, &user.FirstName, &user.LastName, &user.ProfileType, &user.AboutMe, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecords
		}
		return nil, err
	}

	return user, nil
}

/*
returns a user from DB by the email
*/
func (dbm *DBModel) GetUserBySession(uuid string) (*models.UserBase, string, time.Time, error) {
	q := `SELECT ` + UserBaseFields + `, s.id, s.expirySession 
	FROM  users INNER JOIN sessions s ON users.id=s.userID WHERE s.id=?`

	user := &models.UserBase{}
	var uuidInDB sql.NullString
	var expirySessionInDB sql.NullTime
	row := dbm.DB.QueryRow(q, uuid)
	err := row.Scan(&user.ID, &user.UserName, &uuidInDB, &expirySessionInDB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", expirySessionInDB.Time, models.ErrNoRecords
		}
		return nil, "", expirySessionInDB.Time, err
	}

	return user, uuidInDB.String, expirySessionInDB.Time, nil
}

/*
adds the user to DB
*/
func (dbm *DBModel) AddUser(user *models.User) (string, error) {
	userID, err := helpers.GenerateNewUUID()
	if err != nil {
		return "", err
	}
	user.ID = userID
	err = dbm.InsertUser(user)

	if err != nil {
		errUnique := dbm.CheckUserByName(user.UserName)
		if errUnique == nil {
			return "", models.ErrUniqueUserName
		}
		errUnique = dbm.CheckUserByEmail(user.Email)
		if errUnique == nil {
			return "", models.ErrUniqueUserEmail
		}
	}

	return userID, nil
}

/*
adds a session uuid to the user with the given ID
*/
func (dbm *DBModel) AddUserSession(userID string, expired time.Time) (string, error) {
	var err error

	sessionID, err := helpers.GenerateNewUUID()
	if err != nil {
		return "", err
	}

	q := `INSERT INTO sessions (id, userID, expirySession,agent) VALUES (?,?,?,?)`
	_, err = dbm.DB.Exec(q, sessionID, userID, expired, "")
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

/*
deletes the user's session id
*/
func (dbm *DBModel) DeleteUsersSession(id string) error {
	q := `DELETE FROM sessions WHERE id=?`
	res, err := dbm.DB.Exec(q, id)
	if err != nil {
		return err
	}

	return dbm.checkUnique(res)
}

/*
check if a user with the given name exists,  returns nil only if there is exactly one user
*/
func (dbm *DBModel) CheckUserByName(userName string) error {
	err := dbm.checkExisting("users", "userName", userName)
	if errors.Is(err, sql.ErrNoRows) {
		return models.ErrNoRecords
	}
	return err
}

/*
check if a user with the given email exists, returns nil only if there is exactly one user
*/
func (dbm *DBModel) CheckUserByEmail(email string) error {
	err := dbm.checkExisting("users", "email", email)
	if errors.Is(err, sql.ErrNoRows) {
		return models.ErrNoRecords
	}
	return err
}

/*
inserts the new user into DB. It doesn't do any check of unique data. But if DB have some restricts, it will return an error
*/
func (dbm *DBModel) InsertUser(user *models.User) error {
	q := `INSERT INTO users  (id, userName, email, password_hash, dateCreate, dateBirth, gender, firstName, lastName, profileType, aboutMe) VALUES (?,?,?,?,?,?,?,?,?,?, ?)`
	_, err := dbm.DB.Exec(q, user.ID, user.UserName, user.Email, user.Password, user.DateCreate, user.DateBirth, user.Gender, user.FirstName, user.LastName, user.ProfileType, user.AboutMe)
	if err != nil {
		return err
	}

	return nil
}

/*
changes an email of the user with the given id
*/
func (dbm *DBModel) ChangeUsersEmail(id string, email string) error {
	err := dbm.changeUsersField(id, "email", email)
	if err != nil {
		errUnique := dbm.CheckUserByEmail(email)
		if errUnique == nil {
			return models.ErrUniqueUserEmail
		}
	}
	return err
}

/*
changes a password of the user with the given id
*/
func (dbm *DBModel) ChangeUsersPassword(id string, password_hash string) error {
	return dbm.changeUsersField(id, "password_hash", password_hash)
}

/*
changes ProfileType of the user with the given id
*/
func (dbm *DBModel) SetUsersProfileType(id string, profileType int) error {
	// it doesnt't use changeUsersField, because profileType is int
	q := `UPDATE users SET profileType=? WHERE id=?`
	res, err := dbm.DB.Exec(q, profileType, id)
	if err != nil {
		return err
	}

	return dbm.checkUnique(res)
}

/*
changes a field in the users table for the user with the given id
*/
func (dbm *DBModel) changeUsersField(id string, field, value string) error {
	return dbm.setFieldStringWhereId("users", field, value, id)
}
