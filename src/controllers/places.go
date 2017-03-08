package controllers

import (
	"database/sql"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"hiking_trails/src/models"
	"log"
)

func PlacesControllerCreate(place models.Place, render render.Render, db *sql.DB, logger *log.Logger) {
	err := models.Save(&place, db)

	if err != nil {
		LogAndRenderError500(logger, render, "Failed to insert place into database", err)
		return
	}

	render.JSON(201, place)
}

func PlacesControllerRead(params martini.Params, render render.Render, db *sql.DB, logger *log.Logger) {
	id, err := MustGetIdFromParameters(params, logger)
	if err != nil {
		renderErrorAsJson(err, render, logger)
		return
	}

	place := models.NewPlace()
	place.Id = id
	err = models.Load(place, db)

	if err != nil {
		renderErrorAsJson(err, render, logger)
		return
	}

	render.JSON(200, place)
}

func PlacesControllerUpdate(params martini.Params, place models.Place,
	render render.Render, db *sql.DB, logger *log.Logger) {

	id, err := MustGetIdFromParameters(params, logger)
	if err != nil {
		renderErrorAsJson(err, render, logger)
		return
	}

	if id != place.Id {
		err = models.NewAPIError(400, "Not allowed to change place id", nil)
	}

	if err == nil {
		err = models.Update(&place, db)
	}

	if err != nil {
		renderErrorAsJson(err, render, logger)
		return
	}

	render.JSON(200, place)
}

func PlacesControllerDelete(params martini.Params, render render.Render, db *sql.DB,
	logger *log.Logger) {

	id, err := MustGetIdFromParameters(params, logger)
	if err != nil {
		renderErrorAsJson(err, render, logger)
		return
	}

	err = models.Delete(&models.Place{}, id, db)

	if err != nil {
		renderErrorAsJson(err, render, logger)
		return
	}

	render.JSON(204, "")
}

func PlacesControllerList(render render.Render, db *sql.DB, logger *log.Logger) {
	places, err := models.LoadPlaces(db, nil)
	if err != nil {
		LogAndRenderError500(logger, render, "Got error when trying to list places", err)
		return
	}

	render.JSON(200, places)
}
