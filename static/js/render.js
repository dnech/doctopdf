$( document ).ready(function() {

    // GET DATA
    var data = {id: "", config: {}, data:{}};
    try {data = JSON.parse(window.TEMPLATE.data);} catch(e) {}
    console.log("DATA:", data);

    // RENDER TEMPLATE
    $("[type='template']").each(
        function (index, elem) {
            console.log("ELEMENT:", index, elem);

            var target = $($(elem).data("target"));
            var template = $.templates(elem);

            if(target.length) {
                target.html(template.render(data));
            } else {
                $(template.render(data)).insertAfter($(elem));
            }
        }
    );

    if (window.TEMPLATE.render) {
        window.TEMPLATE.render(data);
    }

});