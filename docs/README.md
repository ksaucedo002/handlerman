## Example
```go
package main

import (
	"log"

	handlerman "github.com/ksaucedo002/handlermain"
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