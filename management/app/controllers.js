'use strict';

// Page header that displays the totals for the cluster
function HeaderController($scope, Services) {
    $scope.template = 'partials/header.html';

    Services.query({}, function (d) {
        $scope.services = d.length;
    });
}

// Dashboard includes overall information with graphs and services 
// for the cluster
function DashboardController($scope, Services) {
    /*
    Host.metrics({
        name: "b8f6b1166755"
    }, function (data) {
        var mem = function (d) {
            return (d.memory.used / d.memory.total) * 100;
        };
        var cpu = function (d) {
            return d.load_1;
        };

        newAreaChart(data, mem, '#chart-memory', 'mem %');
        newAreaChart(data, cpu, '#chart-cpu', 'load 1');
    });
    */
}

// Services display information about the cluster services that are running
function ServicesController($scope, $routeParams, Services) {
    $scope.template = 'partials/services.html';

    Services.query({
        name: $routeParams.id
    }, function (data) {
        $scope.services = data;
    });
}

// this needs to move to some super start init func
function toggleStartSidebar() {
    $('.ui.sidebar').sidebar({
        overlay: true
    })
        .sidebar('toggle');
}