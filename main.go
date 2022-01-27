package main

import (
	"database/sql"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func ConexionDB() (conexion *sql.DB) {
	driver := "mysql"
	usuario := "root"
	contrasena := ""
	nombreDB := "sistema"

	conexion, err := sql.Open(driver, usuario+":"+contrasena+"@tcp(127.0.0.1)/"+nombreDB)

	if err != nil {
		panic(err.Error())
	}
	return conexion
}

type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {

	http.HandleFunc("/", Index)            // //Enviar a función Index
	http.HandleFunc("/crear", Crear)       //Enviar a función Crear
	http.HandleFunc("/insertar", Insertar) //Enviar a función Insertar
	http.HandleFunc("/borrar", Borrar)     //Enviar a función Borrar

	http.ListenAndServe(":8080", nil)

}

func Index(w http.ResponseWriter, r *http.Request) {

	conexionEstablecida := ConexionDB()
	registros, err := conexionEstablecida.Query("SELECT * FROM empleados") //consultar registros

	if err != nil {
		panic(err.Error())
	}
	empleado := Empleado{}
	arregloEmpleado := []Empleado{}

	for registros.Next() {
		var id int
		var nombre, correo string
		err = registros.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		arregloEmpleado = append(arregloEmpleado, empleado)

		//fmt.Println(arregloEmpleado)
	}
	plantillas.ExecuteTemplate(w, "index", arregloEmpleado) //redirigir a pagina principal index
}

func Crear(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "Hola Mundo")
	plantillas.ExecuteTemplate(w, "crear", nil) //redirigir al formulario crear
}

func Insertar(w http.ResponseWriter, r *http.Request) { //guardar la información de empleados que viene del formulario crear

	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := ConexionDB()
		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO empleados (nombre,correo) VALUES (?,?)")

		if err != nil {
			panic(err.Error())
		}
		insertarRegistros.Exec(nombre, correo)

		http.Redirect(w, r, "/", 301) //redirige a la pagina principal
	}
}

func Borrar(w http.ResponseWriter, r *http.Request) {
	IdEmpleado := r.URL.Query().Get("id")

	conexionEstablecida := ConexionDB()
	borrarRegistros, err := conexionEstablecida.Prepare("DELETE FROM Empleados WHERE id=?")

	if err != nil {
		panic(err.Error())
	}
	borrarRegistros.Exec(IdEmpleado)

	http.Redirect(w, r, "/", 301) //redirige a la pagina principal
}
