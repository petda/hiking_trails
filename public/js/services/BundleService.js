(function (angular) {
  'use strict';
  var URL = '/api/v1/bundles';

  function bundleService($rootScope, $http) {
    var service = {bundles: {}};

    service.Bundles = function () {
      return $http.get(URL, '').then(getBundlesFromResponse);
    };


    function getBundlesFromResponse(response) {
      var bundle;

      for (var i = 0; i < response.data.length; i++) {
        bundle = response.data[i];
        service.bundles[bundle.id] = bundle;
      }
    }

    service.createBundle = function (bundle) {
      return $http.post(URL, bundle).then(getBundleFromResponse);
    }

    service.updateBundle = function (bundle) {
      return $http.put(URL+"/"+bundle.id, bundle).then(getBundleFromResponse);
    }

    function getBundleFromResponse(response) {
      var currentBundle = service.bundles[response.data.id];

      if (currentBundle) {
	currentBundle.name = response.data.name;
	currentBundle.info = response.data.info;
	currentBundle.image = response.data.image;
      } else {
	service.bundles[response.data.id] = response.data;
      }
    }


    service.deleteBundle = function(bundle) {
      var removeBundleSuccessfull = function (response) {
        delete service.bundles[bundle.id];
	return bundle;
      };

      return $http.delete(URL+"/"+bundle.id).then(removeBundleSuccessfull);
    }

    return service;
  }

  angular.module('appModules').factory('BundleService', bundleService);
})(angular);
