console.log("INFO: Processing aorta infos.");

let baseurl="../../webapp/ai_data/output/ccta/json/";
let coronaryData=[];

console.log(baseurl+'lad_max.json');
coronaryData["lad_max"]=getJson(baseurl+'lad_max.json');
coronaryData["lad_min"]=getJson(baseurl+'lad_min.json');


coronaryData["lcx_max"]=getJson(baseurl+'lcx_max.json');
coronaryData["lcx_min"]=getJson(baseurl+'lcx_min.json');


coronaryData["lm_max"]=getJson(baseurl+'lm_max.json');
coronaryData["lm_min"]=getJson(baseurl+'lm_min.json');


coronaryData["rca_max"]=getJson(baseurl+'rca_max.json');
coronaryData["rca_min"]=getJson(baseurl+'rca_min.json');


function drawChart(id,dataMax,dataMin,descrMax,descrMin) {
    let xlabel=['Start'];
    for (i=1;i<dataMax.length-1;i++) {
        xlabel.push('');
    }
    xlabel.push('End');
    var chartData= {
        labels  : xlabel,
        datasets: [
            {
                label               : descrMax,
                fillColor           : 'rgb(0,13,255)',
                strokeColor         : 'rgb(255,0,6)',
                pointColor          : 'rgba(210, 214, 222, 1)',
                pointStrokeColor    : '#c1c7d1',
                pointHighlightFill  : '#ff0006',
                pointHighlightStroke: 'rgba(220,220,220,1)',
                data                : dataMax
            },
            {
                label               : descrMin,
                fillColor           : 'rgba(60,141,188,0.9)',
                strokeColor         : 'rgba(43,255,41,0.8)',
                pointColor          : '#3b8bba',
                pointStrokeColor    : 'rgba(60,141,188,1)',
                pointHighlightFill  : '#fff',
                pointHighlightStroke: 'rgba(60,141,188,1)',
                data                : dataMin
            }
        ]
    };
    var lineChartCanvas          = $('#'+id).get(0).getContext('2d');
    var lineChart                = new Chart(lineChartCanvas);
    var lineChartOptions         = getChartOptions();
    lineChart.Line(chartData, lineChartOptions);
}




function getChartOptions() {
    var lineChartOptions = {
        //Boolean - If we should show the scale at all
        showScale               : true,
        //Boolean - Whether grid lines are shown across the chart
        scaleShowGridLines      : false,
        //String - Colour of the grid lines
        scaleGridLineColor      : 'rgba(0,0,0,.05)',
        //Number - Width of the grid lines
        scaleGridLineWidth      : 10,
        //Boolean - Whether to show horizontal lines (except X axis)
        scaleShowHorizontalLines: false,
        //Boolean - Whether to show vertical lines (except Y axis)
        scaleShowVerticalLines  : true,
        //Boolean - Whether the line is curved between points
        bezierCurve             : true,
        //Number - Tension of the bezier curve between points
        bezierCurveTension      : 0.3,
        //Boolean - Whether to show a dot for each point
        pointDot                : false,
        //Number - Radius of each point dot in pixels
        pointDotRadius          : 4,
        //Number - Pixel width of point dot stroke
        pointDotStrokeWidth     : 1,
        //Number - amount extra to add to the radius to cater for hit detection outside the drawn point
        pointHitDetectionRadius : 200,
        //Boolean - Whether to show a stroke for datasets
        datasetStroke           : true,
        //Number - Pixel width of dataset stroke
        datasetStrokeWidth      : 2,
        //Boolean - Whether to fill the dataset with a color
        datasetFill             : false,
        //String - A legend template
        //legendTemplate          : '<ul class="<%=name.toLowerCase()%>-legend"><% for (var i=0; i<datasets.length; i++){%><li><span style="background-color:<%=datasets[i].lineColor%>"></span><%if(datasets[i].label){%><%=datasets[i].label%><%}%></li><%}%></ul>',
        //Boolean - whether to maintain the starting aspect ratio or not when responsive, if set to false, will take up entire container
        maintainAspectRatio     : true,
        //Boolean - whether to make the chart responsive to window resizing
        responsive              : true
    };
    return lineChartOptions;
}


drawChart('Chart_lad',coronaryData['lad_max'],coronaryData['lad_min'],'长轴','短轴');
drawChart('Chart_lad_ct',getJson(baseurl+'lad_ctmax.json'),getJson(baseurl+'lad_ctmin.json'),'Max','Min');
drawChart('Chart_lcx',coronaryData['lcx_max'],coronaryData['lcx_min'],'长轴','短轴');
drawChart('Chart_lcx_ct',getJson(baseurl+'lcx_ctmax.json'),getJson(baseurl+'lcx_ctmin.json'),'Max','Min');
drawChart('Chart_lm',coronaryData['lm_max'],coronaryData['lm_min'],'长轴','短轴');
drawChart('Chart_lm_ct',getJson(baseurl+'lm_ctmax.json'),getJson(baseurl+'lm_ctmin.json'),'Max','Min');
drawChart('Chart_rca',coronaryData['rca_max'],coronaryData['rca_min'],'长轴','短轴');
drawChart('Chart_rca_ct',getJson(baseurl+'rca_ctmax.json'),getJson(baseurl+'rca_ctmin.json'),'Max','Min');
