package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
    dbPath := os.Getenv("SQLITE_DB_PATH")
    if dbPath == "" {
        dbPath = "./data/app.db"
    }
    os.MkdirAll("./data", 0755)

    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (id INTEGER PRIMARY KEY, content TEXT)`)
    if err != nil {
        log.Fatal(err)
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        fmt.Fprintln(w, `<!DOCTYPE html><html lang='ja'><head><meta charset='utf-8'><title>掲示板</title><style>
        body { font-family: sans-serif; background: #f7f7f7; }
        .container { max-width: 500px; margin: 40px auto; background: #fff; border-radius: 8px; box-shadow: 0 2px 8px #0001; padding: 24px; }
        h1 { text-align: center; }
        ul { padding: 0; }
        li { list-style: none; border-bottom: 1px solid #eee; padding: 8px 0; }
        form { display: flex; gap: 8px; margin-bottom: 16px; }
        input[type=text] { flex: 1; padding: 8px; border: 1px solid #ccc; border-radius: 4px; }
        button { padding: 8px 16px; border: none; background: #215EBF; color: #fff; border-radius: 4px; cursor: pointer; }
        button:hover { background: #174a8c; }
        </style></head><body><div class='container'>`)
        fmt.Fprintln(w, `<h1>ハンズオン掲示板 v2.0.0</h1>`)
        fmt.Fprintln(w, `<form method='POST' action='/add'><input type='text' name='msg' placeholder='メッセージを入力'><button type='submit'>投稿</button></form>`)
        fmt.Fprintln(w, `<ul>`)
        rows, err := db.Query("SELECT id, content FROM messages ORDER BY id DESC")
        if err != nil {
            fmt.Fprintf(w, "<li style='color:red;'>%s</li>", err.Error())
        } else {
            defer rows.Close()
            for rows.Next() {
                var id int
                var content string
                rows.Scan(&id, &content)
                fmt.Fprintf(w, "<li><b>#%d</b>: %s</li>", id, content)
            }
        }
        fmt.Fprintln(w, `</ul></div></body></html>`)
    })

    http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
        var msg string
        if r.Method == "POST" {
            r.ParseForm()
            msg = r.FormValue("msg")
        } else {
            msg = r.URL.Query().Get("msg")
        }
        if msg == "" {
            http.Error(w, "msg required", 400)
            return
        }
        _, err := db.Exec("INSERT INTO messages(content) VALUES(?)", msg)
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    fmt.Println("Listening on port", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
