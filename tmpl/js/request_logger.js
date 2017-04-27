function RequestLogger() {
    //compile the template
    var source = $("#request-entry").html();
    var template = Handlebars.compile(source);

    this.num = 0

    function updateTitle() {
        $(document).attr("title", "NEW REQUEST!!");
        setTimeout(function () { $(document).attr("title", "MMock Console"); }, 2000);
    }

    function getRequestTime(timestamp) {
        var requestTime = new Date(timestamp*1000);
        var datetime = requestTime.getDate() + "/" +
            (requestTime.getMonth() + 1) + "/" +
            requestTime.getFullYear() + " @ " +
            requestTime.getHours() + ":" +
            requestTime.getMinutes() + ":" +
            requestTime.getSeconds();
        return datetime;
    }

    function updateLastRequestDate(timestamp) {
        $('#last_updated').text(getRequestTime(timestamp));
    }

    function getContext(num, data) {
        var status = data.response.statusCode;
        var method = data.request.method;
        var path = data.request.path;
        var request = syntaxHighlight(JSON.stringify(data.request, undefined, 4));
        var response = syntaxHighlight(JSON.stringify(data.response, undefined, 4));
        var log = syntaxHighlight(JSON.stringify(data.result, undefined, 4));
        var color = getColorByStatus(status)

        return { request_num: num, request: request, response: response, rlog: log, request_date: getRequestTime(data.time), request_code: status, request_method: method, request_path: path, request_color: color };
    }

    this.logEntry = function (data) {
        this.num++;
        var context = getContext(this.num, data);
        var html = template(context);
        $("#request-table tbody").prepend(html);
        updateTitle();
        updateLastRequestDate(data.time);

    };
}

$(document).ready(function () {
    $("#btnClearLog").click(function () {
        $("#request-table tbody").empty();
    });

    $('#request-data').on('click', 'table tr', function () {
        var id = $(this).data("target");
        $("#" + id).toggle();
    });


});
