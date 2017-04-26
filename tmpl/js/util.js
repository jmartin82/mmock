    function getColorByStatus(statusCode) {
        if (statusCode === 200 || statusCode === 201) {
            return "success";
        } else if (statusCode === 404) {
            return "danger";
        } else {
            return "warning";
        }
    }

    function getColorByMethod(method) {
        if (method === 'GET') {
            return "info";
        } else if (method === 'POST') {
            return "success";
        } else if (method === 'PUT') {
            return "success";
        } else if (method === 'DELETE') {
            return "danger";
        } else {
            return "warning";
        }
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
