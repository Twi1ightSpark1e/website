{{if eq .MarkdownVisibility "pre"}}
  <div class="border rounded-top p-3">
    <h4 class="mb-0">{{.MarkdownTitle}}</h4>
  </div>
  <div class="border rounded-bottom p-3 mb-2">
    {{.MarkdownContent}}
  </div>
{{end}}
