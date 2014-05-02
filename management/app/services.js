'use strict';

angular.module('citadel.services', ['ngResource'])
    .factory('Host', function ($resource) {
        return $resource('/api/hosts/:name/', {}, {
            query: {
                method: 'GET',
                isArray: true
            },
        });
    });