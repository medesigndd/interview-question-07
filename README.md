# interview-question-07
Product Code Management — Go REST API + Vue 3 + SQLite

## Requirements
- Go 1.25+ (or current stable release)
- Node.js 20+ and npm

## Run
Terminal 1:
```bash
cd backend
go mod tidy
go test ./...
go run ./cmd/api
```
Terminal 2:
```bash
cd frontend
npm install
npm run dev
```
Open http://localhost:5173. Backend runs at http://localhost:8080.

Product code must be exactly 30 characters from A-Z and 0-9. Display format: `XXXXX-XXXXX-XXXXX-XXXXX-XXXXX-XXXXX`. Duplicates return HTTP 409. SQLite file is created at `backend/data/products.db` when running from backend directory.

## API
- GET `/api/v1/products`
- POST `/api/v1/products` JSON `{"code":"ABCDE12345ABCDE12345ABCDE12345"}`
- DELETE `/api/v1/products/{id}`
- GET `/healthz`

## GitHub
After creating the repository, push with:
```bash
git init
git add .
git commit -m "Build product code manager interview solution"
git branch -M main
git remote add origin https://github.com/medesigndd/interview-question-07.git
git push -u origin main
```
