#!/bin/bash
# scripts/dev/check_debug_prints.sh
# Falha se encontrar prints de debug proibidos no código Go

set -e

if grep -rE 'fmt\.Print|log\.Print|println|print\(' --exclude-dir=.git --exclude-dir=vendor --exclude='scripts/dev/check_debug_prints.sh' .; then
  echo "ERRO: Encontrado uso proibido de prints de debug!" >&2
  exit 1
else
  echo "Nenhum print de debug proibido encontrado."
fi
