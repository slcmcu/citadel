'use strict';

// Page header that displays the totals for the cluster
function HeaderController($scope, Hosts) {
    $scope.template = 'partials/header.html';

    Hosts.query({}, function (d) {
        $scope.hosts = d.length;
        var cpus = 0;
        var memory = 0;

        angular.forEach(d, function (v) {
            cpus += v.cpus;
            memory += v.memory;
        });

        $scope.cpus = cpus;
        $scope.memory = memory;
    });
}

// Dashboard includes overall information with graphs and services 
// for the cluster
function DashboardController($scope, Hosts) {
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

function HostsController($scope, $routeParams, Hosts) {
    $scope.template = 'partials/hosts.html';

    Hosts.query({
        name: $routeParams.id
    }, function (data) {
        $scope.hosts = data;
    });
}

function ContainersController($scope, $routeParams, Containers) {
    $scope.template = 'partials/containers.html';

    Containers.query({
        name: $routeParams.id
    }, function (data) {
        $scope.containers = data;
    });
}

// this needs to move to some super start init func
function toggleStartSidebar() {
    $('.ui.sidebar').sidebar({
        overlay: true
    })
        .sidebar('toggle');
}