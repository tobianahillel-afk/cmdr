package main

import (
	"embed"
	"net/http"
)

//go:embed analyst.html analyst.css analyst.js
var analystAssets embed.FS

// The analyst page is a local, synthetic UI; the existing server-side Case,
// Timeline Entry and linked-source checks remain the sole data authority.
func installAnalystDemoRoutes(mux *http.ServeMux) {
	for _, asset := range []struct{ route, file, mime string }{
		{"/analyst", "analyst.html", "text/html; charset=utf-8"},
		{"/analyst.css", "analyst.css", "text/css; charset=utf-8"},
		{"/analyst.js", "analyst.js", "text/javascript; charset=utf-8"},
	} {
		data, err := analystAssets.ReadFile(asset.file)
		if err != nil {
			panic("missing embedded analyst demo asset: " + asset.file)
		}
		mux.HandleFunc(asset.route, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("Referrer-Policy", "no-referrer")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self'; connect-src 'self'; base-uri 'none'; form-action 'none'; frame-ancestors 'none'")
			if r.Method != http.MethodGet && r.Method != http.MethodHead {
				w.Header().Set("Allow", "GET, HEAD")
				http.Error(w, "method unavailable", http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", asset.mime)
			if r.Method == http.MethodGet {
				_, _ = w.Write(data)
			}
		})
	}
}
