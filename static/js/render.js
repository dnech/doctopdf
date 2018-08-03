$( document ).ready(function() {

    // GET DATA
    var data = {id: "", config: {}, data:{}};
    try {
        data = JSON.parse(window.TEMPLATE.data);
        if (data['data']) {
            for (var key in data['data']) {
                try {
                    data['data'][key] = JSON.parse(data['data'][key]);
                } catch (e) {
                    //console.error("DATA '" + key + "':", data['data'][key]);
                    //console.error("DATA '" + key + "':", e);
                }
            }
        }
    } catch(e) {
        console.error("DATA:", window.TEMPLATE.data);
        console.error("DATA:", e);
    }
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