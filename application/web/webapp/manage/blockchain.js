// 初始化卡片
class nodestatus {
    tps=0;
    nodes = 0;

    constructor() {
        this.nodelistRef = $("#nodelist")
    }

    nodeUpdate(name, height, ip, alive,count) {
        let status = $(`#${name}`)
        if (status.length === 0) {
            status = $("<div/>").attr("id", name).addClass("col-sm-4 mt-2")
            status.append(`<div class="position-relative p-3 bg-gray" style="height: 120px"><div class="ribbon-wrapper ribbon-lg"><div id="status" class="ribbon bg-success text-lg">正常</div></div><b><a>${name}</a></b><br /><br />H:<a id="height">${height}</a> IP:<a id="ip">${ip}</a>&nbsp;</div>`)
            this.nodelistRef.append(status)
        } else {
            if (alive) {
                $(`#${name} #status`).removeClass("bg-warning").addClass("bg-success").text("正常")
            } else {
                $(`#${name} #status`).removeClass("bg-success").addClass("bg-warning").text("异常")
            }

            $(`#${name} #height`).text(height)
            $(`#${name} #ip`).text(ip)
        }
    }

    updateStatus() {
        $.get("/api/v1/blockchain/nodelist", (resp) => {
            if (resp.code !== 200) {
                console.log("response error", resp.message)
            } else {
                let nodelist = resp.data
                if (nodelist.length !== this.nodes) {
                    this.nodelistRef.empty()
                    this.nodes = nodelist.length
                }

                for (let i = 0; i < this.nodes; i++) {
                    let n = nodelist[i]
                    this.nodeUpdate(n.Name, n.Height, n.IP, n.Alive,i)
                }
            }
        })
    }

    updateTps() {
        $.get("/api/v1/blockchain/tps",(resp)=>{
            if (resp.code !== 200) {
                console.log("response error",resp.message)
            } else {
                this.tps = resp.data
            }
        })
    }
}

let s = new nodestatus()

function bsUpdate() {
    s.updateStatus()
    s.updateTps()
}

// tps monitor
$(function () {
    let data = [[0, 0]], maxPoints = 100;
    const updateInterval = 1000;//Fetch data ever x milliseconds
    const updateBlockInterval = 1000;//Fetch data ever x milliseconds
    let realtime = 'on'; //If == to on then fetch data every x seconds. else stop fetching

    function update() {
        updateTps();
        let tps_plot = $.plot('#tps_plot', [
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
        if (data.length > maxPoints) {
            data = data.slice(1);
        }
        let n = data.length;
        let time = data[n - 1][0] + 1;
        let tmp = [time, s.tps];
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
