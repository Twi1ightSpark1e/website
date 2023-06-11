package util

import (
	"log"
	"net/http"
)

type Theme string
const (
	ThemeLight Theme = "light"
	ThemeDark        = "dark"
)

func HandleThemeToggle(w http.ResponseWriter, r *http.Request, logger *log.Logger) bool {
	r.ParseForm()

	if !r.Form.Has("theme") {
		return false
	}

	cookie := http.Cookie{
		Name: "theme",
		Value: r.FormValue("theme"),
		Path: "/",
		MaxAge: 2147483647,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)

	return true
}

func GetTheme(r *http.Request) Theme {
	var theme Theme = ThemeDark

	cookie, err := r.Cookie("theme")
	if err == nil {
		theme = Theme(cookie.Value)
	}

	return theme
}
