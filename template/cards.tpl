<!doctype html>
<html lang="en">

<head>
  {{template "base/head" .}}
</head>

{{template "base/body" .}}
  <div class="row gx-0 justify-content-center">
    {{range .Cards}}
    <div class="card bg-light border-light m-3" style="width: 18rem;">
      <div class="card-body text-dark">
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
    <h4 class="text-center">No content for you, sorry \_(ツ)_/</h4>
    {{end}}
  </div>
{{template "base/footer" .}}

</html>
