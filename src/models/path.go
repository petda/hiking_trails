package models

import (
	"database/sql"
	"fmt"
	"github.com/martini-contrib/binding"
	"net/http"
)

// places (array) An array of places along the path (trail).
// id (int) Path id.
// name (string) The name of the path.
// info (string) Description of the path.
// length (string) Path length in km.
// polyline (array) Path as an array of geo coordinates objects with lat and lng.
// duration (string) Path hiking time in hours.
// image (string) URL to an image describing the trail.

type Path struct {
	Id       int64          `json:"id"`
	Name     string         `json:"name"       binding:"required"`
	Info     string         `json:"info"`
	Length   string         `json:"length"`
	Polyline GEOCoordinates `json:"polyline"`
	Duration string         `json:"duration"`
	Places   []*Place       `json:"places"`
	ImageURL string         `json:"image"`
	BundleId int64          `json:"bundleId,omitempty"`
}

func NewPath() *Path {
	path := &Path{}
	path.Places = make([]*Place, 0)
	path.Polyline = NewGEOCoordinates()

	return path
}

// The martini binding plugin does not use pointer targets. Therefore define
// validate method on struct.
func (path Path) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	// FIXME A lot more validation here.
	validateStringLength("name", path.Name, 1, 255, &errors)
	validateStringLength("info", path.Info, 0, 255, &errors)
	validateStringLength("length", path.Length, 0, 255, &errors)
	validateStringLength("image", path.ImageURL, 0, 255, &errors)
	validateStringLength("duration", path.Duration, 0, 255, &errors)

	return errors
}

func (path *Path) Type() string {
	return "path"
}

func (path *Path) DatabaseTable() string {
	return "paths"
}

func (path *Path) RequireTransaction() bool {
	return true
}

func (path *Path) Save(execer SQLExecer) error {
	result, err := execer.Exec("INSERT INTO paths(name, info, length, polyline, duration, image_url, bundle_id) VALUES(?,?,?,?,?,?,?)",
		path.Name,
		path.Info,
		path.Length,
		path.Polyline.AsBytes(),
		path.Duration,
		path.ImageURL,
		path.BundleId)

	if err != nil {
		return NewAPIError(500, "Failed to create path", err)
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		return NewAPIError(500, "Failed to retrieve last inserted id when saving path", err)
	}

	path.Id = lastInsertedId

	for _, place := range path.Places {
		place.PathId = path.Id

		// FIXME Use prepered statements here instead?
		err = place.Save(execer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (path *Path) Load(queryer SQLQueryer) error {
	polylineData := make([]byte, 0)

	err := queryer.QueryRow("SELECT id, name, info, length, polyline, duration, image_url, bundle_id FROM paths WHERE id=?", path.Id).
		Scan(&path.Id, &path.Name, &path.Info, &path.Length, &polylineData, &path.Duration, &path.ImageURL, &path.BundleId)

	if err == sql.ErrNoRows {
		return NewAPIError(404, fmt.Sprintf("No path with id %d exist", path.Id), nil)
	} else if err != nil {
		return NewAPIError(500, fmt.Sprintf("Failed to load path with id %d", path.Id), err)
	}

	polyline, err := GEOCoordinatesFromBytes(polylineData)
	if err != nil {
		return NewAPIError(500, fmt.Sprintf("Failed to load path with id %d", path.Id), err)
	}

	path.Polyline = polyline

	filter := map[string]interface{}{"path_id": path.Id}
	places, err := LoadPlaces(queryer, filter)
	if err != nil {
		return err
	}

	path.Places = places

	return nil
}

func (path *Path) Update(execer SQLExecer) error {
	result, err := execer.Exec("UPDATE paths SET name=?, info=?, length=?, polyline=?, duration=?, image_url=?, bundle_id=? WHERE id=?",
		path.Name,
		path.Info,
		path.Length,
		path.Polyline.AsBytes(),
		path.Duration, path.ImageURL,
		path.BundleId,
		path.Id)

	if err != nil {
		return NewAPIError(500, fmt.Sprintf("Failed to update path with id %d", path.Id), err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return NewAPIError(500, fmt.Sprintf("Failed to update path with id %d", path.Id), err)
	}

	if rowsAffected == 0 {
		return NewAPIError(404, fmt.Sprintf("No path with id %d exist", path.Id), nil)
	}

	return nil
}

func (path *Path) Delete(execer SQLExecer) error {
	return nil
}

func LoadPathsFromDatabase(transaction SQLQueryer, bundleId int64) ([]*Path, error) {
	var rows *sql.Rows
	var err error
	paths := make([]*Path, 0)

	if bundleId == 0 {
		rows, err = transaction.Query("SELECT id, name, info, length, polyline, duration, image_url, bundle_id FROM paths")
	} else {
		rows, err = transaction.Query("SELECT id, name, info, length, polyline, duration, image_url, bundle_id FROM paths WHERE bundle_id=?", bundleId)

	}

	if err != nil {
		return nil, NewAPIError(500, "Failed to load paths from database", err)
	}
	defer rows.Close()

	for rows.Next() {
		path := &Path{}
		polylineData := make([]byte, 0)

		err = rows.Scan(&path.Id, &path.Name, &path.Info, &path.Length, &polylineData, &path.Duration, &path.ImageURL, &path.BundleId)
		if err != nil {
			return nil, err
		}

		polyline, err := GEOCoordinatesFromBytes(polylineData)
		if err != nil {
			return nil, NewAPIError(500, "Failed to load path from row", err)
		}

		path.Polyline = polyline

		filter := map[string]interface{}{"path_id": path.Id}
		places, err := LoadPlaces(transaction, filter)
		if err != nil {
			return nil, NewAPIError(500, "Failed to load path from row", err)
		}

		path.Places = places
		paths = append(paths, path)
	}

	if err := rows.Err(); err != nil {
		return nil, NewAPIError(500, "Failed to load paths from database", err)
	}

	return paths, nil
}
