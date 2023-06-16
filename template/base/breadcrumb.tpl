<nav class="sticky-top px-3 py-2 mb-3 bg-body-secondary rounded-bottom d-flex justify-content-between" aria-label="breadcrumb">
  <ol class="breadcrumb mb-0">
    {{range .Breadcrumb}}
      <li class="breadcrumb-item"><a href="{{.Address}}">{{.Title}}</a></li>
    {{end}}
    <li class="breadcrumb-item active">{{.LastBreadcrumb}}</li>
  </ol>
  <form method="post" action="">
    <button name="theme" value={{if eq .ThemeSwitch "light"}}"dark"{{else}}"light"{{end}} title="Switch theme" type="submit" class="btn px-2 py-0">
      <i class="bi bi-{{if eq .ThemeSwitch "light"}}moon-stars{{else}}sun{{end}}-fill"></i>
    </button>
  </form>
</nav>
