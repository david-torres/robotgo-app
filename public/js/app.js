var ws = new WebSocket("ws://" + window.location.host + "/ws");

// init websocket
ws.onmessage = function (e) {
    var message = JSON.parse(e.data);
    console.log("Websocket Message", message.msg);
    var p = document.createElement('p');
    p.innerHTML = message.msg;
    document.getElementById('app').appendChild(p);
    scrollDiv();
};

function scrollDiv() {
    var appDiv = document.getElementById("app");
    appDiv.scrollTop = appDiv.scrollHeight - appDiv.clientHeight;
}