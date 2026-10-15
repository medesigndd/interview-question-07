package products

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var codePattern = regexp.MustCompile(`^[A-Z0-9]{30}$`)

type Product struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`
	Formatted string    `json:"formattedCode"`
	CreatedAt time.Time `json:"createdAt"`
}
type API struct{ db *sql.DB }

func InitDB(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS products (id INTEGER PRIMARY KEY AUTOINCREMENT, code TEXT NOT NULL UNIQUE CHECK(length(code)=30), created_at TEXT NOT NULL)`)
	return err
}
func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	a := &API{db: db}
	mux.HandleFunc("GET /api/v1/products", a.list)
	mux.HandleFunc("POST /api/v1/products", a.create)
	mux.HandleFunc("DELETE /api/v1/products/{id}", a.delete)
}
func FormatCode(code string) string {
	parts := make([]string, 0, 6)
	for i := 0; i < len(code); i += 5 {
		parts = append(parts, code[i:i+5])
	}
	return strings.Join(parts, "-")
}
func validCode(code string) bool { return codePattern.MatchString(code) }
func (a *API) list(w http.ResponseWriter, r *http.Request) {
	rows, err := a.db.Query(`SELECT id,code,created_at FROM products ORDER BY id DESC`)
	if err != nil {
		writeError(w, 500, "could not list products")
		return
	}
	defer rows.Close()
	result := []Product{}
	for rows.Next() {
		var p Product
		var created string
		if err := rows.Scan(&p.ID, &p.Code, &created); err != nil {
			writeError(w, 500, "could not read product")
			return
		}
		p.Formatted = FormatCode(p.Code)
		p.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
		result = append(result, p)
	}
	if err := rows.Err(); err != nil {
		writeError(w, 500, "could not read products")
		return
	}
	writeJSON(w, 200, result)
}
func (a *API) create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	var input struct {
		Code string `json:"code"`
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&input); err != nil {
		writeError(w, 400, "invalid JSON body")
		return
	}
	input.Code = strings.TrimSpace(input.Code)
	if !validCode(input.Code) {
		writeError(w, 400, "product code must contain exactly 30 uppercase English letters or digits (A-Z, 0-9)")
		return
	}
	now := time.Now().UTC()
	res, err := a.db.Exec(`INSERT INTO products(code,created_at) VALUES(?,?)`, input.Code, now.Format(time.RFC3339Nano))
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			writeError(w, 409, "product code already exists")
			return
		}
		writeError(w, 500, "could not save product")
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		writeError(w, 500, "could not save product")
		return
	}
	writeJSON(w, 201, Product{ID: id, Code: input.Code, Formatted: FormatCode(input.Code), CreatedAt: now})
}
func (a *API) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		writeError(w, 400, "invalid product id")
		return
	}
	res, err := a.db.Exec(`DELETE FROM products WHERE id=?`, id)
	if err != nil {
		writeError(w, 500, "could not delete product")
		return
	}
	n, err := res.RowsAffected()
	if err != nil {
		writeError(w, 500, "could not delete product")
		return
	}
	if n == 0 {
		writeError(w, 404, "product not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
