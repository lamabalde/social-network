package controllers

import (
	"net/http"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
)

/*
	func getCategories(app *application.Application, viewVars *viewVarsMap) (allCategories []*models.Category, err error) {
		allCategories, err = app.DBModel.GetCategories()
		if err != nil {
			return nil, fmt.Errorf("getting data (set of categories) from DB failed: %w", err)
		}
		(*viewVars)["AllCategories"] = allCategories
		return
	}
*/
func getFilters(r *http.Request, user *models.User) (filter models.Filter, err error) {
	// get category filters
	uQ := r.URL.Query()

	// get author's filters
	if user != nil {
		if uQ.Get(F_AUTHORID) != "" {
			filter.AuthorID = user.ID
		}
		if uQ.Get(F_LIKEBY) != "" {
			filter.LikedByUserID = user.ID
		}
		if uQ.Get(F_DISLIKEBY) != "" {
			filter.DisLikedByUserID = user.ID
		}

	}
	return
}
