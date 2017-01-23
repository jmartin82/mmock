function RequestLogger() {
    //compile the template
    var source = $("#request-entry").html();
    var template = Handlebars.compile(source);

    this.num = 0;
    this.type = "macintosh";
    this.color = "red";

    function getColorByStatus(statusCode) {
        if (statusCode == 200 || statusCode == 201) {
            return "success";
        } else if (statusCode == 404) {
            return "danger";
        } else {
            return "warning";
        }
    }

    function updateTitle() {
        $(document).attr("title", "NEW REQUEST!!");
        setTimeout(function () { $(document).attr("title", "MMock Console"); }, 2000);
    }

    function syntaxHighlight(json) {
        json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
        return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function (match) {
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

    function updateLastRequestDate() {
        $('#last_updated').text(getCurrentTime());
    }

    function getContext(num, data) {
        var status = data.response.statusCode;
        var method = data.request.method;
        var path = data.request.path;
        var request = syntaxHighlight(JSON.stringify(data.request, undefined, 4));
        var response = syntaxHighlight(JSON.stringify(data.response, undefined, 4));
        var log = syntaxHighlight(JSON.stringify(data.result, undefined, 4));
        var color = getColorByStatus(status)

        return { request_num: num, request: request, response: response, rlog: log, request_date: getCurrentTime(), request_code: status, request_method: method, request_path: path, request_color: color };
    }

    this.logEntry = function (data) {
        this.num++;
        var context = getContext(this.num, data);
        var html = template(context);
        $("#request-table tbody").prepend(html);
        updateTitle();
        updateLastRequestDate();

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
