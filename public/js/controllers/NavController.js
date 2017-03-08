(function (angular) {
  'use strict';

  function NavController($scope, $rootScope, $location, $cookies, $http, $state, AuthService) {
    $scope.loginForm = {username:"", password:""}
    $scope.login = AuthService.login;

    $scope.login = function (loginForm) {
      AuthService.login(loginForm).then(loginSuccess, loginFailure)
    }

    function loginSuccess(response) {
      $scope.loginForm = {username: "", password:""};
      $state.go("admin")
    }

    function loginFailure(response) {
      $scope.loginForm = {username: "", password:""};
    }

    $scope.logout = function() {
      AuthService.logout().then(logoutSuccess, logoutFailure);
    }

    function logoutSuccess(response) {
      $state.go("map")
    }

    function logoutFailure(response) {
    }

    $scope.loggedIn = AuthService.loggedIn;
  }

  angular.module('appModules').controller('NavController', NavController);
})(angular);
