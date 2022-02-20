# HandlerMan
Handlerman, es una librería a la cual le proporcionamos un objeto y esta crea las principales 
operaciones CRUD.

## Dependencias:
- Gorm Connextion
- Echo Router

## Operaciones
- Create: `POST`
    - `localhost:91/pasepath`

- Find All: `GET`
    - `localhost:91/pasepath`

- Find by Identifier: `GET`
    - `localhost:91/pasepath/:identifier`

- UPDATE: `PUT`   : _id o identifier es necesarios_
    - `localhost:91/pasepath`

- DELETE: `DELETE`
    - `localhost:91/pasepath/:identifier`


    ## Example
```go
package main

import (
	"log"

	"github.com/ksaucedo002/handlerman"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)
type Personas struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
	Edad   uint   `json:"edad"`
}
type Product struct {
	Codigo string  `json:"codigo"`
	Nombre string  `json:"nombre"`
	Precio float64 `json:"precio"`
}

func main() {
	dsn := "host=localhost user=postgres password=kevin002 dbname=prueba port=5432 sslmode=disable"
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		}})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Database opened!!")
	e := echo.New()



	h := handlerman.NewHandlerMan(e.Group("/persona"), conn)
	h.Start(Personas{})



	handlerProduct := handlerman.NewHandlerMan(e.Group("/producto"), conn)
    ///Nombre del campo primary key, en la tabla y en el modelo
	handlerProduct.Start(Product{}, handlerman.WithKeyFieldName("codigo", "Codigo", false))

	
	e.Start("localhost:91")
}
```