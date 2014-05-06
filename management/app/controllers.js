'use strict';

// Page header that displays the totals for the cluster
function HeaderController($scope, Host, Container) {
    $scope.template = 'partials/header.html';

    Host.query({}, function (d) {
        $scope.hosts = d.length;
    });

    Container.query({}, function (d) {
        $scope.runningContainers = d.length;
    });
}

// Dashboard includes overall information with graphs and services 
// for the cluster
function DashboardController($scope, Host) {
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
}

// Containers controller aggregates the container running 
// information for all containers running on the cluster
function ContainersController($scope, Container) {
    $scope.namespace = 'crosbymichael';

    Container.query({}, function (d) {
        $scope.containers = d;

        var mapReduce = function (get) {
            if ($scope.containers.length === 0) {
                return 0;
            }
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

    });

    $scope.addImage = function () {
        $('#image-modal')
            .modal('setting', {
                onApprove: function () {
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
                },
            })
            .modal('show');
    };
}

function ContainerController($scope, $routeParams) {
    $scope.data = [
        {
            key: 'Instances',
            value: 16
        },
        {
            key: 'Cpu',
            value: 17.3 + ' %'
        },
        {
            key: 'Memory',
            value: 16 * 2048 + ' mb'
        },
        {
            key: 'Network Transfer',
            value: 23 * 1024 + ' mb'
        },
        {
            key: 'Disk',
            value: 16 * 23283 + ' mb'
        },
        {
            key: 'On Hosts',
            value: 16
        },
        {
            key: 'Average Run Time',
            value: 32 + ' days'
        },
    ];

    $scope.instances = [

    ];
}

// Hosts display information about the cluster services that are running 
// on the hosts
function HostsController($scope, Host) {
    $scope.template = 'partials/hosts.html';

    Host.query({}, function (data) {
        $scope.hosts = data;
    });
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
}

// this needs to move to some super start init func
function toggleStartSidebar() {
    $('.ui.sidebar').sidebar({
        overlay: true
    })
        .sidebar('toggle');
}