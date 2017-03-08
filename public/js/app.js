(function (angular) {
  'use strict';

  // angular.module('appModules', ['ngCookies', 'ui.router', 'ui.bootstrap'])
  angular.module('appModules', ['ngCookies', 'ui.router'])
    .config(function ($stateProvider, $urlRouterProvider, $locationProvider, $httpProvider) {

    // FIXME using html5 routes doesn't really work right now. This one is related to routes.js line 54
    // $locationProvider.html5Mode(true);

    $urlRouterProvider.otherwise('/404');

    $stateProvider.state('#', {
      url:            '',
      controller:     'MapController',
      templateUrl:    'html/partials/map.html'
    });

    $stateProvider.state('map', {
      url:            '/map',
      controller:     'MapController',
      templateUrl:    'html/partials/map.html'
    });

    $stateProvider.state('admin', {
      url:            '/admin',
      controller:     'AdminController',
      templateUrl:    'html/partials/admin.html'
    });

    $stateProvider.state('login', {
      url:            '/login',
      controller:     'LoginController',
      templateUrl:    'html/partials/login.html'
    });

    $stateProvider.state('404', {
      url:            '/404',
      templateUrl:    'html/partials/404.html'
    });

  })
  .run(function ($rootScope, $q, $location) {
    $rootScope.$on('$stateChangeStart', function (event, toState, toParams, fromState, fromParams) {
      $rootScope.error = null;
    });
  });

})(angular);
