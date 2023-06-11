<!doctype html>
<html lang="en">

<head>
  {{template "base/head" .}}
  <style>
    img {
      width: 75% !important;
      content: url("data:image/svg+xml;base64,{{.Image}}");
    }
  </style>
</head>

{{template "base/body" .}}
  <div class="row justify-content-center">
    <h3 class="text-center">Last update: {{.Timestamp}}</h3>
    <img></img>
  </div>
{{template "base/footer" .}}

</html>
