// 初始化卡片
$(function() {
    $("#center-1-group-name").text("昌平中心");
    $("#center-2-group-name").text("天坛医院");
    $("#center-3-group-name").text("朝阳医院");
});

var NNData={
    Height:0,
    Tps:0,
    NodeStatus:[false,false,false,false,false,false],
    NodeHeight:[0,0,0,0,0,0],
    GatewayStatus:[true,true,true]

};
// tps monitor
$(function () {
    var data = [[0,0]],maxPoints = 100;
    var updateInterval = 1000;//Fetch data ever x milliseconds
    var updateBlockInterval = 1000;//Fetch data ever x milliseconds
    var realtime = 'on'; //If == to on then fetch data every x seconds. else stop fetching


    function update() {
        updateTps();
        updateHeight();
        var tps_plot = $.plot('#tps_plot', [
                {
                    data: data,
                }
            ],
            {
                grid: {
                    borderColor: '#f3f3f3',
                    borderWidth: 1,
                    tickColor: '#f3f3f3'
                },
                series: {
                    color: '#3c8dbc',
                    lines: {
                        lineWidth: 2,
                        show: true,
                        fill: true,
                    },
                },
                yaxis: {
                    min: 0,
                    max: 1500,
                    show: true
                },
                xaxis: {
                    show: true
                }
            }
        );

        // Since the axes don't change, we don't need to call plot.setupGrid()
        tps_plot.draw();
        if (realtime === 'on') {
            setTimeout(update, updateInterval);
            bsStatusUpdate()
        }
    }

    function updateTps() {
        if (data.length>maxPoints) {
            data=data.slice(1);
        }
        let n=data.length;
        let tps=NNData.Tps;
        let time=data[n-1][0]+1;
        let tmp=[time,tps];
        data.push(tmp);
    }

    function updateHeight() {
        bcNodes.forEach(function (v, i) {
            let num = i + 1;
            let group = Math.ceil(num / 2);
            let nodenum = Math.floor(num / group);
            let selector = "#center-" + group + "-node" + nodenum + "-";

            $(selector + "load").text((NNData.Tps / 18000).toFixed(1) + "%");
            $(selector + "height").text(NNData.Height);
        });
    }

    if (realtime === 'on') {
        update();
    }

    $('#realtime .btn').click(function () {
        if ($(this).data('toggle') === 'on') {
            realtime = 'on';
        } else {
            realtime = 'off';
        }
        update();
    });

    // NNFunction
    function NNUpdateNNdata() {
        NNData=getJson("/monitor/bcstatus");
        setTimeout(NNUpdateNNdata, updateBlockInterval);
    }

    setTimeout(NNUpdateNNdata, updateBlockInterval);

});

var bcGateways=["172.16.1.151","172.16.1.158","172.16.1.159"];
var bcNodes=["172.16.1.152","172.16.1.153","172.16.1.154","172.16.1.155","172.16.1.156","172.16.1.157"];
function bsStatusUpdate() {
    bcGateways.forEach(function (v, i) {
        let num = i + 1;
        let selector = "#center-" + num + "-gateway-";

        $(selector + "name").text("Gateway#" + num);
        $(selector + "ip").text(v);
        $(selector + "status").text("正常");

    });

    bcNodes.forEach(function (v, i) {
        let num = i + 1;
        let group = Math.ceil(num / 2);
        let nodenum = Math.floor(num / group);
        let selector = "#center-" + group + "-node" + nodenum + "-";

        $(selector + "name").text("Node#" + num);
        $(selector + "ip").text(v);
        if (NNData.NodeStatus[i]) {
            $(selector + "status").text("正常");
        } else {
            $(selector + "status").text("离线");
        }
        $(selector + "height").text(NNData.NodeHeight[i]);

    });
}

bsStatusUpdate();