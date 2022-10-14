function mychart(id,id_legend,data,ymax,id_image,path) {
    var legendOptions = {
        position: "ne",
        show: true,
        noColumns: 3,
        container: document.getElementById(id_legend),
    };

    var options = {
        legend: legendOptions,
        grid: {
            hoverable: true,
            borderColor: '#f3f3f3',
            borderWidth: 1,
            tickColor: '#f3f3f3'
        },
        series: {
            shadowSize: 1,
            lines: {
                show: true
            },
        },
        yaxes: [
            {
                min: 0,
                max: ymax,
                position: 'left',
                axisLabel: '测量长度（px）',
                show: true,
            }
        ],
        xaxes: [
            {
                position: 'bottom',
                axisLabel: '投影长度（px)',
                show: true,
            },
        ],
    };

    var plot = setupGraph('#' + id, data, options);

    //drawGraph(plot,data);

    function setupGraph(id, data, options) {
        //console.log("setupGraph", id, data, options);
        plot = $.plot($(id), data, options);
        $('<div class="tooltip-inner" id="line-chart-tooltip"></div>').css({
            position: 'absolute',
            display: 'none',
            opacity: 0.8
        }).appendTo('body');
        $(id).bind('plothover', function (event, pos, item) {
            if (item) {
                var index=this.id.split("_")[2];
                var d=globaldata[index];
                id_image=d.id_image;
                path=d.path;

                var x = item.datapoint[0].toFixed(2),
                    y = item.datapoint[1].toFixed(2)

                $('#line-chart-tooltip').html(item.series.label + ' of ' + x + ' = ' + y)
                    .css({
                        top: item.pageY + 5,
                        left: item.pageX + 5
                    })
                    .fadeIn(200);
                if (id_image) {
                    imagechange(id_image,Math.round(x),path);
                }
            } else {
                $('#line-chart-tooltip').hide();
            }
        });
        return plot;
    }

    function drawGraph(plot, data) {
        plot.setData(data);
        plot.setupGrid();
        plot.draw();
        requestAnimationFrame(drawGraph);
    }
}