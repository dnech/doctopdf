$( document ).ready(function() {

    function appendFormData(form){
        var unindexed_array = form.serializeArray();
        var indexed_array = {};

        $.map(unindexed_array, function(n, i){
            indexed_array[n['name']] = n['value'];
        });

        var input = $("<input>").attr("type", "hidden").attr("name", "data").val(JSON.stringify(indexed_array));
        form.append($(input));
    }

    $("form[data-type='pdf']").each(
        function (index, elem) {
            console.log("PDF FORM:", index, elem);
            $(elem).submit(function( event ) {
                appendFormData($(this));
                //event.preventDefault();
            });
        }
    );

});