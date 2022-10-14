var globaldata={};
globaldata.index=0;

ses=[[1,2,3,6],[4,7,6,8],[3,2,1,4],[2,1,3,6],[5,8,4,7],[2,1,6,3]]
tps=["A","","A","A","","A"]

select=Math.floor(Math.random()*ses.length)
ses2=ses[select]
tps2=tps[select]
console.log("ses:",ses2)

let title="原始目标"
ses2.forEach((a,i)=>{
    let url=`/webapp/ai_data/ct/ai_data/A${a}/json/`
    console.log("R",url)
    showResult(title,tps2,url);
    title="相似结果"
})


//
// showResult("原始目标","A","/webapp/ai_data/ct/ai_data/A1/json/");
// showResult("相似结果1","A","/webapp/ai_data/ct/ai_data/A2/json/");
// showResult("相似结果2","A","/webapp/ai_data/ct/ai_data/A3/json/");
// showResult("相似结果3","A","/webapp/ai_data/ct/ai_data/A5/json/");
// showResult("结果4","A","/webapp/ai_data/ct/ai_data/A5/json/");
// showResult("结果5","A","/webapp/ai_data/ct/ai_data/A6/json/");
// showResult("结果6","A","/webapp/ai_data/ct/ai_data/A7/json/");

/*
//showResult("F目标","B","/webapp/ai_data/ct/ai_data/B1/json/");
//showResult("F结果1","B","/webapp/ai_data/ct/ai_data/B2/json/");
showResult("目标","B","/webapp/ai_data/ct/ai_data/B3/json/");
showResult("结果1","B","/webapp/ai_data/ct/ai_data/B4/json/");
//showResult("F结果4","B","/webapp/ai_data/ct/ai_data/B5/json/");
showResult("结果2","B","/webapp/ai_data/ct/ai_data/B6/json/");
showResult("结果3","B","/webapp/ai_data/ct/ai_data/B7/json/");

 */

function showResult(title, type, path) {
    var index = globaldata.index;
    globaldata.index=globaldata.index+1;
    var id_3d = "result_3d_"+index;
    var id_chart1 = "result_c1_"+index;
    var id_chart1_legend = "result_c1l_"+index;
    var id_chart2 = "result_c2_"+index;
    var id_chart2_legend = "result_c2l_"+index;
    var id_image="result_img_"+index;

    var imgsrc = path+"../RESULT/jpgs/section_";

    globaldata[index]={};
    globaldata[index].id_image=id_image;
    globaldata[index].path=path;
    globaldata[index].imgsrc=imgsrc;

    var str = title + "（分型："+type+"）";

    var htmlcont = '<div class="row"><div class="card card-primary col-sm-12"><div class="card-header with-border">' +
        '<h3 class="card-title">' + str + '</h3>' +
        '<div class="card-tools"><button type="button" class="btn btn-tool" ai_data-card-widget="collapse"><i class="fas fa-minus"></i></button></div></div>' +
        '<div class="card-body chart-responsive"><div class="row" style="height:400px">' +
        // '<div class="col-sm-2"><div class="btn btn-info">转至影像详情</div><br /><div id="' + id_3d + '" style="width: 100%;height: 300px"></div></div>'+
        '<div class="col-sm-4"><img id="'+id_image+'" src="" height="300px" width="300px" style="margin-top: 50px"></div>'+
        '<div class="col-sm-4"><div id="' + id_chart1_legend + '" class="legend"></div><br /><div id="' + id_chart1 + '" style="height:360px"></div></div>' +
        '<div class="col-sm-4"><div id="' + id_chart2_legend + '" class="legend"></div><br /><div id="' + id_chart2 + '" style="height:360px"></div></div>' +
        '</div></div></div></div>';

    $("#mycontent").append($(htmlcont));

    imagechange(id_image,2,path);

    sec1(path,id_3d,id_chart1,id_chart1_legend,id_chart2,id_chart2_legend);

}


function sec1(url,id_3d,id_chart1,id_chart1l,id_chart2,id_chart2l,id_image,path) {
    //document.getElementById(idHeader+'type').text=type;
    //document.getElementById(idHeader+"img-raw").src=imgraw;
    let cta_data = {};
    cta_data.url = url;

    // show3D(id_3d, getJson(cta_data.url + "mask.json"), '#ff244e');

    cta_data.r_long = getJson(cta_data.url + "section_long.json");
    cta_data.r_short = getJson(cta_data.url + 'section_short.json');
    cta_data.sec_std = getJson(cta_data.url + 'section_std.json');
    let section_mean = getJson(cta_data.url + 'section_mean.json');

    let chart1A_data = [], chart1B_data = [], chart1C_data = [];
    let chart2A_data = [], chart2B_data = [];

    for (let i = 0; i < cta_data.r_long.length; i++) {
        chart1A_data.push([i, cta_data.r_long[i]]);
        chart1B_data.push([i, cta_data.r_short[i]]);
        chart1C_data.push([i, cta_data.r_long[i] - cta_data.r_short[i]]);

        chart2A_data.push([i, cta_data.sec_std[i]]);
        chart2B_data.push([i, section_mean[i]]);
    }

    var line_data1 = {data: chart1A_data, color: '#3c8dbc', label: "截面最长直径"};
    var line_data2 = {data: chart1B_data, color: '#00c0ef', label: "截面最短直径"};
    var line_data3 = {data: chart1C_data, color: 'red', label: "长短径差",};
    mychart(id_chart1,id_chart1l, [line_data1, line_data2, line_data3],50,id_image,path);

    var line_data1 = {data: chart2A_data, color: '#3c8dbc', label: "标准差"};
    var line_data2 = {data: chart2B_data, color: '#00c0ef', label: "平均值"};
    mychart(id_chart2,id_chart2l, [line_data1, line_data2],500,id_image,path);

}

function imagechange(id,index,path) {
    var imgsrc = path+"../RESULT/jpgs/section_";
    var selector="#"+id;
    var url= imgsrc + index+".jpg";
    $(selector).attr("src",url);
}