var globaldata={};
globaldata.index=0;

let DB_CCTA="http://192.168.2.102:8000";
//let DB_CCTA="http://192.168.1.11:8080";
let url=window.location.href;

let StudiesUID=getQueryString(url,"StudiesUID");
let SeriesUID=getQueryString(url,"SeriesUID");
let ObjectUID=getQueryString(url,"ObjectUID");

console.log(ObjectUID);

showResult("目标",StudiesUID,SeriesUID,ObjectUID);
let head=ObjectUID.slice(0,53);
let foot=ObjectUID.slice(53);
showResult("检索结果1",StudiesUID,SeriesUID,head+(parseInt(foot)+2));
showResult("检索结果2",StudiesUID,SeriesUID,head+(parseInt(foot)+17));
showResult("检索结果3",StudiesUID,SeriesUID,head+(parseInt(foot)+4));


function getQueryString(url,key) {
    var reg = new RegExp("(^|&)" + key + "=([^&]*)(&|$)", "i");
    var r = url.match(reg);
    if ( r != null ){
        return unescape(r[2]);
    }else{
        return null;
    }
}

function showResult(title, StudiesUID,SeriesUID,ObjectUID) {

    let dcm4cheeWado = DB_CCTA+"/dcm4chee-arc/aets/AS_RECEIVED/wado?requestType=WADO";
    let wadoAddURL = "&studyUID=" + StudiesUID + "&seriesUID=" + SeriesUID + "&objectUID=" + ObjectUID;
    let jpegParams = '&contentType=image/jpeg&frameNumber=1';
    let u= dcm4cheeWado + wadoAddURL + jpegParams;
    console.log(u);

    var str = title;

    var htmlcont = '<div class="row"><div class="card card-primary col-sm-12"><div class="card-header with-border">' +
        '<h3 class="card-title">' + str + '</h3>' +
        '<div class="card-tools"><button type="button" class="btn btn-tool" data-card-widget="collapse"><i class="fas fa-minus"></i></button></div></div>' +
        '<div class="card-body chart-responsive"><div class="row" style="height:540px">' +
        '<img src="'+u+'" height="512px" width="512px">'+
        '</div></div></div></div>';

    $("#mycontent").append($(htmlcont));

    //sec1(path,id_3d,id_chart1,id_chart1_legend,id_chart2,id_chart2_legend);

}