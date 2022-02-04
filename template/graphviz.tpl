<html lang="en">

<head>
  {{template "head.tpl" .}}
  <style>
    img {
      width: 75% !important;
      content: url("data:image/svg+xml;base64,{{.Image}}");
    }
  </style>
</head>

{{template "body.tpl" .}}
  <div class="row justify-content-center">
    <h3 class="text-light text-center">Last update: {{.Timestamp}}</h3>
    <img></img>
  </div>
{{template "footer.tpl" .}}

</html>
