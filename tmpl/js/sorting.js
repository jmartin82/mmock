function Sorting(tableSelector, tableElements) {

    let selectorUp = '.glyphicon-chevron-up';
    let selectorDown = '.glyphicon-chevron-down';
    let selectedClassName = 'selected';

    function sortTable(order, itemsCount){
        let rows = $(tableSelector + ' tbody  tr').get();

        rows.sort(function(a, b) {

            let A = getVal(a);
            let B = getVal(b);

            if(A < B) {
                return -1*order;
            }
            if(A > B) {
                return 1*order;
            }
            return 0;
        });

        function getVal(elm){
            let v = $(elm).children('td').eq(itemsCount).text().toUpperCase();
            if($.isNumeric(v)){
                v = parseInt(v,10);
            }
            return v;
        }

        $.each(rows, function(index, row) {
            $(tableSelector).children('tbody').append(row);
        });
    }

    function attachListener(item, index, elements){
        $(item).find(selectorUp).on('click', function(){
            let itemsCount = $(item).prevAll().length;
            sortTable(-1 ,itemsCount);
            markAsSelected($(this));
        });
        $(item).find(selectorDown).on('click', function(){
            let itemsCount = $(item).prevAll().length;
            sortTable(1 ,itemsCount);
            markAsSelected($(this));
        });
    }

    function markAsSelected(object) {
        clearSelectedClass();
        object.addClass(selectedClassName);
    }

    function clearSelectedClass() {
        $(selectorUp).removeClass(selectedClassName);
        $(selectorDown).removeClass(selectedClassName);
    }

    function init() {
        tableElements.forEach(attachListener);
    }

    return {
        init : init
    }
}