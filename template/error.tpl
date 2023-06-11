{{template "base/html" .}}

<head>
  {{template "base/head" .}}
</head>

{{template "base/body" .}}
  <div class="justify-content-center">
    <h4 class="text-center">{{.Error}}</h4>
  </div>
{{template "base/footer" .}}

</html>
