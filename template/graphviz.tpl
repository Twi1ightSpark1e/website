<!doctype html>
<html lang="en">
  <head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta http-equiv="refresh" content="30">

    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">

    <!-- Own styles -->
    <link rel="stylesheet" href="/static/common.css"></link>

    <style>
        img {
            width: 75%;
            content: url("data:image/svg+xml;base64,{{.Image}}");
        }
    </style>
  </head>
  <body class="bg-dark">
    <div class="container">
      <nav class="sticky-top" aria-label="breadcrumb">
        <ol class="breadcrumb">
          {{range .Breadcrumb}}
            <li class="breadcrumb-item"><a href="{{.Address}}">{{.Title}}</a></li>
          {{end}}
          <li class="breadcrumb-item active">{{.LastBreadcrumb}}</li>
        </ol>
      </nav>
      <div class="row justify-content-center">
        <h3 class="text-light text-center">Last update: {{.Timestamp}}</h3>
        <img></img>
      </div>
    </div>
  </body>
</html>
