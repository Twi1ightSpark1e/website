<!doctype html>
<html>

<head>
  {{template "head.tpl" .}}
  <link rel="stylesheet" type="text/css" href="/static/autoindex.css"></link>
</head>

{{template "body.tpl" .}}
  {{if .List}}
  {{template "inline-markdown-pre.tpl" .}}
  <div class="toolbar mb-1 container row row-cols-1 row-cols-xl-3 mx-0 px-0 justify-content-between">
    <div class="col">
        {{if .AllowDownload}}
        <div class="input-group w-auto">
            <span class="input-group-text">Download as:</span>
            <a href="?type=tar" type="button" class="btn btn-light" role="button">.tar</a>
            <a href="?type=gz" type="button" class="btn btn-light" role="button">.tar.gz</a>
            <a href="?type=zst" type="button" class="btn btn-light" role="button">.tar.zst</a>
        </div>
        {{end}}
    </div>
    <form method="get" class="col">
      <div class="d-flex flex-row">
        <div class="input-group" role="group">
          <input class="form-control" type="text" placeholder="Find file" aria-label="Find file" name="query" value="{{.FindQuery}}">
          <input type="checkbox" class="btn-check" id="btn-matchcase" autocomplete="off" name="matchcase" {{if .FindMatchCase}}checked{{end}}>
          <label class="btn btn-outline-light" for="btn-matchcase" title="Match case">Aa</label>
          <input type="checkbox" class="btn-check" id="btn-regex" autocomplete="off" name="regex" {{if .FindRegex}}checked{{end}}>
          <label class="btn btn-outline-light" for="btn-regex" title="Use regular expression">.*</label>
          <button type="submit" class="btn btn-light">Find</button>
          {{if .FindQuery}}
          <a href="{{.URL}}" type="button" class="btn btn-light" role="button" title="Clear">X</a>
          {{end}}
        </div>
      </div>
    </form>
    <form method="post" enctype="multipart/form-data" class="col">
      {{if .AllowUpload}}
      <div class="d-flex flex-row">
        <input class="form-control me-1" type="file" id="file" name="file">
        <input class="btn btn-light" type="submit" value="Upload">
      </div>
      {{end}}
    </form>
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
            <a href="{{.Name}}" class="text-light">
              {{if .IsDir}}<img class="me-2 directory-icon"/>{{else}}<img class="me-2 file-icon"/>{{end}}
              {{.Name}}
            </a>
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
