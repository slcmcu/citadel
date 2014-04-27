'use strict';

function HeaderController($scope) {
    $scope.template = 'partials/header.html';

    $scope.hosts = 202;
    $scope.runningContainers = 17338;
}

function DashboardController($scope) {
    $scope.template = 'partials/dashboard.html';

    var labels = __generateDates();
    var cpuData = __generateRandomData(66, 100, 24);
    var memoryData = __generateRandomData(61, 81, 24);

    newLineChart('#chart-cpu', labels, cpuData);
    newLineChart('#chart-memory', labels, memoryData);
}