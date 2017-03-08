package controllers

import (
	"database/sql"
	"github.com/martini-contrib/render"
	"hiking_trails/src/middleware"
	"hiking_trails/src/models"
	"log"
	"net/http"
	"time"
)

func UsersControllerLogin(form models.LoginForm, store *middleware.SessionStore, render render.Render, db *sql.DB, logger *log.Logger, response http.ResponseWriter) {
	user := &models.User{}

	err := user.LoadFromUsername(form.Username, db)
	if err != nil {
		err := models.NewAPIError(401, "Invalid username and password", nil)
		renderErrorAsJson(err, render, logger)
		return
	}

	if !user.IsCorrectPassword(form.Password) {
		err := models.NewAPIError(401, "Invalid username and password", nil)
		renderErrorAsJson(err, render, logger)
		return
	}

	session := middleware.NewSession(store)
	session.Set("userId", user.Id)
	session.Set("isAdministrator", user.IsAdministrator)
	session.Create()

	cookie := &http.Cookie{
		Name:    "SessionId",
		Value:   session.Id,
		Path:    "/",
		Expires: time.Now().Add(48 * time.Hour),
	}
	http.SetCookie(response, cookie)

	responseData := map[string]string{"SessionId": session.Id}
	render.JSON(200, responseData)
}

func UsersControllerLogout(session *middleware.Session, render render.Render, response http.ResponseWriter) {
	session.Delete()

	// Set session cookie to expired(1970...) so browser knows to remove it.
	cookie := &http.Cookie{
		Name:    "SessionId",
		Value:   "",
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Time{},
	}

	http.SetCookie(response, cookie)

	render.JSON(200, "")
}
