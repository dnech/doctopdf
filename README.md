# doctopdf
Templates document to pdf converter service

```
// Install wkhtmltopdf in your system from https://wkhtmltopdf.org/downloads.html
// Create file "config.json"
//----------------------------------------------------------------------------------------------------------------------
{
  "host": ":8888",
  "static": "http://localhost:8888/static/",
  "public": "http://localhost:8888",
  "pdf": {
    "dpi": 300,
    "orientation": "Portrait",
    "grayscale": false,
    "pagesize": "A4",
    "margin_top": 5,
    "margin_bottom": 5,
    "margin_right": 5,
    "margin_left": 20
  }
}
//----------------------------------------------------------------------------------------------------------------------

// Example template in file "template/test.html"
// Example config for template in file "template/test.json"

// run server "doctopdf"
// See in browser html result "localhost:8888"
// See in browser html result "localhost:8888/doc/test"
// See in browser pdf result  "localhost:8888/pdf/test"

```

How it works:
- send in GET or POST method json param "data",
- this param integrate in templates as javascript params
- javascript parse json and replace template field