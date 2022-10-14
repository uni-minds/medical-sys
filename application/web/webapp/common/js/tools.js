function analysisURL(href) {
    let url ={};
    let data = href.split('?')[1].split('&');
    data.forEach(e=>{
        let p=e.split('=');
        url[p[0]]=p[1];
    });
    return url
}