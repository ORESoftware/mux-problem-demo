#!/usr/bin/env bash


set -e;

# ng build --prod

rm huru.zip | cat

### !!!!! REMEMBER TO BUMP version !!!!

zip -r huru.zip . \
-x "zip/*"             \
-x "bin/*"             \
-x "pkg/*"             \
-x ".vscode/*"           \
-x ".idea/*"           \
-x "scripts/*"         \
-x ".git/*"            \
-x "*/.git/*"            \
-x .gitignore          
