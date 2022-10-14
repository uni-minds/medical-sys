function unGzip(b64Data) {
    let binData = atob(b64Data);
    // Convert binary string to character-number array
    let charData = binData.split('').map((x) => {
        return x.charCodeAt(0);
    });
    let unpackData = pako.inflate(new Uint8Array(charData));
    // Convert gunzipped byteArray back to ascii string:
    let strData = ""
    let step = 50000
    for (let c = 0; c < unpackData.length; c += step) {
        let str = String.fromCharCode.apply(null, new Uint16Array(unpackData.slice(c, c + step)))
        // console.log("count",c)
        strData = strData.concat(str)
    }
    // strData = String.fromCharCode.apply(null, new Uint16Array(unpackData));
    return strData;
}

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
    console.log("LD:",url)
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