# doctopdf
Templates document to pdf converter service

// Install wkhtmltopdf in your system from https://wkhtmltopdf.org/downloads.html
// Create folder "static"
// Create folder "template"
// Create folder "temp"
// Create file "config.json"
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

// Create file "template/test.html"
<html>
<head>
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.2/css/bootstrap.min.css">
</head>
<body>
    <div class="alert alert-primary" role="alert">This is a primary alert—check it out!</div>
    <div class="alert alert-secondary" role="alert"> This is a secondary alert—check it out! </div>
    <div class="alert alert-success" role="alert"> This is a success alert—check it out! </div>
    <div class="alert alert-danger" role="alert"> This is a danger alert—check it out! </div>
    <div class="alert alert-warning" role="alert"> This is a warning alert—check it out! </div>
    <div class="alert alert-info" role="alert"> This is a info alert—check it out! </div>
    <div class="alert alert-light" role="alert"> This is a light alert—check it out! </div>
    <div class="alert alert-dark" role="alert"> This is a dark alert—check it out! </div>
</body>
</html>

// run server "doctopdf"
// See in browser html result "localhost:5555/doc/test"
// See in browser pdf result  "localhost:5555/pdf/test"
