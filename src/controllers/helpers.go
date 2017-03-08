package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"hiking_trails/src/models"
	"log"
	"strconv"
)

func MustGetIdFromParameters(params martini.Params, logger *log.Logger) (int64, error) {
	idString, exist := params["id"]
	if !exist {
		logger.Panicf("Parameter 'id' is not present in params. Router must be misconfigured.")
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		return 0, models.NewAPIError(400, fmt.Sprintf("%s is not a valid id.", idString), nil)
	}

	return int64(id), nil
}

func MustGetLastInsertedId(result sql.Result, logger *log.Logger) int64 {
	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		logger.Panicf("Failed to retrieve last inserted id: %s", err)
	}

	return lastInsertedId
}

func MustGetRowsAffected(result sql.Result, logger *log.Logger) int64 {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Panicf("Failed to retrieve affected rows: %s", err)
	}

	return rowsAffected
}

func LogAndRenderError500(logger *log.Logger, render render.Render, message string, causedBy error) {
	logger.Printf("%s: %s", message, causedBy)
	apiError := models.NewAPIError(500, "Internal Server Error", causedBy)
	render.JSON(500, apiError)
}

func printAsJson(data interface{}) {
	asJson, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Failed to print data(%#v) as json: %s", data, err)
	} else {
		fmt.Printf("%s\n", string(asJson))
	}
}

func renderErrorAsJson(err error, render render.Render, logger *log.Logger) {
	apiError, isApiError := err.(*models.APIError)

	if !isApiError {
		apiError = models.NewAPIError(500, "Internal Server Error", err)
	}

	apiError.RenderAsJson(render, logger)
}
