<html lang="en">
  <head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta http-equiv="refresh" content="30">

    <!-- Bootstrap CSS -->
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3" crossorigin="anonymous">

    <!-- Own styles -->
    <link rel="stylesheet" href="/static/common.css"></link>

    <style>
        img {
            width: 75% !important;
            content: url("data:image/svg+xml;base64,{{.Image}}");
        }
    </style>
  </head>
  <body class="bg-dark">
    <div class="container">
      <nav class="sticky-top p-3 mb-3 bg-light rounded-bottom" aria-label="breadcrumb">
        <ol class="breadcrumb mb-0">
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
