package workqueue

import (
	"bytes"
	"html/template"
)

// The browser is a second representation of the SAME server-authorized
// projection, not a second route, Saved View engine or trusted client.
var queuePage = template.Must(template.New("command-work-queue").Parse(`<!doctype html>
<html lang="fr">
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Command — Work Queue</title></head>
<body>
<header><h1>Work Queue</h1><p>Tenant : {{.Tenant}} — Environnement : {{.Environment}}</p></header>
<nav aria-label="Vues système">
{{range .Views}}<a href="/command/work-queue?view={{.Key}}"{{if .Selected}} aria-current="page"{{end}}>{{.Label}}</a> {{end}}
</nav>
<main>
<h2>Vue : {{.Current}}</h2>
<p role="status">État : {{.State}}</p>
{{if .Diagnostics}}<ul aria-label="Limites de la projection">{{range .Diagnostics}}<li>{{.}}</li>{{end}}</ul>{{end}}
{{if .Rows}}
<table><caption>Éléments autorisés — première page (50 maximum)</caption>
<thead><tr><th scope="col">Type</th><th scope="col">Identifiant</th><th scope="col">Responsable</th><th scope="col">État</th><th scope="col">Mis à jour (UTC)</th></tr></thead>
<tbody>{{range .Rows}}<tr><td>{{.Kind}}</td><th scope="row">{{.ID}}</th><td>{{.Owner}}</td><td>{{.Status}}</td><td>{{.UpdatedAt}}</td></tr>{{end}}</tbody>
</table>
{{else}}<p>Aucun élément autorisé dans cette vue.</p>{{end}}
</main></body></html>`))

type queuePageView struct {
	Key View
	Label string
	Selected bool
}

func renderQueuePage(result Result) ([]byte, error) {
	labels := [...]string{"Tout", "Incidents", "Tâches", "Non assignés", "Risque SLA", "Mon travail"}
	views := make([]queuePageView, 0, len(SystemViews))
	for i, key := range SystemViews {
		views = append(views, queuePageView{Key:key, Label:labels[i], Selected:key==result.View})
	}
	var out bytes.Buffer
	err := queuePage.Execute(&out, struct {
		Result
		Current View
		Views []queuePageView
	}{Result:result, Current:result.View, Views:views})
	if err != nil { return nil, err }
	return out.Bytes(), nil
}
