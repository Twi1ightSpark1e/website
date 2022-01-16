<!doctype html>
<html lang="en">
  <head>
    <title>{{.LastBreadcrumb}}</title>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">

    <!-- Own styles -->
    <link rel="stylesheet" href="/static/common.css"></link>
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
        {{range .Cards}}
        <div class="card bg-light border-light m-3" style="width: 18rem;">
          <div class="card-body">
            <h5 class="card-title">{{.Title}}</h5>
            <p class="card-text">{{.Description}}</p>
          </div>
          <div class="card-footer">
            {{range .Links}}
            <a href="{{.Address}}" class="card-link">{{.Title}}</a>
            {{end}}
          </div>
        </div>
        {{else}}
        <h4 class="text-light text-center">No content for you, sorry \_(ツ)_/</h4>
        {{end}}
      </div>
    </div>
  </body>
</html>
