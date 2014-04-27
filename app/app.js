'use strict';

angular.module('citadel', ['ngRoute', 'citadel.services', 'citadel.filters'])
    .config(['$routeProvider',
        function ($routeProvider) {
            $routeProvider.when('/', {
                templateUrl: 'partials/dashboard.html',
                controller: 'DashboardController'
            });

            $routeProvider.otherwise({
                redirectTo: '/'
            });
    }]);