(function (angular) {
  'use strict';

  function AdminController($scope, $rootScope, $http, $cookies, $state, AuthService, BundleService, PathService, PlaceService) {

    var markerClickListener;
    var markerDragListener;
    var polylineClickListener;
    var markers = {};
    var polylines = {};
    var temporaryMarker;
    var temporaryPolyline;
    var infoWindow = new google.maps.InfoWindow();
    var mapOptions = {
      zoom: 8,
      center: new google.maps.LatLng(63,20), // Center on location from example.
      mapTypeId: google.maps.MapTypeId.ROADMAP
    }

    var adminMap = new google.maps.Map(document.getElementById('adminMap'), mapOptions);

    $scope.bundles = BundleService.bundles;
    $scope.temporaryBundle = undefined;
    $scope.temporaryPath = undefined;
    $scope.temporaryPlace = undefined;
    $scope.viewMode = "listBundlesMode";


    function showBundleList() {
      $scope.viewMode = "listBundlesMode";
    }


    $scope.newBundle = function () {
      $scope.temporaryBundle = {name: "",
                                info: "",
                                image: "",
                                paths:[],
                               };
      $scope.viewMode = "addBundleMode";
    }


    $scope.editBundle = function (bundle) {
      $scope.temporaryBundle = {id: bundle.id,
                                name: bundle.name,
                                info: bundle.info,
                                image: bundle.image,
                                paths: [],
                               };
      $scope.viewMode = "editBundleMode";
    }


    $scope.cancelEditBundle = function () {
      $scope.temporaryBundle = undefined;
      $scope.viewMode = "listBundlesMode";
    }


    $scope.saveBundle = function () {
      if ($scope.temporaryBundle.id === undefined) {
        createBundle($scope.temporaryBundle);
      } else {
        updateBundle($scope.temporaryBundle);
      }
    }


    function createBundle(bundle) {
      BundleService.createBundle(bundle).then(createBundleSuccess, createBundleFailure);
    }


    function createBundleSuccess() {
      $scope.temporaryBundle = undefined;
      showInfo("Created bundle.");
      $scope.viewMode = "listBundlesMode";
    }


    function createBundleFailure(response) {
      handleError("Failed to create bundle: ", response);
    }


    function updateBundle(bundle) {
      BundleService.updateBundle(bundle).then(updateBundleSuccess, updateBundleFailure);
    }


    function updateBundleSuccess() {
      showInfo("Updated bundle.");
      $scope.viewMode = "listBundlesMode";
    }


    function updateBundleFailure(response) {
      handleError("Failed to update bundle: ", response)
    }


    $scope.deleteBundle = function (bundle) {
      BundleService.deleteBundle(bundle).then(deleteBundleSuccess, deleteBundleFailure);
    }


    function deleteBundleSuccess(deletedBundle) {
      for (var i = 0; i < deletedBundle.paths.length; i++){
	      var path = deletedBundle.paths[i];

	      removePathPolylineIfItExists(path);

	      for (var j = 0; j < path.places.length; j++){
	        removePlaceMarkerIfItExists(path.places[j]);
	      }
      }

      showInfo("Deleted bundle.");
      $scope.viewMode = "listBundlesMode";
    }


    function deleteBundleFailure(response) {
      handleError("Failed to delete bundle: ", response);
    }


    $scope.newPath = function (bundle) {
      $scope.temporaryPath = {
        name: "",
        info: "",
        image: "",
        length:"",
        polyline: [],
        duration: "",
        places: [],
        bundleId: bundle.id,
      }

      hidePathsAndPlaces();

      polylineClickListener = adminMap.addListener('click', function(event) {
        addLatLngPolyline(event, adminMap);
      });

      $scope.viewMode = "addPathMode";
    }


    $scope.editPath = function (path) {
      $scope.temporaryPath = {
        id: path.id,
        name: path.name,
        info: path.info,
        image: path.image,
        length: path.length,
        polyline: angular.copy(path.polyline),
        duration: path.duration,
        places: angular.copy(path.places),
        bundleId: path.bundleId,
      }

      hidePathsAndPlaces();

      temporaryPolyline = new google.maps.Polyline({
        map: adminMap,
        path: $scope.temporaryPath.polyline,
        strokeColor: '#000000',
        strokeOpacity: 1.0,
        strokeWeight: 3
      });

      zoomToFitAndCenterOnCoordinates(adminMap, temporaryPolyline.getPath().b)

      polylineClickListener = adminMap.addListener('click', function(event) {
        addLatLngPolyline(event, adminMap);
      });

      $scope.viewMode = "editPathMode";
    }


    $scope.savePath = function () {
      removeTemporaryPolylineIfItExists();

      showPathsAndPlaces();

      if ($scope.temporaryPath.id === undefined) {
        createPath($scope.temporaryPath)
      } else {
        updatePath($scope.temporaryPath)
      }
    }

    function createPath(path) {
      PathService.createPath(path).then(createPathSuccess, createPathFailure)
    }

    function createPathSuccess(createdPath) {
      var bundle = $scope.bundles[createdPath.bundleId];

      if (bundle === undefined) {
	      // Parent Bundle object deleted while waiting for request to return.
	      $scope.viewMode = "listBundlesMode";
	      return
      }

      bundle.paths.push(createdPath);
      addOrUpdatePolyline(createdPath);

      $scope.temporaryPath = undefined;
      $scope.viewMode = "listBundlesMode";
      showInfo("Created path.");
    }


    function createPathFailure(response) {
      handleError("Failed to create path.", response);
    }


    function updatePath(path) {
      PathService.updatePath(path).then(updatePathSuccess, updatePathFailure)
    }


    function updatePathSuccess(updatedPath) {
      var path = pathFromId(updatedPath.id);

      if (path === undefined) {
	      // Path object deleted while waiting for request to return.
	      $scope.viewMode = "listBundlesMode";
	      return
      }

      path.name = updatedPath.name;
      path.info = updatedPath.info;
      path.image = updatedPath.image;
      path.length = updatedPath.length;
      path.polyline = updatedPath.polyline;
      path.duration = updatedPath.duration;

      addOrUpdatePolyline(path);

      $scope.temporaryPath = undefined;
      $scope.viewMode = "listBundlesMode";
      showInfo("Updated path.");
    }


    function updatePathFailure(response) {
      handleError("Failed to update path: ", response)
    }


    function handleError(info, response) {
      showAlarm(info + ": " + response.status + " " + response.data.message);

      if (response.status === 401) {
        delete $cookies.SessionId;
        AuthService.logout();
        $state.go("map");
      }
    }


    $scope.deletePath = function (path) {
      PathService.deletePath(path).then(deletePathSuccess, deletePathFailure)
    }


    function deletePathSuccess(deletedPath) {
      for (var key in $scope.bundles) {
        if (!$scope.bundles.hasOwnProperty(key)) {
          continue
        }

        var paths = $scope.bundles[key].paths;

        for (var i = 0; i < paths.length; i++){
          if (paths[i].id === deletedPath.id) {
            paths.splice(i, 1)
          }
        }
      }

      var polyline = polylines[deletedPath.id];
      if (polyline !== undefined) {
	      polyline.setMap(null);
	      delete polylines[deletedPath.id];
      }

      for (var i = 0; i < path.places.length; i++){
	      removePlaceMarkerIfItExists(path.places[i]);
      }

      showInfo("Deleted path.");
      $scope.viewMode = "listBundlesMode";
    }


    function deletePathFailure(response) {
      handleError("Failed to delete path: ", response);
    }


    $scope.cancelEditPath = function () {
      removeTemporaryPolylineIfItExists();
      showPathsAndPlaces();
      $scope.temporaryPath = undefined;
      $scope.viewMode = "listBundlesMode";
    }


    $scope.newPlace = function (path) {
      $scope.temporaryPlace = {
        name: "",
        info: "",
        image: "",
        radius:1,
        position: {lat: 0, lng:0},
        media: [],
        pathId: path.id,
      }

      hidePathsAndPlaces();

      markerClickListener = adminMap.addListener('click', function(event) {
	      var latLng = {
	        lat: event.latLng.lat(),
	        lng: event.latLng.lng(),
	      };

        placeMarkerAndPanTo(latLng, adminMap);

	      $scope.temporaryPlace.position.lat = latLng.lat;
	      $scope.temporaryPlace.position.lng = latLng.lng;
	      $scope.$apply();

    	  // Should only be possible to add one marker.
    	  google.maps.event.removeListener(markerClickListener);
	      markerClickListener = undefined;
      });

      $scope.viewMode = "addPlaceMode";
    }


    $scope.editPlace = function (place) {
      $scope.temporaryPlace = {
        id: place.id,
        name: place.name,
        info: place.info,
        image: place.image,
        radius: place.radius,
        position: angular.copy(place.position),
        media: angular.copy(place.media),
        pathId: place.pathId,
      }

      hidePathsAndPlaces();

      placeMarkerAndPanTo($scope.temporaryPlace.position, adminMap);

      $scope.viewMode = "editPlaceMode";
    }


    $scope.cancelEditPlace = function () {
      removeTemporaryMarkerIfItExists();
      showPathsAndPlaces();
      $scope.temporaryPlace = undefined;
      $scope.viewMode = "listBundlesMode";
    }

    $scope.savePlace = function () {
      removeTemporaryMarkerIfItExists();
      showPathsAndPlaces();

      if ($scope.temporaryPlace.id === undefined) {
        createPlace($scope.temporaryPlace);
      } else {
        updatePlace($scope.temporaryPlace);
      }
    }


    function createPlace(place) {
      PlaceService.createPlace(place).then(createPlaceSuccess, createPlaceFailure);
    }


    function createPlaceSuccess(createdPlace) {
      var path = pathFromId(createdPlace.pathId);

      if (path === undefined) {
	      // Parent Path object deleted while waiting for request to return.
	      $scope.viewMode = "listBundlesMode";
	      return
      }

      path.places.push(createdPlace);
      addOrUpdatePlacemarker(createdPlace);

      $scope.temporaryPlace = undefined;
      showInfo("Created place.");
      $scope.viewMode = "listBundlesMode";
    }


    function createPlaceFailure(response) {
      handleError("Failed to create place.", response);
    }


    function updatePlace(place) {
      PlaceService.updatePlace(place).then(updatePlaceSuccess, updatePlaceFailure);
    }


    function updatePlaceSuccess(updatedPlace) {
      var place = placeFromId(updatedPlace.id);

      if (place === undefined) {
	      // Place object deleted while waiting for request to return.
	      $scope.viewMode = "listBundlesMode";
	      return
      }

      place.name = updatedPlace.name;
      place.info = updatedPlace.info;
      place.radius = updatedPlace.radius;
      place.position = updatedPlace.position;
      place.media = updatedPlace.media;

      addOrUpdatePlacemarker(place);

      $scope.temporaryPlace = undefined;

      showInfo("Updated place.");
      $scope.viewMode = "listBundlesMode";
    }


    function updatePlaceFailure(response) {
      handleError("Failed to update place: ", response);
    }


    function handleError(info, response) {
      showAlarm(info + ": " + response.status + " " + response.data.message);

      if (response.status === 401) {
        delete $cookies.SessionId;
        AuthService.logout();
        $state.go("map");
      }
    }


    $scope.deletePlace = function (place) {
      PlaceService.deletePlace(place).then(deletePlaceSuccess, deletePlaceFailure);
    }


    function deletePlaceSuccess(deletedPlace) {
      var paths;
      var places;

      for (var key in $scope.bundles) {
        if (!$scope.bundles.hasOwnProperty(key)) {
          continue
        }

        paths = $scope.bundles[key].paths
        for (var i = 0; i < paths.length; i++){

          places = paths[i].places;
          for (var j = 0; j < places.length; j++){
            if (places[j].id === deletedPlace.id) {
              places.splice(j, 1)
            }
    	    }
        }
      }

      var marker = markers[deletedPlace.id];
      if (marker !== undefined) {
	      marker.setMap(null);
	      delete markers[deletedPlace.id];
      }

      showInfo("Deleted place.");
      $scope.viewMode = "listBundlesMode";
    }


    function deletePlaceFailure(response) {
      handleError("Failed to delete place: ", response);
    }


    function addPlaceMarkers() {
      for (var id in $scope.bundles) {
        if (!$scope.bundles.hasOwnProperty(id)) {
          continue
        }

        var paths = $scope.bundles[id].paths

        for (var i = 0; i < paths.length; i++){
          var places = paths[i].places

          for (var j = 0; j < places.length; j++){
            addOrUpdatePlacemarker(places[j]);
          }
        }
      }
    }


    function addOrUpdatePlacemarker(place){
      removePlaceMarkerIfItExists(place);

      var marker = new google.maps.Marker({
        map: adminMap,
        position: new google.maps.LatLng(place.position.lat, place.position.lng),
        title: place.name
      });

      marker.content = '<div class="infoWindowContent">' + place.info +
        '<br/><br/><b>Latitude:</b> ' + place.position.lat + '<br/><b>Longitude:</b> ' +
        place.position.lng + '</div>';

      markers[place.id] = marker;

      google.maps.event.addListener(marker, 'click', function(){
        var content = '<h2>' + marker.title + '</h2>' + marker.content;
        infoWindow.setContent(content);

        // Position is set automatically to marker position if provided.
        infoWindow.open(adminMap, marker);
      });
    }


    function removePathPolylineIfItExists(path) {
      var polyline = polylines[path.id];

      if (polyline === undefined) {
	      return;
      }

      polyline.setMap(null);
      delete polylines[path.id];
    }


    function removePlaceMarkerIfItExists(place){
      var marker = markers[place.id];

      if (marker === undefined) {
	      return;
      }

      marker.setMap(null);
      delete markers[place.id];
    }


    function removeTemporaryMarkerIfItExists(){
      if (temporaryMarker !== undefined) {
	      temporaryMarker.setMap(null);
	      temporaryMarker = undefined;
      }

      if (markerClickListener !== undefined) {
	      google.maps.event.removeListener(markerClickListener);
	      markerClickListener = undefined;
      }

      if (markerDragListener !== undefined) {
	      google.maps.event.removeListener(markerDragListener);
	      markerDragListener = undefined;
      }
    }


    function removeTemporaryPolylineIfItExists(){
      if (temporaryPolyline !== undefined) {
	      temporaryPolyline.setMap(null);
	      temporaryPolyline = undefined;
      }

      if (polylineClickListener !== undefined) {
	      google.maps.event.removeListener(polylineClickListener);
	      polylineClickListener = undefined;
      }
    }


    function addPathPolylines() {
      for (var key in $scope.bundles) {
        if (!$scope.bundles.hasOwnProperty(key)) {
          continue
        }

        var paths = $scope.bundles[key].paths
        for (var i = 0; i < paths.length; i++){
          addOrUpdatePolyline(paths[i])
        }
      }
    }


    function addOrUpdatePolyline(path) {
      if (polylines[path.id] !== undefined) {
        polylines[path.id].setMap(null);
      }

      var polyline = new google.maps.Polyline({
        map: adminMap,
        path: path.polyline,
        geodesic: true,
        strokeColor: '#FF0000',
        strokeOpacity: 1.0,
        strokeWeight: 2,
        title: path.name
      });

      polyline.addListener('click', function(event){
        var polylineCoordinates = polyline.getPath().b
        var content = '<h2>' + path.name + '</h2>' + path.info + '<br/><br/><b>Duration:</b> ' + path.duration + ' Hours' + '<br/><b>Length:</b> ' + path.length + ' km';

        zoomToFitAndCenterOnCoordinates(adminMap, polylineCoordinates)

        infoWindow.setContent(content);
        infoWindow.setPosition(polylineCoordinates[0]);
        infoWindow.open(adminMap);
      });

      polylines[path.id] = polyline;
    }


    function addEditablePolyline(path) {
      if (path.polyline.length < 1) {
        return
      }

      path.mapPolyline = new google.maps.Polyline({
        map: adminMap,
        path: path.polyline,
        geodesic: true,
        strokeColor: '#FF0000',
        strokeOpacity: 1.0,
        strokeWeight: 2,
        title: path.name,
        editable: true,
        dragable: true
      });
    }


    function zoomToFitAndCenterOnCoordinates(map, coordinates) {
      var bounds = new google.maps.LatLngBounds();

      for (var i = 0; i < coordinates.length; i++) {
        bounds.extend(coordinates[i]);
      }

      map.fitBounds(bounds);
    }


    $scope.openInfoWindow = function(e, element){
      e.preventDefault();
      // NOTE element should be either marker or polyline.
      google.maps.event.trigger(element, 'click', {});
    }


    $scope.showPathOnMap = function(path){
      // There seems to be a bug in the google maps api, that require {} as argument,
      // or else an error is written to the console.
      google.maps.event.trigger(polylines[path.id], 'click', {});
    }


    $scope.showPlaceOnMap = function(place){
      google.maps.event.trigger(markers[place.id], 'click');
    }


    // Handles click events on a map, and adds a new point to the Polyline.
    function addLatLngPolyline(event, map) {
      if (temporaryPolyline === undefined) {
        temporaryPolyline = new google.maps.Polyline({
          map: map,
          strokeColor: '#000000',
          strokeOpacity: 1.0,
          strokeWeight: 3
        });
      }

      var polylinePath = temporaryPolyline.getPath();

      // Because path is an MVCArray, we can simply append a new coordinate
      // and it will automatically appear.
      polylinePath.push(event.latLng);
      $scope.temporaryPath.polyline.push(event.latLng)

      // Angular does not know about the google maps api event listener,
      // so we must trigger update of scope manually.
      $scope.$apply();
    }


    function placeMarkerAndPanTo(latLng, map) {
      removeTemporaryMarkerIfItExists();

      temporaryMarker = new google.maps.Marker({
        position: latLng,
        map: adminMap,
        editable: true,
        draggable: true,
      });

      infoWindow.setContent('<h2>' + "Drag me to change position" + '</h2>');
      infoWindow.open(adminMap, temporaryMarker);

      $scope.temporaryPlace.position.lat = latLng.lat;
      $scope.temporaryPlace.position.lng = latLng.lng;

      markerDragListener = temporaryMarker.addListener("dragend", markerDragend);
    }


    function hidePathsAndPlaces() {
      setMapOnAll(markers, null);
      setMapOnAll(polylines, null);
    }


    function showPathsAndPlaces() {
      setMapOnAll(markers, adminMap);
      setMapOnAll(polylines, adminMap);
    }


    function setMapOnAll(items, map) {
      for (var id in items) {
        items[id].setMap(map);
      }
    }


    function markerDragend(event){
      $scope.temporaryPlace.position.lat = event.latLng.lat();
      $scope.temporaryPlace.position.lng = event.latLng.lng();

      // Angular does not know about the google maps api event listener,
      // so we must trigger update of scope manually.
      $scope.$apply();
    }


    function showAlarm(message) {
      showMessage(message, "alarm");
    }


    function showWarning(message) {
      showMessage(message, "warning");
    }


    function showInfo(message) {
      showMessage(message, "info");
    }


    function showMessage(message, level) {
      var messageClass;

      switch (level) {
      case "info":
	      messageClass = "bg-success";
	      break;
      case "warning":
	      messageClass = "bg-warning";
	      break;
      case "alarm":
	      messageClass = "bg-danger";
	      break;
      default:
	      messageClass = "bg-danger";
      }

      // Setting length to 0 clears all previous messages.
      $rootScope.messages.length = 0;
      $rootScope.messages.push({message: message, messageClass: messageClass});
    }


    function pathFromId(pathId) {
      for (var id in $scope.bundles) {
	      var paths = $scope.bundles[id].paths;

	      for (var i = 0; i < paths.length; i++) {
	        if (paths[i].id === pathId) {
	          return paths[i];
	        }
	      }
      }

      return undefined;
    }


    function placeFromId(placeId) {
      for (var id in $scope.bundles) {
	      var paths = $scope.bundles[id].paths;

	      for (var i = 0; i < paths.length; i++) {
	        var places = paths[i].places;

	        for (var j = 0; j < places.length; j++) {
	          if (places[j].id === placeId) {
	            return places[j];
	          }
	        }
	      }
      }

      return undefined;
    }

    // Load bundles at instantiation of controller.
    BundleService.Bundles().then(addPlaceMarkers).then(addPathPolylines);
  }

  angular.module('appModules').controller('AdminController', AdminController);
})(angular);
