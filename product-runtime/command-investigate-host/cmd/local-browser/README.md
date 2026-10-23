# CMDR — navigateur local Command → Investigate (données synthétiques)

Cette commande exécute réellement les handlers Go **existants** de la Work Queue Command, de la Timeline Investigate et de leur routeur same-origin. Elle est destinée à valider le parcours dans un **vrai navigateur local**, et ne doit **jamais** être exposée sur un réseau, placée derrière un reverse proxy ou alimentée par des secrets et données de clients.

Depuis `product-runtime/command-investigate-host/` :

```bash
go run ./cmd/local-browser
```

Le programme ouvre uniquement un listener TCP4 `127.0.0.1` sur un port aléatoire, puis imprime une URL complète et deux identifiants de démonstration fraîchement générés (32 octets cryptographiques par mot de passe) dans la console du processus. Ouvrir l'URL en navigateur ; le navigateur demande des identifiants **HTTP Basic natifs** et réutilise cette autorisation sur les deux routes same-origin. Utiliser **analyst** pour consulter un Incident lié à un Case, suivre le lien Case et explorer les versions `event-v1` (remplacée) / `event-v2` avec les filtres HTML du lot K. Utiliser **commander** pour prouver que l'accès à Command ne donne pas le droit de lire le Case, même en saisissant l'URL manuellement. Arrêt : Ctrl+C ; un redémarrage génère deux nouveaux secrets.

Sécurité et limites : **demo synthétique seulement**, aucun IAM client ni source SIEM réelle, aucun TLS (HTTP Basic autorisé uniquement sur loopback), aucune écriture, aucun audit persistant, aucun stockage, aucune API payante. La route refuse les autres Host, les identités incorrectes et les privilèges implicites. Les mots de passe ne sont jamais placés dans l'URL, les cookies, la réponse HTML ni les logs de requêtes ; la console locale est considérée comme un secret temporaire. Un navigateur partagé peut conserver les credentials Basic jusqu'à fermeture de session. Ne pas utiliser comme modèle de déploiement sécurisé ou preuve de production.
