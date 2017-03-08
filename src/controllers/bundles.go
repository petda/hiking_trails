package controllers

import (
	"database/sql"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"hiking_trails/src/models"
	"log"
)

func BundlesControllerCreate(bundle models.Bundle, render render.Render, db *sql.DB, logger *log.Logger) {
	err := models.Save(&bundle, db)

	if err != nil {
		LogAndRenderError500(logger, render, "Failed to insert bundle into database", err)
		return
	}

	render.JSON(201, bundle)
}

func BundlesControllerRead(params martini.Params, render render.Render, db *sql.DB, logger *log.Logger) {
	id, err := MustGetIdFromParameters(params, logger)
	if err != nil {
		renderErrorAsJson(err, render, logger)
		return
	}

	bundle := models.NewBundle()
	bundle.Id = id
	err = models.Load(bundle, db)

	if err != nil {
		renderErrorAsJson(err, render, logger)
		return
	}

	render.JSON(200, bundle)
}

func BundlesControllerUpdate(params martini.Params, bundle models.Bundle,
	render render.Render, db *sql.DB, logger *log.Logger) {

	id, err := MustGetIdFromParameters(params, logger)
	if err != nil {
		renderErrorAsJson(err, render, logger)
		return
	}

	if id != bundle.Id {
		err = models.NewAPIError(400, "Not allowed to change bundle id", nil)
	}

	if err == nil {
		err = models.Update(&bundle, db)
	}

	if err != nil {
		renderErrorAsJson(err, render, logger)
		return
	}

	render.JSON(200, bundle)
}

func BundlesControllerDelete(params martini.Params, render render.Render, db *sql.DB,
	logger *log.Logger) {

	id, err := MustGetIdFromParameters(params, logger)
	if err == nil {
		err = models.Delete(&models.Bundle{}, id, db)
	}

	if err != nil {
		renderErrorAsJson(err, render, logger)
		return
	}

	render.JSON(204, "")
}

func BundlesControllerList(render render.Render, db *sql.DB, logger *log.Logger) {
	transaction, err := db.Begin()
	if err != nil {
		err = models.NewAPIError(500, "Failed to begin transaction when reading bundle", err)
		renderErrorAsJson(err, render, logger)
		return
	}

	bundles, err := models.LoadBundles(transaction)

	if err == nil {
		err = transaction.Commit()
	} else {
		transaction.Rollback()
	}

	if err != nil {
		err = models.NewAPIError(500, "Failed to load bundles from database", err)
		renderErrorAsJson(err, render, logger)
		return
	}

	render.JSON(200, bundles)
}
