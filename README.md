# Vigie – Test technique Backend (Go)

### Compiler
```bash
go build -o orders
```
#### Sur Windows, Go génère parfois un fichier sans extension (orders au lieu de orders.exe).
#### Si c’est le cas, vous pouvez simplement le renommer :
```bash
rename-item orders orders.exe
```
### Exécuter le programme
```bash
.\orders orders.json
.\orders -from=2024-11-01 orders.json 
```
### Lancer les tests unitaires
```bash
go test
```
