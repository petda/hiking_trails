(function (angular) {
  'use strict';

  function MapController($scope, BundleService) {
    var markers = {};
    var polylines = {};
    var infoWindow = new google.maps.InfoWindow();
    var mapOptions = {
      zoom: 8,
      center: new google.maps.LatLng(63,20),
      mapTypeId: google.maps.MapTypeId.ROADMAP
    }

    var nonAdminMap = new google.maps.Map(document.getElementById('map'), mapOptions);

    $scope.bundles = BundleService.bundles;


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



    function addPlaceMarkers() {
      for (var id in $scope.bundles) {
        if (!$scope.bundles.hasOwnProperty(id)) {
          continue
        }

        var paths = $scope.bundles[id].paths

        for (var i = 0; i < paths.length; i++){
          var places = paths[i].places

          for (var j = 0; j < places.length; j++){
            addPlaceMarker(places[j]);
          }
        }
      }
    }


    function addPlaceMarker(place){
      removePlaceMarkerIfItExists(place);

      var marker = new google.maps.Marker({
        map: nonAdminMap,
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
        infoWindow.open(nonAdminMap, marker);
      });
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
        map: nonAdminMap,
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

        zoomToFitAndCenterOnCoordinates(nonAdminMap, polylineCoordinates)

        infoWindow.setContent(content);
        infoWindow.setPosition(polylineCoordinates[0]);
        infoWindow.open(nonAdminMap);
      });

      polylines[path.id] = polyline;
    }

    // function addMarkers() {
    //   for (var key in $scope.bundles) {
    //     if (!$scope.bundles.hasOwnProperty(key)) {
    //       continue
    //     }

    //     var paths = $scope.bundles[key].paths

    //     for (var i = 0; i < paths.length; i++){
    //       var places = paths[i].places

    //       for (var j = 0; j < places.length; j++){
    //         addMarker(places[i]);
    //       }
    //     }
    //   }
    // }

    // function addMarker(place){
    //   var marker = new google.maps.Marker({
    //     map: nonAdminMap,
    //     position: new google.maps.LatLng(place.position.lat, place.position.lng),
    //     title: place.name
    //   });

    //   marker.content = '<div class="infoWindowContent">' + place.info +
    //     '<br/><br/>Latitude: ' + place.position.lat + ' Longitude:' +
    //     place.position.lng + '</div>';

    //   $scope.markers.push(marker);

    //   google.maps.event.addListener(marker, 'click', function(){
    //     var content = '<h2>' + marker.title + '</h2>' + marker.content;
    //     infoWindow.setContent(content);

    //     // Position is set automatically to marker position if provided.
    //     infoWindow.open(nonAdminMap, marker);
    //   });
    // }

    // function addPolylines() {
    //   for (var key in $scope.bundles) {
    //     if (!$scope.bundles.hasOwnProperty(key)) {
    //       continue
    //     }

    //     var paths = $scope.bundles[key].paths
    //     for (var i = 0; i < paths.length; i++){
    //       addPolyline(paths[i])
    //     }
    //   }
    // }

    // function addPolyline(path) {
    //   if (path.polyline.length < 1) {
    //     return
    //   }

    //   var polyline = new google.maps.Polyline({
    //     map: nonAdminMap,
    //     path: path.polyline,
    //     geodesic: true,
    //     strokeColor: '#FF0000',
    //     strokeOpacity: 1.0,
    //     strokeWeight: 2,
    //     title: path.name
    //   });

    //   polyline.addListener('click', function(event){
    //     var polylineCoordinates = polyline.getPath().b
    //     var content = '<h2>' + path.name + '</h2>' + path.info;

    //     zoomToFitAndCenterOnCoordinates(nonAdminMap, polylineCoordinates)

    //     infoWindow.setContent(content);
    //     infoWindow.setPosition(polylineCoordinates[0]);
    //     infoWindow.open(nonAdminMap);
    //   });

    //   $scope.polylines.push(polyline);
    // }

    function zoomToFitAndCenterOnCoordinates(map, coordinates) {
      var bounds = new google.maps.LatLngBounds();

      for (var i = 0; i < coordinates.length; i++) {
        bounds.extend(coordinates[i]);
      }

      map.fitBounds(bounds);
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


    // Load bundles at instantiation of controller.
    BundleService.Bundles().then(addPlaceMarkers).then(addPathPolylines);
  }

  angular.module('appModules').controller('MapController', MapController);
})(angular);
