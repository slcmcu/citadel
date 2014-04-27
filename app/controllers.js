'use strict';

// Page header that displays the totals for the cluster
function HeaderController($scope) {
    $scope.template = 'partials/header.html';

    $scope.hosts = 202;
    $scope.runningContainers = 17338;
    $scope.start = toggleStartSidebar;
}

// Dashboard includes overall information with graphs and services 
// for the cluster
function DashboardController($scope) {
    var labels = __generateLabels(25),
        cpuData = __generateRandomData(66, 100, 24),
        memoryData = __generateRandomData(61, 81, 24);

    newLineChart('#chart-cpu', labels, cpuData);
    newLineChart('#chart-memory', labels, memoryData);
}

// Containers controller aggregates the container running 
// information for all containers running on the cluster
function ContainersController($scope) {
    $scope.namespace = 'crosbymichael';

    $scope.containers = [
        {
            name: 'rethinkdb',
            instances: 16,
            cpu: 17,
            memory: 16 * 2048,
            status: 'healthy'
        },
        {
            name: 'redis',
            instances: 96,
            cpu: 7,
            memory: 96 * 512,
            status: 'healthy'
        },
        {
            name: 'sentry',
            instances: 16,
            cpu: 20,
            memory: 16 * 68,
            status: 'healthy'
        },
        {
            name: 'postgres',
            instances: 8,
            cpu: 7,
            memory: 8 * 243,
            status: 'healthy'
        },
        {
            name: 'mysql',
            instances: 1,
            cpu: 7,
            memory: 643,
            status: 'healthy'
        },
        {
            name: 'phabricator',
            instances: 1,
            cpu: 13,
            memory: 516,
            status: 'healthy'
        },
        {
            name: 'nsqd',
            instances: 100,
            cpu: 9,
            memory: 100 * 26,
            status: 'healthy'
        },
        {
            name: 'nsqadmin',
            instances: 1,
            cpu: 4,
            memory: 20,
            status: 'sick'
        }
    ];

    var mapReduce = function (get) {
        return $scope.containers.map(get).reduce(function (prev, curr, i, array) {
            return prev + curr;
        });
    };

    $scope.count = $scope.containers.length;
    $scope.totalMemory = mapReduce(function (i) {
        return i.memory;
    });
    $scope.totalCpu = mapReduce(function (i) {
        return i.cpu;
    });
    $scope.instances = mapReduce(function (i) {
        return i.instances;
    });
}

// Services display information about the cluster services that are running 
// on the hosts
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

// New Image allows as user to input the image name and pulls the image
// down to all hosts in the cluster
function NewImageController($scope) {
    $scope.template = 'partials/new-image.html';

    $scope.deployName = '';
    $scope.deployImage = function () {
        var progress = $('#deploy-progress');
        var i = 0;
        var update = function () {
            if (i < 100) {
                progress.css('width', i + '%');
                i = i + (100 / 202);
                setTimeout(update, 50);
            }
        };
        update();
    };
}

// this needs to move to some super start init func
function toggleStartSidebar() {
    $('.ui.sidebar').sidebar({
        overlay: true
    })
        .sidebar('toggle');
}