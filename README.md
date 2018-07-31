# doctopdf
Templates document to pdf converter service

```
// Install wkhtmltopdf in your system from https://wkhtmltopdf.org/downloads.html
// Create file "config.json"
//----------------------------------------------------------------------------------------------------------------------
{
  "host": ":5555",
  "static": "http://localhost:5555/static/",
  "pdf": {
    "dpi": 300,
    "orientation": "Portrait",
    "grayscale": false,
    "pagesize": "A4"
  }
}
//----------------------------------------------------------------------------------------------------------------------

// Create file "template/test.html"
//----------------------------------------------------------------------------------------------------------------------
<html>
<head>
    <meta charset="utf-8">
    <title>${template.id}</title>
    <link rel="stylesheet" href="${template.static}/css/bootstrap.min.css">
    <script src="${template.static}/js/jquery-3.3.1.min.js"></script>
    <script src="${template.static}/js/jsrender.min.js"></script>
    <script src="${template.static}/js/render.js"></script>
</head>
<body>
    <script>
        window.TEMPLATE = {
            id: '${template.id}',
            static: '${template.static}',
            public: '${template.public}',
            data: '${template.data}',
            render: function(data) {
                console.log("After render", data);
            }
        };
    </script>

    <div class="alert alert-primary" role="alert">Use <a href="https://www.jsviews.com/#jsrapi">JsRender Api</a></div>

    <div id="render" class="alert alert-info" role="alert"></div>

    <script type="template" data-target="#render">
    <div>
        <b><em>Id:</em></b> "{{:id}}"<br/>
        <b><em>Config:</em></b><br/>
        <ul>
            {{props config}}
                <li>"{{>key}}": "{{>prop}}"</li>
            {{/props}}
        </ul>
        <b><em>Data:</em></b><br/>
        <ul>
        {{props data}}
            <li>"{{>key}}": "{{>prop}}"</li>
        {{/props}}
        </ul>
    </div>

    <div>
        <b><em>Name:</em></b> "{{:data.name}}"<br/>
        <b><em>Email:</em></b> "{{:data.email}}"<br/>
    </div>

    </script>


</body>
</html>
//----------------------------------------------------------------------------------------------------------------------

// Create file "template/test.json"
//----------------------------------------------------------------------------------------------------------------------
{
    "dpi": 300,
    "orientation": "Portrait",
    "grayscale": true,
    "pagesize": "A4"
}
//----------------------------------------------------------------------------------------------------------------------

// run server "doctopdf"
// See in browser html result "localhost:5555"
// See in browser html result "localhost:5555/doc/test"
// See in browser pdf result  "localhost:5555/pdf/test"

```

How it works:
- send in GET or POST method json param "data",
- this param integrate in templates as javascript params
- javascript parse json and replace template field