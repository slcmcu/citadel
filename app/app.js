'use strict';

angular.module('citadel', ['ngRoute', 'citadel.services', 'citadel.filters', 'citadel.directives'])
    .config(['$routeProvider',
        function ($routeProvider) {
            $routeProvider.when('/', {
                templateUrl: 'partials/dashboard.html',
                controller: 'DashboardController'
            });
            $routeProvider.when('/images', {
                templateUrl: 'partials/images.html',
                controller: 'ImagesController'
            });
            $routeProvider.otherwise({
                redirectTo: '/'
            });
    }]);