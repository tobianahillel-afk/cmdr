package workqueue

import (
	"bytes"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
)

// This bounded Incident Detail Class-0 slice is a distinct read-only
// presentation of the existing Command owner and selected-source projections.
// It does NOT implement Incident transitions, editing, action requests,
// durable audit, or the canonical production Incident Detail screen.
var incidentDetailPage = template.Must(template.New("incident-detail").Parse(`<!doctype html>
<html lang="fr"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Command — Incident Detail</title></head>
<body><header><h1>Incident Detail — aperçu autorisé en lecture seule</h1><p>Tenant : {{.Tenant}} · Environnement : {{.Environment}}</p></header>
<main aria-label="Dossier Incident"><p><strong>Incident {{.Selected.ID}}</strong></p>
<p>État : {{if .Selected.Status}}{{.Selected.Status}}{{else}}inconnu{{end}} · Responsable : {{if .Selected.Owner}}{{.Selected.Owner}}{{else}}inconnu{{end}} · Dernière mise à jour source : {{.Selected.UpdatedAt}}</p>
<p>Fraîcheur : {{if .Selected.Freshness}}{{.Selected.Freshness.State}}{{if .Selected.Freshness.Source}} · {{.Selected.Freshness.Source}}{{end}}{{else}}inconnue — aucune politique source{{end}}</p>
<h2>Informations métier de l'Incident (propriétaire Command)</h2>
{{if .Overview}}
<p>Titre : {{if .Overview.Title}}{{.Overview.Title}}{{else}}inconnu — non fourni{{end}}</p>
<p>Priorité : {{if .Overview.Priority}}{{.Overview.Priority}}{{else}}inconnue — non fournie{{end}} · Sévérité : {{if .Overview.Severity}}{{.Overview.Severity}}{{else}}inconnue — non fournie{{end}}</p>
<p>Impact métier : {{if .Overview.BusinessImpact}}{{.Overview.BusinessImpact}}{{else}}inconnu — non fourni{{end}}</p>
<p>Prochaine action Command : {{if .Overview.NextAction}}{{.Overview.NextAction}}{{else}}inconnue — non fournie{{end}}{{if .Overview.NextAction}} · échéance : {{if .Overview.NextActionDueAt}}{{.Overview.NextActionDueAt}}{{else}}non précisée{{end}}{{end}}</p>
<p>Provenance Command : {{.Overview.SourceRef}} · version {{.Overview.SourceVersion}} · observé {{.Overview.ObservedAt}}</p>
{{else}}<p>Informations métier indisponibles : aucune projection Incident autorisée fournie. Ne pas déduire titre, priorité, sévérité, impact ni prochaine action.</p>{{end}}
<h2>Coordination et informations indépendamment autorisées</h2>
<dl>{{range .Factors}}<dt>{{.Name}}</dt><dd>{{.Value}}{{if .Source}} · {{.Source}}{{end}}{{if .ObservedAt}} · {{.ObservedAt}}{{end}}</dd>{{end}}</dl>
<p>SLA : {{if .SelectedSLA}}{{.SelectedSLA.State}}{{if .SelectedSLA.PolicySource}} · source {{.SelectedSLA.PolicySource}}{{end}}{{else}}inconnu — source SLA indisponible ou non autorisée{{end}}</p>
<p>Propriétaire canonique : {{if .SelectedCanonicalOwnership}}{{if .SelectedCanonicalOwnership.ProductOwner}}{{.SelectedCanonicalOwnership.ProductOwner}}{{else}}inconnu{{end}}{{else}}non fourni{{end}}</p>
<h2>Alertes d'origine — sources Command autorisées (lecture seule)</h2>
{{if .AlertOrigins}}
<ul>{{range .AlertOrigins}}<li>Alert {{.ID}} · sévérité {{if .Severity}}{{.Severity}}{{else}}inconnue{{end}} · état {{if .Status}}{{.Status}}{{else}}inconnu{{end}} · disposition {{if .Disposition}}{{.Disposition}}{{else}}inconnue{{end}} · source {{.SourceRef}} · version {{.SourceVersion}} · observée {{.ObservedAt}}</li>{{end}}</ul>
{{else}}<p>Origine Alert non disponible ou non autorisée. Aucune association n’est déduite de l’Incident.</p>{{end}}
<h2>Résultats Govern liés — lecture indépendante (aucune action)</h2>
{{if .GovernResults}}
<ul>{{range .GovernResults}}<li>Result {{.ID}} · Response Run {{.RunID}} · état Govern {{if .Lifecycle}}{{.Lifecycle}}{{else}}inconnu{{end}} · outcome {{if .Outcome}}{{.Outcome}}{{else}}inconnu{{end}} · vérification {{if .Verification}}{{.Verification}}{{else}}non renseignée{{end}} · risque résiduel {{if .ResidualRisk}}{{.ResidualRisk}}{{else}}non renseigné{{end}} · provenance {{.SourceRef}} version {{.SourceVersion}} · créé {{.CreatedAt}}</li>{{end}}</ul>
{{else}}<p>Result Govern non disponible ou non autorisé. Aucun succès ni risque résiduel n’est déduit.</p>{{end}}
<h2>Investigations liées (lecture seule)</h2>
{{if .SelectedLinkedCases}}{{if .SelectedLinkedCases.Cases}}
<ul>{{range .SelectedLinkedCases.Cases}}<li>Case {{.CaseID}} — {{.Status}} · <a href="/case-timeline?case_id={{urlquery .CaseID}}&amp;origin_incident={{urlquery $.Selected.ID}}&amp;origin_view=detail" rel="noreferrer">Ouvrir le Case (autorisation distincte)</a></li>{{end}}</ul>
{{else}}<p>Aucun Case autorisé visible.</p>{{end}}{{else}}<p>Relations inconnues — source de liaison non fournie.</p>{{end}}
<p>Statut du dossier : partiel. Décisions et actions de modification non disponibles, autres objets Govern non fournis par cet adaptateur ; aucune valeur n’est inventée.</p>
<nav aria-label="Retour à la Work Queue"><a href="{{.ReturnHref}}">Retour à l’Incident dans la Work Queue (nouvelle autorisation)</a></nav>
<p>Aucune action de modification, d’export ou de réponse n’est accessible depuis cette vue.</p>
</main></body></html>`))

type incidentDetailView struct {
	Result
	Factors []factorLine
	Overview *AuthorizedIncidentOverview
	AlertOrigins []AuthorizedIncidentAlert
	GovernResults []AuthorizedGovernResult
	ReturnHref string
}

// NewIncidentDetailHTTPHandler is a separate read-only canonical Command
// Incident route. It delegates all owner decisions and optional data-source
// permissions to the SAME Project and enrichSelectedResult as Work Queue.
func NewIncidentDetailHTTPHandler(sessions SessionResolver, owner Owner) http.Handler {
	return NewIncidentDetailHTTPHandlerWithResults(sessions, owner, nil, nil)
}

// NewIncidentDetailHTTPHandlerWithResults injects the Shared link reader and
// Govern Result.read owner separately. A Command Incident grant is never
// used as permission to inspect Govern's source object.
func NewIncidentDetailHTTPHandlerWithResults(sessions SessionResolver, owner Owner,
	links IncidentResultReferencesOwner, govern GovernResultOwner) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; base-uri 'none'; form-action 'none'; frame-ancestors 'none'")
		w.Header().Set("X-Frame-Options", "DENY")
		if r == nil || r.URL == nil || r.URL.Path != "/command/incident-detail" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if sessions == nil || owner == nil {
			http.Error(w, "incident detail unavailable", http.StatusServiceUnavailable)
			return
		}
		if !strings.Contains(r.Header.Get("Accept"), "text/html") {
			http.Error(w, "incident detail requires HTML", http.StatusNotAcceptable)
			return
		}
		query, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil || len(query) != 1 || len(query["incident"]) != 1 {
			http.Error(w, "invalid incident detail selection", http.StatusBadRequest)
			return
		}
		id := query.Get("incident")
		if len(id) == 0 || len(id) > 128 || !utf8.ValidString(id) ||
			invalidScope(id) || strings.HasPrefix(id, "Task:") {
			http.Error(w, "invalid incident detail selection", http.StatusBadRequest)
			return
		}
		if r.Context().Err() != nil {
			http.Error(w, "incident detail unavailable", http.StatusServiceUnavailable)
			return
		}
		scope, err := sessions.ResolveSession(r.Context(), r)
		if err != nil || invalidScope(scope.Principal) || invalidScope(scope.Tenant) ||
			invalidScope(scope.Environment) {
			http.Error(w, "incident detail session unavailable", http.StatusUnauthorized)
			return
		}
		// FocusSelection locates the target exclusively in the entire fresh
		// tenant/env/principal CanRead-authorized bounded Incident cohort.
		// A client-supplied ID is never permission, cursor or scope.
		result, err := Project(r.Context(), owner, Request{
			Principal: scope.Principal, Tenant: scope.Tenant,
			Environment: scope.Environment, View: Incidents,
			FocusSelection: "Incident:" + id,
		})
		if err != nil {
			if err == ErrCursor {
				http.NotFound(w, r)
			} else {
				http.Error(w, "incident detail unavailable", http.StatusServiceUnavailable)
			}
			return
		}
		for i := range result.Rows {
			if result.Rows[i].Kind == Incident && result.Rows[i].ID == id {
				result.Selected = &result.Rows[i]
				break
			}
		}
		if result.Selected == nil {
			http.NotFound(w, r)
			return
		}
		if err := enrichSelectedResult(r.Context(), owner, scope, &result); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		// The original Incident owner alone can release its canonical header.
		// An unavailable optional provider produces an explicit partial view;
		// a malformed or failed provider aborts before any HTML is emitted.
		overview, err := readAuthorizedIncidentOverview(r.Context(), owner, scope, id)
		if err != nil {
			http.Error(w, "incident detail source unavailable", http.StatusServiceUnavailable)
			return
		}
		if r.Context().Err() != nil {
			http.Error(w, "incident detail unavailable", http.StatusServiceUnavailable)
			return
		}
		// Incident.read alone cannot reveal any Alert: recheck relation scope,
		// then independently authorize every Alert.read before HTML output.
		alertOrigins, err := readAuthorizedIncidentAlertOrigins(r.Context(), owner, scope, id)
		if err != nil {
			http.Error(w, "incident alert origin unavailable", http.StatusServiceUnavailable)
			return
		}
		// The link is not permission. The Govern owner must approve every
		// Result.read, revalidate its identity and return only bounded fields.
		governResults, err := readAuthorizedIncidentGovernResults(r.Context(),
			links, govern, scope, id)
		if err != nil {
			http.Error(w, "incident result source unavailable", http.StatusServiceUnavailable)
			return
		}
		data := incidentDetailView{Result: result,
			Overview: overview,
			AlertOrigins: alertOrigins,
			GovernResults: governResults,
			Factors: selectedFactorLines(result.SelectedFactors),
			ReturnHref: "/command/work-queue?restore=1&selected=" + url.QueryEscape("Incident:"+id) + "&view=incidents",
		}
		var page bytes.Buffer
		if err := incidentDetailPage.Execute(&page, data); err != nil {
			http.Error(w, "incident detail unavailable", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(page.Bytes())
	})
}
