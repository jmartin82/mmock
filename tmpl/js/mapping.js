function Mapping(domain) {

    this.domain = domain
    //compile the template
    var mapping_source = $("#mapping-entry").html();
    var mapping_template = Handlebars.compile(mapping_source);

    var fillList = function() {
        $.getJSON("http://" + domain + "/api/mapping", function(data) {
            var items = [];
            $("#mapping-table tbody").empty();
            $.each(data, function(key, mapping) {
                var html = mapping_template(mapping);
                $("#mapping-table tbody").append(html);
            });
        });

    };

    this.fillList = fillList;

    $('#mapping-data').on('click', '.btn-view-mapping', function() {
        var uri = $(this).data("uri");
        $.getJSON("http://" + domain + "/api/mapping/" + uri, function(data) {
            content = JSON.stringify(data, null, 4)
            BootstrapDialog.show({
                title: 'Mapping definition',
                message: "<div class=\"pre\">" + syntaxHighlight(content) + "</div>"
            });
        });
    });


    $('#mapping-data').on('click', '.btn-edit-mapping', function() {
        var uri = $(this).data("uri");
        var endpoint = "http://" + domain + "/api/mapping/" + uri;
        $.getJSON(endpoint, function(data) {
            content = JSON.stringify(data)
            BootstrapDialog.show({
                title: 'Mapping edit',
                message: "<textarea id='text-update-mapping' style='min-width:100%;height:300px'>" + content + "</textarea>",
                buttons: [{
                    label: 'Cancel',
                    action: function(dialog) {
                        dialog.close();
                    }
                }, {
                    label: 'Save',
                    action: function(dialog) {
                       var content = $('#text-update-mapping').val();
                       $.ajax({
                            type: 'PUT',
                            url: endpoint,
                            data: content,
                            success: function(data) {  dialog.close(); fillList() },
                            error: function(data) { alert("Error: "+JSON.stringify(data)); dialog.close(); },
                            contentType: "application/json",
                            dataType: 'json'
                        });
                    }
                }]
            });
        });
    });





}