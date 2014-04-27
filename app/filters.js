'use strict';

angular.module('citadel.filters', [])
    .filter('servicestatus', function () {
        return function (status) {
            if (status === 'sick') {
                return 'error';
            }
            return 'positive';
        };
    });