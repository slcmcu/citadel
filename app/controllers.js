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

    var maxHosts = 202;

    $scope.scale = 1;
    $scope.cpuProfile = "low";
    $scope.memory = 52;
    $scope.image = '';
    $scope.cores = 1;

    $scope.images = [
        'redis',
        'rethinkdb',
        'nsqd',
        'nsqadmin'
    ];

    $scope.coresUp = function () {
        $scope.cores++;
    };
    $scope.coresDown = function () {
        if ($scope.cores > 1) {
            $scope.cores--;
        }
    };
    $scope.scaleUp = function () {
        if ($scope.scale < maxHosts) {
            $scope.scale++;
        }
    };
    $scope.scaleDown = function () {
        if ($scope.scale > 1) {
            $scope.scale--;
        }
    };

    $scope.launchContainers = function () {
        // have to get the value this way because binding to a checkbox
        // is near impossible or I am doing something wrong
        var isService = $('#is-service-checkbox')[0].checked
        console.log(isService);
        console.log($scope.cpuProfile);
    };
}

// this needs to move to some super start init func
function toggleStartSidebar() {
    $('.ui.sidebar').sidebar({
        overlay: true
    })
        .sidebar('toggle');

    $('.ui.dropdown')
        .dropdown();

    $('.ui.checkbox')
        .checkbox();

    $('.ui.radio.checkbox')
        .checkbox();
}