function getChart(id) {
    var ctx = $(id).get(0).getContext('2d');
    return new Chart(ctx);
}

function newLineChart(id, labels, data) {
    var chart = getChart(id);
    var dataset = {
        fillColor: "rgba(151,187,205,0.5)",
        strokeColor: "rgba(151,187,205,1)",
        pointColor: "rgba(151,187,205,1)",
        pointStrokeColor: "#fff",
        data: data
    };
    chart.Line({
        labels: labels,
        datasets: [dataset]
    }, {
        scaleStepWidth: 1,
        pointDotRadius: 1,
        pointDot: false
        //scaleOverride: true
        // scaleSteps: labels.length
    });
}

// temp function for mocking the hours in a day
function __generateDates() {
    var out = [];
    for (var i = 1; i < 25; i++) {
        out[i - 1] = i.toString();
    }
    return out;
}

function __generateRandomData(min, max, count) {
    var out = [];
    for (var i = 1; i <= count; i++) {
        var v = Math.floor(Math.random() * (max - min + 1) + min);
        out.push(v);
    }
    return out;
}