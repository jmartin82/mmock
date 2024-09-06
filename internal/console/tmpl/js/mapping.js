function Mapping(domain) {

    this.domain = domain
    //compile the template
    var mapping_source = $("#mapping-entry").html();
    var mapping_template = Handlebars.compile(mapping_source);

    var fillList = function() {
        $.getJSON(window.location.protocol + "//" + domain + "/api/mapping", function(data) {
            var items = [];
            $("#mapping-table tbody").empty();
            $.each(data, function(key, mapping) {
                mapping['status_color'] =  getColorByStatus(mapping.response.statusCode);
                mapping['method_color'] =  getColorByMethod(mapping.request.method);
                var html = mapping_template(mapping);
                $("#mapping-table tbody").append(html);
            });
        });

    };

    this.fillList = fillList;

    $('#mapping').on('click', '.btn-create-mapping', function() {
        BootstrapDialog.show({
            title: 'Mapping create',
            message: $('<textarea id="text-create-mapping" style="min-width:570px;min-height:300px"></textarea>'),
            buttons: [{
                label: 'Cancel',
                action: function(dialog) {
                    dialog.close();
                }
            }, {
                label: 'Save',
                cssClass: 'btn-primary',
                action: function(dialog) {
                   var content = $('#text-create-mapping').val();
                   var uri = JSON.parse(content).URI;
                   var endpoint = window.location.protocol + "//" + domain + "/api/mapping/" + uri;
                   $.ajax({
                        type: 'POST',
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

    $('#mapping-data').on('click', '.btn-view-mapping', function() {
        var uri = $(this).data("uri");
        $.getJSON(window.location.protocol + "//" + domain + "/api/mapping/" + uri, function(data) {
            var content = JSON.stringify(data, null, "\t")
            BootstrapDialog.show({
                title: 'Mapping definition',
                message: '<pre>' + syntaxHighlight(content) + '</pre>'
            });
        });
    });


    $('#mapping-data').on('click', '.btn-edit-mapping', function() {
        var uri = $(this).data("uri");
        var endpoint = window.location.protocol + "//" + domain + "/api/mapping/" + uri;
        $.getJSON(endpoint, function(data) {
            var content = JSON.stringify(data, null,"\t")
            var $text = $('<textarea id="text-update-mapping" style="min-width:570px;min-height:300px"></textarea>');
            $text.val(content)
            
            BootstrapDialog.show({
                title: 'Mapping edit',
                message: $text,
                buttons: [{
                    label: 'Cancel',
                    action: function(dialog) {
                        dialog.close();
                    }
                }, {
                    label: 'Save',
                    cssClass: 'btn-primary',
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
