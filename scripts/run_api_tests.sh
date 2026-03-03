#!/bin/sh
# Executa todos os testes de integração da API
set -e
for f in test/*.json; do
  echo "Atualizando timestamp em $f"
  # Atualiza todos os campos "timestamp" para o horário UTC atual
  tmpf=$(mktemp)
  jq --arg now "$(date -u +%Y-%m-%dT%H:%M:%SZ)" '(.messages[] | select(.timestamp != null) | .timestamp) |= $now' "$f" > "$tmpf"
  mv "$tmpf" "$f"
done

for f in test/*.json; do
  echo "Testando $f"
  curl -s -X POST http://localhost:8080/analyze-feed -H 'Content-Type: application/json' -d @$f | jq .
done
