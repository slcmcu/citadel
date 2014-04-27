'use strict';

function HeaderController($scope) {
    $scope.template = 'partials/header.html';

    $scope.hosts = 202;
    $scope.runningContainers = 17338;
}

function DashboardController($scope) {
    $scope.template = 'partials/dashboard.html';

    $scope.cpuUsage = 48;
    $scope.memoryUsage = 65;
    $scope.diskUsage = 20;
}