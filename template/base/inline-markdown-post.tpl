{{if eq .MarkdownVisibility "post"}}
  <div class="border rounded-top p-3 overflow-auto">
    <h4 class="mb-0">{{.MarkdownTitle}}</h4>
  </div>
  <div class="border rounded-bottom p-3 mb-2 overflow-auto">
    {{.MarkdownContent}}
  </div>
{{end}}
