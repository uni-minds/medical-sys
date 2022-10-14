// 初始化卡片

function bsUpdate() {
    $.get("/api/v1/blockchain/nodelist",(resp)=>{
        if (resp.code !== 200) {
            console.log("response error",resp.message)
        } else {
            let nodelist = resp.data
            for (let i=0;i<nodelist.length;i++) {
                let num = i + 1;
                let selector = `#center-1-node${num}-`;
                $(selector + "name").text(nodelist[i].Name);
                $(selector + "ip").text(nodelist[i].IP);
                $(selector + "height").text(nodelist[i].BlockHeight);
                $(selector + "status").text("正常");
            }
        }
    })

    $.get("/api/v1/blockchain/tps",(resp)=>{
        if (resp.code !== 200) {
            console.log("response error",resp.message)
        } else {
            TPS = resp.data
        }
    })

}

let TPS = 0

// tps monitor
$(function () {
    let data = [[0, 0]], maxPoints = 100;
    const updateInterval = 1000;//Fetch data ever x milliseconds
    const updateBlockInterval = 1000;//Fetch data ever x milliseconds
    let realtime = 'on'; //If == to on then fetch data every x seconds. else stop fetching


    function update() {
        updateTps();
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
            bsUpdate()
        }
    }

    function updateTps() {
        if (data.length>maxPoints) {
            data=data.slice(1);
        }
        let n=data.length;
        let time=data[n-1][0]+1;
        let tmp=[time,TPS];
        data.push(tmp);
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
});
