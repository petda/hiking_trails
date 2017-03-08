package controllers

import (
	"database/sql"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"hiking_trails/src/models"
	"log"
)

func PathsControllerCreate(path models.Path, render render.Render, db *sql.DB, logger *log.Logger) {
	err := models.Save(&path, db)

	if err != nil {
		err = models.NewAPIError(500, "Failed to insert path into database", err)
		renderErrorAsJson(err, render, logger)
		return
	}

	render.JSON(201, path)
}

func PathsControllerRead(params martini.Params, render render.Render, db *sql.DB, logger *log.Logger) {
	id, err := MustGetIdFromParameters(params, logger)
	if err != nil {
		err = models.NewAPIError(500, "Failed to load path from database", err)
		renderErrorAsJson(err, render, logger)
		return
	}

	path := models.NewPath()
	path.Id = id
	err = models.Load(path, db)

	if err != nil {
		renderErrorAsJson(err, render, logger)
		return
	}

	render.JSON(200, path)
}

func PathsControllerUpdate(params martini.Params, path models.Path,
	render render.Render, db *sql.DB, logger *log.Logger) {

	id, err := MustGetIdFromParameters(params, logger)
	if err != nil {
		renderErrorAsJson(err, render, logger)
		return
	}

	if id != path.Id {
		err = models.NewAPIError(400, "Not allowed to change path id", nil)
	}

	if err == nil {
		err = models.Update(&path, db)
	}

	if err != nil {
		renderErrorAsJson(err, render, logger)
		return
	}

	render.JSON(200, path)
}

func PathsControllerDelete(params martini.Params, render render.Render, db *sql.DB,
	logger *log.Logger) {

	id, err := MustGetIdFromParameters(params, logger)
	if err != nil {
		renderErrorAsJson(err, render, logger)
		return
	}

	err = models.Delete(&models.Path{}, id, db)

	if err != nil {
		renderErrorAsJson(err, render, logger)
		return
	}

	render.Text(204, "")
}

func PathsControllerList(render render.Render, db *sql.DB, logger *log.Logger) {

	transaction, err := db.Begin()
	if err != nil {
		LogAndRenderError500(logger, render, "Failed to begin transaction when reading  paths", err)
		return
	}

	paths, err := models.LoadPathsFromDatabase(transaction, 0)

	if err == nil {
		err = transaction.Commit()
	} else {
		transaction.Rollback()
	}

	if err != nil {
		LogAndRenderError500(logger, render, "Failed to read paths from database", err)
		return
	}

	render.JSON(200, paths)
}
