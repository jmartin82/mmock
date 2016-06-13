var count = 0;
var requests = [];

$(document).ready(function() {
    $("#btnClearLog").click(function() {
        $("#groupConsole").empty();
        $("#tirecap").hide();
        $("#hdrequest").html("");
        $("#hdresponse").html("");
        $("#hdlog").html("");
    });
});

function incrementCount() {
    return count++;
}

function showDetails(id) {
    logDetails(requests[id]);
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
    var request = JSON.stringify(json.request, undefined, 4);
    var response = JSON.stringify(json.response, undefined, 4);
    var log = JSON.stringify(json.result, undefined, 4);
    var status = json.response.statusCode;
    $("#tirecap").show();
    $("#tirecap").attr('class', 'alert alert-' + getColorByStatus(status));
    $("#tistatus").html(status);
    $("#tirequest").html(json.request.method + " " + json.request.path);
    $("#hdrequest").html(syntaxHighlight(request));
    $("#hdresponse").html(syntaxHighlight(response));
    $("#hdlog").html(syntaxHighlight(log));
}


function logRequest(json) {
    var status = json.response.statusCode;
    var id = incrementCount();
    var datetime = getCurrentTime();
    var fullLog = datetime + " <- " + json.request.method + " " + json.request.path;
    requests[id] = json;
    $("#groupConsole").append('<li class="list-group-item list-group-item-' + getColorByStatus(status) + '" onclick="showDetails(' + id + ');return false">' + fullLog + '</li>');
    logDetails(json)
    clearOldLogs();
    scrollLogsDown();
}


function clearOldLogs() {
    var logItemSize = $("#groupConsole li").size();
    if (logItemSize > 50) {
        $('#groupConsole li:first').remove();
    }
}

function scrollLogsDown() {
    if ($("#chkAutoScroll").is(':checked')) {
        $("#groupConsole").scrollTop($("#groupConsole").get(0).scrollHeight);
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