'use strict';

angular.module('citadel.services', ['ngResource'])
    .factory('Hosts', function ($resource) {
        return $resource('/api/hosts/:name/:action', {}, {
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