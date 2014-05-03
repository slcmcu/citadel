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

function newAreaChart(data, getValue, chartDiv, text) {
    var margin = {
        top: 20,
        right: 20,
        bottom: 30,
        left: 50
    },
        width = 768 - margin.left - margin.right,
        height = 250 - margin.top - margin.bottom;

    var parseDate = d3.time.format("%d-%b-%y").parse;

    var x = d3.time.scale()
        .range([0, width]);

    var y = d3.scale.linear()
        .range([height, 0]);

    var xAxis = d3.svg.axis()
        .scale(x)
        .orient("bottom");

    var yAxis = d3.svg.axis()
        .scale(y)
        .orient("left");

    var area = d3.svg.line()
        .x(function (d) {
            return x(d.time);
        })
        .y(function (d) {
            return y(getValue(d));
        });

    var svg = d3.select(chartDiv).append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
        .append("g")
        .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

    x.domain(d3.extent(data, function (d) {
        return d.time;
    }));
    y.domain([0, d3.max(data, getValue)]);

    svg.append("path")
        .datum(data)
        .attr("class", "line")
        .attr("d", area);

    svg.append("g")
        .attr("class", "x axis")
        .attr("transform", "translate(0," + height + ")")
        .call(xAxis);

    svg.append("g")
        .attr("class", "y axis")
        .call(yAxis)
        .append("text")
        .attr("transform", "rotate(-90)")
        .attr("y", 6)
        .attr("dy", ".71em")
        .style("text-anchor", "end")
        .text(text);
}