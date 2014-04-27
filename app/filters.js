'use strict';

angular.module('citadel.filters', [])
    .filter('servicestatus', function () {
        return function (status) {
            switch (status.toLowerCase()) {
            case 'sick':
                return 'error';
            case 'stale':
                return 'warning';
            default:
                return 'positive';
            }
        };
    });