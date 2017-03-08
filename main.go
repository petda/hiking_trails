package main

import (
	"database/sql"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	_ "github.com/mattn/go-sqlite3"
	"hiking_trails/src/controllers"
	"hiking_trails/src/middleware"
	"hiking_trails/src/models"
	"log"
)

const (
	DATABASE_FILE = "hiking_trails.sqlite3"
)

func main() {
	// os.Remove(DATABASE_FILE)

	db, err := sql.Open("sqlite3", DATABASE_FILE)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	MustEnableForeignKeyChecks(db)

	MustCreateUsersDBTableIfNotExist(db)
	MustCreatePlacesDBTableIfNotExist(db)
	MustCreatePathsDBTableIfNotExist(db)
	MustCreateBundlesDBTableIfNotExist(db)

	models.MustCreateDefaultAdministratorIfMissing(db)

	app := martini.New()

	app.Use(martini.Logger())
	app.Use(martini.Recovery())
	app.Use(martini.Static("public"))
	app.Use(render.Renderer())

	app.Map(db)

	sessionStore := middleware.NewSessionStore()
	app.Map(sessionStore)
	app.Use(middleware.Sessions("store", sessionStore))

	router := martini.NewRouter()
	app.MapTo(router, (*martini.Routes)(nil))
	app.Action(router.Handle)

	router.Post("/api/v1/login", binding.Bind(models.LoginForm{}), controllers.UsersControllerLogin)
	router.Post("/api/v1/logout", controllers.UsersControllerLogout)

	router.Get("/api/v1/bundles", controllers.BundlesControllerList)
	router.Group("/api/v1/bundles", func(router martini.Router) {
		router.Post("", binding.Bind(models.Bundle{}), controllers.BundlesControllerCreate)
		router.Get("/:id", controllers.BundlesControllerRead)
		router.Put("/:id", binding.Bind(models.Bundle{}), controllers.BundlesControllerUpdate)
		router.Delete("/:id", controllers.BundlesControllerDelete)
	}, middleware.AdministratorRequired)

	router.Get("/api/v1/places", controllers.PlacesControllerList)
	router.Group("/api/v1/places", func(router martini.Router) {
		router.Post("", binding.Bind(models.Place{}), controllers.PlacesControllerCreate)
		router.Get("/:id", controllers.PlacesControllerRead)
		router.Put("/:id", binding.Bind(models.Place{}), controllers.PlacesControllerUpdate)
		router.Delete("/:id", controllers.PlacesControllerDelete)
	}, middleware.AdministratorRequired)

	router.Get("/api/v1/paths", controllers.PathsControllerList)
	router.Group("/api/v1/paths", func(router martini.Router) {
		router.Post("", binding.Bind(models.Path{}), controllers.PathsControllerCreate)
		router.Get("/:id", controllers.PathsControllerRead)
		router.Put("/:id", binding.Bind(models.Path{}), controllers.PathsControllerUpdate)
		router.Delete("/:id", controllers.PathsControllerDelete)
	}, middleware.AdministratorRequired)

	app.Run()
}

func MustCreateUsersDBTableIfNotExist(db *sql.DB) {
	sqlStatement := `
	CREATE TABLE IF NOT EXISTS users (id integer not null primary key,
                                    username VARCHAR(255) NOT NULL,
                                    salt VARCHAR(255),
                                    hashed_password VARCHAR(255),
                                    is_administrator BOOLEAN,
                                    CONSTRAINT username_unique UNIQUE (username));
 `

	_, err := db.Exec(sqlStatement)
	if err != nil {
		log.Fatalf("Failed to create 'users' database table: %s", err)
	}
}

func MustCreateBundlesDBTableIfNotExist(db *sql.DB) {
	sqlStatement := `
	CREATE TABLE IF NOT EXISTS bundles (id integer not null primary key,
                                      name VARCHAR(255),
                                      info VARCHAR(255),
                                      image_url VARCHAR(255));
 `

	_, err := db.Exec(sqlStatement)
	if err != nil {
		log.Fatalf("Failed to create 'bundles' database table: %s", err)
	}
}

func MustCreatePathsDBTableIfNotExist(db *sql.DB) {
	sqlStatement := `
	CREATE TABLE IF NOT EXISTS paths (id INTEGER NOT NULL PRIMARY KEY,
                                    name VARCHAR(255),
                                    info VARCHAR(255),
                                    length VARCHAR(255),
                                    duration VARCHAR(255),
                                    image_url VARCHAR(255),
                                    polyline BLOB,
                                    bundle_id INTEGER NOT NULL REFERENCES bundles(id) ON UPDATE CASCADE ON DELETE CASCADE);
 `

	_, err := db.Exec(sqlStatement)
	if err != nil {
		log.Fatalf("Failed to create 'paths' database table: %s", err)
	}
}

func MustCreatePlacesDBTableIfNotExist(db *sql.DB) {
	sqlStatement := `
	CREATE TABLE IF NOT EXISTS places (id integer not null primary key,
                                     name VARCHAR(255),
                                     info VARCHAR(255),
                                     radius BIGINT,
                                     position BLOB,
                                     path_id INTEGER NOT NULL REFERENCES paths(id) ON UPDATE CASCADE ON DELETE CASCADE);
 `

	_, err := db.Exec(sqlStatement)
	if err != nil {
		log.Fatalf("Failed to create 'places' database table: %s", err)
	}
}

func MustEnableForeignKeyChecks(db *sql.DB) {
	_, err := db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatalf("Failed to enable foreign key checks: %s", err)
	}
}
