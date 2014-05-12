'use strict';

angular.module('citadel.services', ['ngResource'])
    .factory('Host', function ($resource) {
        return $resource('/api/nodes/:name/:action', {}, {
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
    })
    .factory('Container', function ($resource) {
        return $resource('/api/containers/:name/:action', {}, {
            query: {
                method: 'GET',
                isArray: true
            }
        });
    });
