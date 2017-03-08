(function (angular) {
  'use strict';
  var URL = '/api/v1/places';

  function placeService($rootScope, $http) {
    var service = {};

    service.createPlace = function (place) {
      var data = _.pick(place, ["name",
                               "info",
                               "radius",
                               "position",
                               "media",
                               "pathId"]);
      console.log("Creating place: ", data);
      return $http.post(URL, data).then(getPlaceFromResponse);
    }

    service.updatePlace = function (place) {
      var data = _.pick(place, ["id",
                               "name",
                               "info",
                               "radius",
                               "position",
                               "media",
                               "pathId"]);
      return $http.put(URL+"/"+place.id, data).then(getPlaceFromResponse);
    }

    function getPlaceFromResponse(response) {
      return response.data;
    }

    service.deletePlace = function(place) {
      return $http.delete(URL+"/"+place.id).then(function () {
        return place;
      });
    }


    return service;
  }


  angular.module('appModules').factory('PlaceService', placeService);
})(angular);
