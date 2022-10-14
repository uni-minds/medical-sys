$(function(window, document) {
    let header=document.getElementById('header');
    let logout=document.createElement('a');
    logout.innerText='退出';
    logout.setAttribute('href','/logout');
    header.appendChild(logout)
}(this, document));