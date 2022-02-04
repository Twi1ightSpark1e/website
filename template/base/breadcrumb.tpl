<nav class="sticky-top p-3 mb-3 bg-light rounded-bottom" aria-label="breadcrumb">
  <ol class="breadcrumb mb-0">
    {{range .Breadcrumb}}
      <li class="breadcrumb-item"><a href="{{.Address}}">{{.Title}}</a></li>
    {{end}}
    <li class="breadcrumb-item active">{{.LastBreadcrumb}}</li>
  </ol>
</nav>
