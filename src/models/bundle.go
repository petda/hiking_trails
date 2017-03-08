package models

import (
	"database/sql"
	"fmt"
	"github.com/martini-contrib/binding"
	"net/http"
)

// id (int) Bundle id.
// name (string) Bundle name.
// image (string) URL to an image describing the bundle.
// info (string) A short descriptive text for the bundle.
// paths (array) Array of path objects (trails) in the bundle.

type Bundle struct {
	Id       int64   `json:"id"`
	Name     string  `json:"name"       binding:"required"`
	Info     string  `json:"info"`
	ImageURL string  `json:"image"`
	Paths    []*Path `json:"paths"`
}

func NewBundle() *Bundle {
	bundle := &Bundle{}
	bundle.Paths = make([]*Path, 0)

	return bundle
}

// The martini binding plugin does not use pointer targets. Therefore define
// validate method on struct.
func (bundle Bundle) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	// FIXME A lot more validation here.
	validateStringLength("name", bundle.Name, 1, 255, &errors)
	validateStringLength("info", bundle.Info, 0, 255, &errors)
	validateStringLength("image", bundle.ImageURL, 0, 255, &errors)

	return errors
}

func (bundle *Bundle) Type() string {
	return "bundle"
}

func (bundle *Bundle) DatabaseTable() string {
	return "bundles"
}

func (bundle *Bundle) RequireTransaction() bool {
	return true
}

func (bundle *Bundle) Save(execer SQLExecer) error {
	result, err := execer.Exec("INSERT INTO bundles(name, info, image_url) VALUES(?,?,?)",
		bundle.Name, bundle.Info, bundle.ImageURL)

	if err != nil {
		return NewAPIError(500, "Failed to create bundle", err)
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		return NewAPIError(500, "Failed to retrieve last inserted id when saving bundle", err)
	}

	bundle.Id = lastInsertedId

	for _, path := range bundle.Paths {
		path.BundleId = bundle.Id

		// FIXME Use prepered statements here instead?
		err = path.Save(execer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (bundle *Bundle) Load(queryer SQLQueryer) error {
	bundle.Paths = make([]*Path, 0)

	err := queryer.QueryRow("SELECT id, name, info, image_url FROM bundles WHERE id=?", bundle.Id).
		Scan(&bundle.Id, &bundle.Name, &bundle.Info, &bundle.ImageURL)

	if err == sql.ErrNoRows {
		return NewAPIError(404, fmt.Sprintf("No bundle with id %d exist", bundle.Id), nil)
	} else if err != nil {
		return NewAPIError(500, fmt.Sprintf("Failed to load bundle with id %d", bundle.Id), err)
	}

	paths, err := LoadPathsFromDatabase(queryer, bundle.Id)
	if err != nil {
		return err
	}
	bundle.Paths = paths

	return nil
}

func (bundle *Bundle) Update(execer SQLExecer) error {
	// FIXME Must update paths and places here to?
	result, err := execer.Exec("UPDATE bundles SET name=?, info=?, image_url=? WHERE id=?",
		bundle.Name, bundle.Info, bundle.ImageURL, bundle.Id)

	if err != nil {
		return NewAPIError(500, fmt.Sprintf("Failed to update bundle with id %d", bundle.Id), err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return NewAPIError(500, fmt.Sprintf("Failed to update bundle with id %d", bundle.Id), err)
	}

	if rowsAffected == 0 {
		return NewAPIError(404, fmt.Sprintf("No bundle with id %d exist", bundle.Id), nil)
	}

	return nil
}

func (bundle *Bundle) Delete(execer SQLExecer) error {
	result, err := execer.Exec("DELETE FROM bundles WHERE id=?", bundle.Id)
	if err != nil {
		return NewAPIError(500, fmt.Sprintf("Failed to delete bundle with id %d", bundle.Id), err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return NewAPIError(500, fmt.Sprintf("Failed to delete bundle with id %d", bundle.Id), err)
	}

	if rowsAffected == 0 {
		return NewAPIError(404, fmt.Sprintf("No bundle with id %d exist", bundle.Id), nil)
	}

	return nil
}

func LoadBundles(transaction *sql.Tx) ([]*Bundle, error) {
	bundles := make([]*Bundle, 0)

	// TODO: Implement support for limit and offset.
	rows, err := transaction.Query("SELECT id, name, info, image_url FROM bundles")
	if err != nil {
		return nil, NewAPIError(500, "Failed to load bundles", err)
	}
	defer rows.Close()

	for rows.Next() {
		bundle := NewBundle()

		err := rows.Scan(&bundle.Id, &bundle.Name, &bundle.Info, &bundle.ImageURL)
		if err != nil {
			return nil, NewAPIError(500, "Failed to load bundle from row", err)
		}

		paths, err := LoadPathsFromDatabase(transaction, bundle.Id)
		if err != nil {
			return nil, NewAPIError(500, "Failed to load bundle from row", err)
		}

		bundle.Paths = paths
		bundles = append(bundles, bundle)
	}

	err = rows.Err()
	if err != nil {
		return nil, NewAPIError(500, "Failed to load bundles", err)
	}

	return bundles, nil
}
