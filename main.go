package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	// Configura los detalles de conexión
	dsn := "root:@tcp(127.0.0.1:3306)/finalprogramacionavanzada"
	var err error
	// Abre la conexión
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Verifica la conexión
	err = db.Ping()
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}

	fmt.Println("Conexión exitosa a la base de datos!")
	usuarios, err := usuarioPorNombre("Dante")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Usuario encontrado: %v\n", usuarios)
	user, err := usuarioPorId(5)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Usuario por id encontrado: %v\n", user)
	userId, err := insertarUsuario(Usuario{
		Nombre: "Dante",
		Mail:   "pruebago@gmail.com",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("El usuario agregado tiene el id: %v\n", userId)
}

type Usuario struct {
	ID     int64
	Nombre string
	Mail   string
}
type Libro struct {
	ID     int64
	Titulo string
	Autor  string
	Dispo  int32
}

func usuarioPorNombre(nombre string) ([]Usuario, error) {
	var usuarios []Usuario
	rows, err := db.Query("SELECT * from usuarios where nombre = ?", nombre)
	if err != nil {
		return nil, fmt.Errorf("usuarioPorNombre %q: %v", nombre, err)
	}
	defer rows.Close()
	for rows.Next() {
		var user Usuario
		if err := rows.Scan(&user.ID, &user.Nombre, &user.Mail); err != nil {
			return nil, fmt.Errorf("usuarioPorNombre %q: %v", nombre, err)
		}
		usuarios = append(usuarios, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("usuarioPorNombre %q: %v", nombre, err)
	}
	return usuarios, nil
}
func usuarioPorId(id int32) (Usuario, error) {
	var usuarios Usuario
	row := db.QueryRow("Select * from usuarios where id = ?", id)
	if err := row.Scan(&usuarios.ID, &usuarios.Nombre, &usuarios.Mail); err != nil {
		if err == sql.ErrNoRows {
			return usuarios, fmt.Errorf("usuarioPorId %d", id)
		}
		return usuarios, fmt.Errorf("usuariosPorId %d: %v", id, err)
	}
	return usuarios, nil
}
func insertarUsuario(user Usuario) (int64, error) {
	res, err := db.Exec("INSERT INTO usuarios(nombre , mail)VALUES(?,?)", user.Nombre, user.Mail)
	if err != nil {
		return 0, fmt.Errorf("insertarUsuario: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("insertarUsuario: %v", err)
	}
	return id, nil
}
