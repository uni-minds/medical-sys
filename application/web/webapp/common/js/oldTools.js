function getJson(url) {
    var req = new XMLHttpRequest();
    if (req) {
        req.open('GET', url, false);
        req.send();
        //console.log(url, req.responseText);
        if (req.responseText === "") {
            return;
        } else {
            return JSON.parse(req.responseText);
        }
    }
}
function loadJson(url) {
    let result={};
    let xmlReqJson = new XMLHttpRequest();
    xmlReqJson.onreadystatechange = function () {
        if (this.readyState == 4 && this.status == 200) {
            result = JSON.parse(this.responseText);
        }
    };
    xmlReqJson.open("GET", url, false);
    xmlReqJson.send();
    return result;
}