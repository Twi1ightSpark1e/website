<html lang="en">

<head>
  {{template "head.tpl" .}}
</head>

{{template "body.tpl" .}}
  <div class="row justify-content-center">
    <h4 class="text-light text-center">{{.Error}}</h4>
  </div>
{{template "footer.tpl" .}}

</html>
