function sortTable(tableSelector, f, n){
	var rows = $(tableSelector + ' tbody  tr').get();

	rows.sort(function(a, b) {

		var A = getVal(a);
		var B = getVal(b);

		if(A < B) {
			return -1*f;
		}
		if(A > B) {
			return 1*f;
		}
		return 0;
	});

	function getVal(elm){
		var v = $(elm).children('td').eq(n).text().toUpperCase();
		if($.isNumeric(v)){
			v = parseInt(v,10);
		}
		return v;
	}

	$.each(rows, function(index, row) {
		$(tableSelector).children('tbody').append(row);
	});
}

var mappingElements = ['#uri', '#desc', '#method', '#path', '#result'];
var defaultSortingElement = '#uri';
var mappingTableSelector = '#mapping-table';
var selectorUp = '.glyphicon-chevron-up';
var selectorDown = '.glyphicon-chevron-down';
var selectedClassName = 'selected';

function attachListener(item, index, elements){
    $(item).find(selectorUp).on('click', function(){
        var n = $(item).prevAll().length;
        sortTable(mappingTableSelector, -1 ,n);
        markAsSelected($(this));
    });
    $(item).find(selectorDown).on('click', function(){
        var n = $(item).prevAll().length;
        sortTable(mappingTableSelector, 1 ,n);
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

function attachListeners() {
    mappingElements.forEach(attachListener);
}

function applyDefaultSorting() {
    $(defaultSortingElement).find(selectorDown).click();
}