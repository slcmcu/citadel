'use strict';

function HeaderController($scope) {
    $scope.template = 'partials/header.html';

    $scope.hosts = 202;
    $scope.runningContainers = 17338;
    $scope.start = toggleStartSidebar;
}

function DashboardController($scope) {
    $scope.template = 'partials/dashboard.html';

    var labels = __generateLabels(25),
        cpuData = __generateRandomData(66, 100, 24),
        memoryData = __generateRandomData(61, 81, 24);

    newLineChart('#chart-cpu', labels, cpuData);
    newLineChart('#chart-memory', labels, memoryData);
}

function ServicesController($scope) {
    $scope.template = 'partials/services.html';

    $scope.services = [
        {
            name: 'api',
            ip: '192.168.56.101',
            status: 'healthy'
        },
        {
            name: 'scheduler',
            ip: '192.168.56.102',
            status: 'healthy'
        },
        {
            name: 'storage',
            ip: '192.168.56.102',
            status: 'healthy'
        },
        {
            name: 'metrics',
            ip: '192.168.56.103',
            status: 'sick'
        }
    ];
}

function StartController($scope) {
    $scope.template = 'partials/start.html';

    $scope.images = [
        'redis',
        'rethinkdb',
        'nsqd',
        'nsqadmin'
    ];
}

function toggleStartSidebar() {
    $('.ui.sidebar').sidebar({
        overlay: true
    })
        .sidebar('toggle');
}