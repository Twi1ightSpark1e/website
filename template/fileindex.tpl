<html>

<head>
  {{template "head.tpl" .}}
  <link rel="stylesheet" href="/static/autoindex.css"></link>
</head>

{{template "body.tpl" .}}
  {{if .List}}
  <div class="download-as">
    <span class="text-light">Download as:</span>
    <a href="?type=tar"><button type="button" class="btn btn-light">.tar</button></a>
    <a href="?type=gz"><button type="button" class="btn btn-light">.tar.gz</button></a>
    <a href="?type=zst"><button type="button" class="btn btn-light">.tar.zst</button></a>
  </div>
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
{{template "footer.tpl" .}}

</html>
