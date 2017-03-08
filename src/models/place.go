package models

import (
	"database/sql"
	"fmt"
	"github.com/martini-contrib/binding"
	"net/http"
)

// name (string) Place name.
// info (string) Place description.
// image (string) URL to image asset of the place.
// radius (int) Radius of place marker.
// position (object) Geo coordinates object with lat and lng properties.
// media (array) Array of additional media objects.

type Place struct {
	Id       int64         `json:"id"`
	Name     string        `json:"name"       binding:"required"`
	Info     string        `json:"info"`
	Radius   int64         `json:"radius"`
	Position GEOCoordinate `json:"position"`
	Media    []Media       `json:"media"`
	PathId   int64         `json:"pathId,omitempty"`
}

func NewPlace() *Place {
	place := &Place{}
	place.Media = make([]Media, 0)

	return place
}

// The martini binding plugin does not use pointer targets. Therefore define
// validate method on struct.
func (place Place) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	// FIXME A lot more validation here.
	validateStringLength("name", place.Name, 1, 255, &errors)
	validateStringLength("name", place.Info, 0, 255, &errors)

	return errors
}

func validateStringLength(field string, value string, min int, max int, errors *binding.Errors) {
	if len(value) < min || len(value) > max {
		*errors = append(*errors, NewBindingRangeError("name", min, max))
	}
}

func (place *Place) Type() string {
	return "place"
}

func (place *Place) DatabaseTable() string {
	return "places"
}

func (place *Place) RequireTransaction() bool {
	return false
}

func (place *Place) Save(execer SQLExecer) error {
	result, err := execer.Exec("INSERT INTO places(name, info, radius, position, path_id) VALUES(?,?,?,?,?)",
		place.Name,
		place.Info,
		place.Radius,
		place.Position.AsBytes(),
		place.PathId)

	if err != nil {
		return NewAPIError(500, "Failed to create place", err)
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("Failed to retrieve last inserted id: %s", err)
	}

	place.Id = lastInsertedId
	return nil
}

func (place *Place) Load(queryer SQLQueryer) error {
	positionData := make([]byte, 0, 4*2)

	err := queryer.QueryRow("SELECT id, name, info, radius, position, path_id FROM places WHERE id=?", place.Id).
		Scan(&place.Id, &place.Name, &place.Info, &place.Radius, &positionData, &place.PathId)

	if err == sql.ErrNoRows {
		return NewAPIError(404, fmt.Sprintf("No place with id %d exist", place.Id), nil)
	} else if err != nil {
		return NewAPIError(500, fmt.Sprintf("Failed to load place with id %d", place.Id), err)
	}

	position, err := GEOCoordinateFromBytes(positionData)
	if err != nil {
		return NewAPIError(500, fmt.Sprintf("Failed to load place with id %d", place.Id), err)
	}

	place.Position = *position
	return nil
}

func (place *Place) Update(execer SQLExecer) error {
	result, err := execer.Exec("UPDATE places SET name=?, info=?, radius=?, position=?, path_id=? WHERE id=?",
		place.Name, place.Info, place.Radius, place.Position.AsBytes(), place.PathId, place.Id)

	if err != nil {
		return NewAPIError(500, fmt.Sprintf("Failed to update place with id %d", place.Id), err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return NewAPIError(500, fmt.Sprintf("Failed to update place with id %d", place.Id), err)
	}

	if rowsAffected == 0 {
		return NewAPIError(404, fmt.Sprintf("No place with id %d exist", place.Id), nil)
	}

	return nil
}

func (place *Place) Delete(execer SQLExecer) error {
	return nil
}

func LoadPlaces(queryer SQLQueryer, filter map[string]interface{}) ([]*Place, error) {
	places := make([]*Place, 0)
	arguments := make([]interface{}, 0)
	queryStatement := "SELECT id, name, info, radius, position, path_id FROM places"

	path_id, path_id_in_filter := filter["path_id"]
	if path_id_in_filter {
		queryStatement += " WHERE path_id=?"
		arguments = append(arguments, path_id.(int64))
	}

	rows, err := queryer.Query(queryStatement, arguments...)

	if err != nil {
		return nil, NewAPIError(500, "Failed to load places", err)

	}
	defer rows.Close()

	for rows.Next() {
		place := NewPlace()
		positionData := make([]byte, 0, 8)

		err := rows.Scan(&place.Id, &place.Name, &place.Info, &place.Radius, &positionData, &place.PathId)
		if err != nil {
			return nil, NewAPIError(500, "Failed to load place from row", err)
		}

		position, err := GEOCoordinateFromBytes(positionData)
		if err != nil {
			return nil, NewAPIError(500, "Failed to load place from row", err)
		}

		place.Position = *position
		places = append(places, place)
	}

	err = rows.Err()
	if err != nil {
		return nil, NewAPIError(500, "Failed to load places from database", err)
	}

	return places, nil
}
