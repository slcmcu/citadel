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
    })
    .filter('tomb', function () {
        return function (bytes) {
            return bytes / 1024 / 1024;
        };
    });