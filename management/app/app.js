'use strict';

angular.module('citadel', ['ngRoute', 'citadel.services', 'citadel.filters', 'citadel.directives'])
    .config(['$routeProvider',
        function ($routeProvider) {
            $routeProvider.when('/', {
                templateUrl: 'partials/dashboard.html',
                controller: 'DashboardController'
            });
            $routeProvider.when('/containers', {
                templateUrl: 'partials/containers.html',
                controller: 'ContainersController'
            });
            $routeProvider.when('/containers/:id', {
                templateUrl: 'partials/container.html',
                controller: 'ContainerController'
            });
            $routeProvider.otherwise({
                redirectTo: '/'
            });
    }]);