Hiking Trails
===========

A small and simple application written to test out the Google Maps API. It consists of a backend server written in Golang, which serves an API and a GUI built with AngularJS and Bootstrap.


Installation
------------

Follow instructions on the [official golang website](https://golang.org/dl/) to install Golang on your operating system. This project has only been tested with Golang 1.8 as of yet.

The backend server uses a Golang wrapper for sqlite3, to connect to an SQlite3 database. This wrapper uses CGO, and therefore some version of gcc needs to be installed on the system.

### Installing gcc on Linux

Install gcc if not already installed on system:

```
sudo apt get install build-essential
```

### Installing gcc on OS X

If you dont have gcc installed on your system then run `xcode-select --install` and follow the guide.

### Installing gcc on Windows

In order to use cgo on Windows, you'll need to install a gcc compiler (for instance, [mingw-w64](http://mingw-w64.org/doku.php)) and have gcc.exe (etc.) in your PATH environment variable.


### Building application

Golang expects the go files to be located under `$GOPATH/src/hiking_trails`.
If this repository has been checked out to some other path, then move the files to the
correct directory first, then change to that directory.

```
mv <current_path> $GOPATH/src/hiking_trails
cd $GOPATH/src/hiking_trails
```

Fetch application go dependencies:

```
go get -v
```

Build application:
```
go build -o hiking_trails
```

Run
------------

Start by changing the string GOOGLE_MAPS_API_KEY in public/index.html to your Google Maps API key. If you don't have an API key, you can get one from [here](https://developers.google.com/maps/documentation/javascript/get-api-key).

To run the webserver simply run `./hiking_trails` in a terminal. Then the GUI will be accesible from a web browser at `localhost:3000`.


Dependencies
------------

### Backend(Golang)
* [Martini](http://github.com/go-martini/martini): Lightweight and easy framework for writing web servers in Golang.
* [Martini binding](http://github.com/martini-contrib/binding): Martini middleware/handler for parsing, validating and binding request data.
* [Martini render](http://github.com/martini-contrib/render): Martini middleware/handler for easily rendering serialized JSON and HTML templates.
* [go-sqlite](http://github.com/mattn/go-sqlite3): Golang sqlite driver.

### Frontend(Javascript/CSS)
* [AngularJS](https://angularjs.org/): Web framework for quicker development of web pages.
* [Bootstrap](http://getbootstrap.com/): Very common CSS framework(includes some javascript).
* [AdminLTE](https://almsaeedstudio.com/themes/AdminLTE/index2.html): Bootstrap CSS theme.
* [Jquery](https://jquery.com/): Javascript utility library(required by bootstrap javascript).
* [Lodash](https://lodash.com/): Javascript utility library.
* [Google maps Javascript API](https://developers.google.com/maps/documentation/javascript/tutorial):

Example API usage(using Curl)
------------


### Get all bundles

```
curl -v http://localhost:3000/api/v1/bundles
```

### Login

```
curl -v -H "Content-Type: application/x-www-form-urlencoded" -X POST "http://localhost:3000/api/v1/login?username=admin&password=admin"
```

Write down the returned session id since it is needed for all API requests that require authentication.


### Create bundle

Create a new file named new_bundle.json, containing the bundle to create in JSON format, then run:
```
curl -v -X POST -d @new_bundle.json --cookie "SessionId=0edb605e2acfbd1de35ad3a14052d2b5375795143cfaf5df000eba5be19b6c8e" http://localhost:3000/api/v1/bundles/1
```

### Update bundle

Copy the new_bundle.json file to a new file named updated_bundle.json. Add an extra field to the file with the id from the previously created bundle. Update some stuff then run:

```
curl -v -X PUT -d @updated_bundle.json --cookie "SessionId=0edb605e2acfbd1de35ad3a14052d2b5375795143cfaf5df000eba5be19b6c8e" http://localhost:3000/api/v1/bundles/1
```

### Delete bundle

```
curl -v -X DELETE --cookie "SessionId=0edb605e2acfbd1de35ad3a14052d2b5375795143cfaf5df000eba5be19b6c8e" http://localhost:3000/api/v1/bundles/1
```

### Logout

```
curl -v -X POST --cookie "SessionId=0edb605e2acfbd1de35ad3a14052d2b5375795143cfaf5df000eba5be19b6c8e" http://localhost:3000/api/v1/logout
```

