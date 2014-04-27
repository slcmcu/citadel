function getChart(id) {
    return new Chart($(id).get(0).getContext('2d'));
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
    });
}

// temp function for mocking the hours in a day
function __generateLabels(count) {
    var out = [];
    for (var i = 1; i < count; i++) {
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