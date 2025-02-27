package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)
type Credential struct {
	ID      int
	SiteApp string
	Usuario string
	Senha   string
}

func OpenDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "senhas.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTable() {
	db, err := OpenDatabase()
	if err != nil {
		log.Fatal("Erro ao abrir o banco de dados:", err)
	}
	defer db.Close()

	sqlStmt := `CREATE TABLE IF NOT EXISTS credenciais (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		siteApp TEXT NOT NULL,
		usuario TEXT NOT NULL,
		senha TEXT NOT NULL
	);`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("Erro ao criar tabela:", err)
	}

	fmt.Println("Tabela 'credenciais' verificada/criada com sucesso!")
}

func InsertCredential(siteApp, usuario, senha string) error {
	db, err := OpenDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO credenciais (siteApp, usuario, senha) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(siteApp, usuario, senha)
	if err != nil {
		return err
	}

	fmt.Println("Nova credencial salva com sucesso!")
	return nil
}

func GetCredentials() ([]Credential, error) {
	db, err := sql.Open("sqlite3", "senhas.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, siteApp, usuario, senha FROM credenciais")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credentials []Credential

	for rows.Next() {
		var c Credential
		err = rows.Scan(&c.ID, &c.SiteApp, &c.Usuario, &c.Senha)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, c)
	}

	return credentials, nil
}

func UpdatePassword(id int, senha string) error {
    db, err := OpenDatabase() 
    if err != nil {
        return err
    }
    defer db.Close()

    _, err = db.Exec("UPDATE credenciais SET senha = ? WHERE id = ?", senha, id)
    return err
}

func DeleteCredential(id int) error {
	db, err := OpenDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM credenciais WHERE id = ?", id)
	return err
}
