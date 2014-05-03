'use strict';

angular.module('citadel.services', ['ngResource'])
    .factory('Host', function ($resource) {
        return $resource('/api/hosts/:name/:action', {}, {
            query: {
                method: 'GET',
                isArray: true
            },
            metrics: {
                method: 'GET',
                isArray: true,
                params: {
                    action: 'metrics',
                    name: "@name"
                }
            },
        });
    });