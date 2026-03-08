#!/bin/bash

# Este script compila el juego a WebAssembly (WASM)
# Necesitás tener Go instalado.

echo "Compilando a WASM..."

# Seteamos las variables de entorno para la compilación cruzada
GOOS=js GOARCH=wasm go build -o 7colors.wasm .

echo "¡Listo! Se generó 7colors.wasm"
echo "Para correrlo, necesitás servir el index.html y 7colors.wasm con un servidor web."
echo "Podés usar por ejemplo: npx serve ."
