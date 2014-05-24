'use strict';

angular.module('citadel.services', ['ngResource'])
    .factory('Services', function ($resource) {
        return $resource('/api/services/:name/:action', {}, {
            query: {
                method: 'GET',
                isArray: true,
                params: {
                    name: "@name"
                }
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