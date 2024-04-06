package queries

import (
	"fmt"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
)

func (dbm *DBModel) GetListOfGroupsUsers(searchQuery string) ([]models.SearchResult, error) {
	results := []models.SearchResult{}
	searchQuery = "%" + searchQuery + "%"
	q := `SELECT id, title, 'groups'
			FROM groups
			WHERE title LIKE ?
			UNION ALL
			SELECT id, userName, 'profile'
			FROM users
			WHERE userName LIKE ?
	`
	rows, err := dbm.DB.Query(q, searchQuery, searchQuery)
	if err != nil {
		return results, fmt.Errorf("GetListOfSearchResults query failed: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var result models.SearchResult
		err := rows.Scan(&result.ID, &result.Name, &result.EntityType)
		if err != nil {
			return results, fmt.Errorf("GetListOfSearchResults row scan failed: %w", err)
		}

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return results, fmt.Errorf("GetListOfSearchResults: row scan failed: %w", err)
	}

	return results, nil
}

func (dbm *DBModel) GetListOfUsersNotFriends(searchQuery, userID string) ([]models.SearchResult, error) {
	results := []models.SearchResult{}
	searchQuery = "%" + searchQuery + "%"
	q := `SELECT u.id, u.userName FROM users u
      LEFT JOIN friends f ON u.id = f.friendUser AND f.mainUser = ?
      INNER JOIN followers fo ON u.id = fo.followerID AND fo.followingID = ?
      WHERE u.userName LIKE ? AND f.friendUser IS NULL`

	rows, err := dbm.DB.Query(q, userID, userID, searchQuery) // selects a portion of users who are currently following you and who are also NOT yet friends
	if err != nil {
		return results, err
	}
	defer rows.Close()
	for rows.Next() {
		var result models.SearchResult
		err := rows.Scan(&result.ID, &result.Name)
		if err != nil {
			return results, fmt.Errorf("GetListOfUsersNotFriends: row scan failed: %w", err)
		}

		results = append(results, result)
	}
	
	if err := rows.Err(); err != nil {
		return results, fmt.Errorf("GetListOfUsersNotFriends: row scan failed: %w", err)
	}

	return results, nil
}

// old query
/* 	q := `SELECT id, userName
		FROM users
		WHERE userName LIKE ? AND id NOT IN (
			SELECT friendUser
			FROM friends
			WHERE mainUser = ?
		) AND id IN (
			SELECT followerID from followers WHERE followingID = ?
		)
` */
