<html>

<head>
    <title>{{.Title}}</title>
    <!-- Required meta tags -->
    <meta charset="utf-8"></meta>
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no"></meta>

    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous"></link>
    <!-- Own styles -->
    <link rel="stylesheet" href="/static/common.css"></link>
    <link rel="stylesheet" href="/static/autoindex.css"></link>
</head>

<body class="bg-dark">
    <div class="container">
        <!-- Breadcrumb -->
        <nav class="sticky-top" aria-label="breadcrumb">
            <ol class="breadcrumb">
                {{range .Breadcrumb}}
                <li class="breadcrumb-item"><a href="{{.Address}}">{{.Title}}</a></li>
                {{end}}
                <li class="breadcrumb-item active">{{.LastBreadcrumb}}</li>
            </ol>
        </nav>
        <!-- Files table -->
        {{if .Error}}
           <h4 class="text-light text-center">{{.Error}}</h4>
        {{else if .List}}
            <table class="table table-striped table-sm table-hover table-dark">
                <thead>
                    <tr>
                        <th scope="col">name</th>
                        <th scope="col">size</th>
                        <th scope="col">date</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .List}}
                    <tr class="filetable-entry">
                        <td>
                            <a href="{{.Name}}" class="text-light">{{if .IsDir}}<img class="mr-2 directory-icon"/>{{else}}<img class="mr-2 file-icon"/>{{end}}{{.Name}}</a>
                        </td>
                        <td>{{.Size}}</td>
                        <td>{{.Date}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        {{end}}
    </div>
</body>

</html>
