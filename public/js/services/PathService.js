(function (angular) {
  'use strict';
  var URL = '/api/v1/paths';

  function pathService($rootScope, $http) {
    var service = {};

    service.createPath = function (path) {
      var data = _.pick(path, ["name",
                               "info",
                               "image",
                               "length",
                               "polyline",
                               "places",
                               "duration",
                               "bundleId"]);
      return $http.post(URL, data).then(getPathFromResponse);
    }

    service.updatePath = function (path) {
      var data = _.pick(path, ["id",
                               "name",
                               "info",
                               "image",
                               "length",
                               "polyline",
                               "places",
                               "duration",
                               "bundleId"]);
      return $http.put(URL+"/"+path.id, data).then(getPathFromResponse);
    }

    function getPathFromResponse(response) {
      return response.data;
    }

    service.deletePath = function(path) {
      return $http.delete(URL+"/"+path.id).then(function () {
        return path;
      });
    }


    return service;
  }


  angular.module('appModules').factory('PathService', pathService);
})(angular);
