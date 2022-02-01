<html>

<head>
    <!-- Required meta tags -->
    <meta charset="utf-8"></meta>
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no"></meta>

    <!-- Bootstrap CSS -->
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3" crossorigin="anonymous">
    <!-- Own styles -->
    <link rel="stylesheet" href="/static/common.css"></link>
    <link rel="stylesheet" href="/static/autoindex.css"></link>
</head>

<body class="bg-dark">
    <div class="container">
        <!-- Breadcrumb -->
        <nav class="sticky-top p-3 mb-3 bg-light rounded-bottom" aria-label="breadcrumb">
            <ol class="breadcrumb mb-0">
                {{range .Breadcrumb}}
                <li class="breadcrumb-item"><a href="{{.Address}}">{{.Title}}</a></li>
                {{end}}
                <li class="breadcrumb-item active">{{.LastBreadcrumb}}</li>
            </ol>
        </nav>
        <!-- Files table -->
        {{if .List}}
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
                            <a href="{{.Name}}" class="text-light">{{if .IsDir}}<img class="me-2 directory-icon"/>{{else}}<img class="me-2 file-icon"/>{{end}}{{.Name}}</a>
                        </td>
                        {{if .IsDir}}<td></td>{{else}}<td>{{.Size}}</td>{{end}}
                        <td>{{.Date}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        {{end}}
    </div>
</body>

</html>
