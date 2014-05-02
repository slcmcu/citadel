'use strict';

angular.module('citadel.services', ['ngResource'])
    .factory('Host', function ($resource) {
        return $resource('/api/hosts/:name/:action', {}, {
            query: {
                method: 'GET',
                isArray: true
            },
            memory: {
                method: 'GET',
                isArray: true,
                params: {
                    action: 'memory'
                }
            },
        });
    });