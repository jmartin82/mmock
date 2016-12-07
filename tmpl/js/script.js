var count = 0;
var requests = [];

$(document).ready(function() {
    $("#btnClearConsole").click(function() {
        $("#groupConsole").empty();
        $("#tirecap").hide();
        $("#hdrequest").html("");
        $("#hdresponse").html("");
        $("#hdpersist").html("");
        $("#hdlog").html("");
    });

    $("#btnClearLog").click(function() {
        $("#groupLog").empty();
    });
});

function incrementCount() {
    return count++;
}

function showDetails(id) {
    setRowSelected(id);
    logDetails(requests[id]);
}

function setRowSelected(id) {
   $(".list-group-item").removeClass('selected_row');
   $("#row-request-"+id).addClass('selected_row');
}

function getColorByStatus(statusCode) {

    if (statusCode == 200 || statusCode == 201) {
        return "success";
    } else if (statusCode == 404) {
        return "danger";
    } else {
        return "warning";
    }
}

function logDetails(json) {
    $("#tirecap").fadeOut(100);

    if(json.response.persisted && json.response.persisted.name == ""){
        delete json.response["persisted"];
    }

    if(json.persist && json.persist.name == ""){
        delete json.persist["name"];
    }
    if(json.persist.delete == false){
        delete json.persist["delete"];
    }
    if(json.persist.amqp && json.persist.amqp.url == ""){
        delete json.persist["amqp"];
    }

    var request = JSON.stringify(json.request, undefined, 4);
    var response = JSON.stringify(json.response, undefined, 4);
    var persist = JSON.stringify(json.persist, undefined, 4);

    var log = JSON.stringify(json.result, undefined, 4);
    var status = json.response.statusCode;
    $("#tirecap").attr('class', 'alert alert-' + getColorByStatus(status));
    $("#tirecap").fadeIn(100);
    $("#tistatus").html(status);
    $("#tirequest").html(json.request.method + " " + json.request.path);
    $("#hdrequest").html(syntaxHighlight(request));
    $("#hdresponse").html(syntaxHighlight(response));
    $("#hdpersist").html(syntaxHighlight(persist));
    $("#hdlog").html(syntaxHighlight(log));
}


function logRequest(json) {
    var status = json.response.statusCode;
    var id = incrementCount();
    var datetime = getCurrentTime();
    var fullLog = datetime + " <- " + json.request.method + " " + json.request.path;
    requests[id] = json;
    $("#groupConsole").append('<li id="row-request-' + id + '" class="list-group-item list-group-item-' + getColorByStatus(status) + '" onclick="showDetails(' + id + ');return false">' + fullLog + '</li>');
    showDetails(id)
    clearOldLogs();
    scrollConsoleDown();
}

function writeLog(text) {
    var datetime = getCurrentTime();
    var fullLog = datetime + " <- " + text;
    $("#groupLog").append('<li class="list-group-item">' + fullLog + '</li>');
    clearOldLogs();
    scrollLogsDown();
}


function clearOldLogs() {
    var logItemSize = $("#groupConsole li").size();
    if (logItemSize > 50) {
        $('#groupConsole li:first').remove();
    }
}

function scrollConsoleDown() {
    if ($("#chkAutoScroll").is(':checked')) {
        $("#divGroupConsole").scrollTop($("#divGroupConsole").get(0).scrollHeight);
    }
}

function scrollLogsDown() {
    if ($("#chkLogAutoScroll").is(':checked')) {
        $("#divGroupLog").scrollTop($("#divGroupLog").get(0).scrollHeight);
    }
}

function getCurrentTime() {
    var currentdate = new Date();
    var datetime = currentdate.getDate() + "/" +
        (currentdate.getMonth() + 1) + "/" +
        currentdate.getFullYear() + " @ " +
        currentdate.getHours() + ":" +
        currentdate.getMinutes() + ":" +
        currentdate.getSeconds();
    return datetime;
}

function syntaxHighlight(json) {
    json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
    return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function(match) {
        var cls = 'number';
        if (/^"/.test(match)) {
            if (/:$/.test(match)) {
                cls = 'key';
            } else {
                cls = 'string';
            }
        } else if (/true|false/.test(match)) {
            cls = 'boolean';
        } else if (/null/.test(match)) {
            cls = 'null';
        }
        return '<span class="' + cls + '">' + match + '</span>';
    });
}