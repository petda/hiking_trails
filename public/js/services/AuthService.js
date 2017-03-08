(function (angular) {
  'use strict';

  function authService($rootScope, $http, $cookies) {
    var service = {sessionId: $cookies.SessionId};

    $rootScope.messages = [];

    service.login = function(loginForm) {
      var promise;

      promise = $http({
        url: 'api/v1/login',
        method: 'POST',
        data: "username=" + loginForm.username + "&password=" + loginForm.password,
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      });

      return promise.then(loginSuccess, loginFailure);
    }

    function loginSuccess(response) {
      service.sessionId = response.data.SessionId;
      $rootScope.messages.length = 0;
      $rootScope.messages.push({message: "Succesfully logged in", messageClass: "bg-success"});
    }

    function loginFailure(response) {
      var message = "Failed to log in: " + response.data.message
      $rootScope.messages.length = 0;
      $rootScope.messages.push({message: message, messageClass: "bg-danger"});
    }

    service.logout = function() {
      return $http.post("api/v1/logout").then(logoutSuccess, logoutFailure);
    }

    function logoutSuccess(response) {
      service.sessionId = "";
      $rootScope.messages.length = 0;
      $rootScope.messages.push({message: "Succesfully logged out", messageClass: "bg-success"});
    }

    function logoutFailure(response) {
      $rootScope.messages.length = 0;
      $rootScope.messages.push({message: "Failed to log", messageClass: "bg-danger"});
    }


    service.loggedIn = function () {
      return service.sessionId && service.sessionId !== "";
    }

    return service;
  }


  angular.module('appModules').factory('AuthService', authService);
})(angular);
