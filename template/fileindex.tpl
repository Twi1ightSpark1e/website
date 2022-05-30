<html>

<head>
  {{template "head.tpl" .}}
  <link rel="stylesheet" href="/static/autoindex.css"></link>
</head>

{{template "body.tpl" .}}
  {{if .List}}
  {{template "inline-markdown-pre.tpl" .}}
  <div class="d-flex flex-column flex-lg-row justify-content-between">
    <div class="mb-2 mb-lg-0 download-as">
      <span>Download as:</span>
      <a href="?type=tar"><button type="button" class="btn btn-light">.tar</button></a>
      <a href="?type=gz"><button type="button" class="btn btn-light">.tar.gz</button></a>
      <a href="?type=zst"><button type="button" class="btn btn-light">.tar.zst</button></a>
    </div>
    {{if .AllowUpload}}
    <form method="post" enctype="multipart/form-data">
      <div class="d-flex flex-row">
        <input class="form-control me-1" type="file" id="file" name="file">
        <input class="btn btn-light" type="submit" value="Upload">
      </div>
    </form>
    {{end}}
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
  {{template "inline-markdown-post.tpl" .}}
  {{end}}
{{template "footer.tpl" .}}

</html>
