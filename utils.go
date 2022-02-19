package handlerman

import (
	"fmt"
	"reflect"

	"github.com/ksaucedo002/answer/errores"
	"github.com/labstack/echo/v4"
)

func jsonBind(c echo.Context, payload interface{}) error {
	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err != nil {
		return errores.NewBadRequestf(nil, errores.ErrInvalidJSON)
	}
	return nil
}
func (h *HandlerMan) getIndentifierValues(sf reflect.Value) (interface{}, error) {
	switch sf.Type().Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		n, ok := (sf.Interface()).(int)
		if !ok {
			err := fmt.Errorf("%s, error assertion int", sf.Type().Name())
			return nil, errores.NewInternalf(err, errores.ErrDatabaseInternal)
		}
		if n == 0 {
			return nil, errores.NewNotFoundf(nil, "error identificador nulo")
		}
		return n, nil
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		n, ok := (sf.Interface()).(uint)
		if !ok {
			err := fmt.Errorf("%s, error assertion uint", sf.Type().Name())
			return nil, errores.NewInternalf(err, errores.ErrDatabaseInternal)
		}
		if n == 0 {
			return nil, errores.NewNotFoundf(nil, "error identificador nulo")
		}
		return n, nil
	case reflect.String:
		n, ok := (sf.Interface()).(string)
		if !ok {
			err := fmt.Errorf("%s, error assertion string", sf.Type().Name())
			return nil, errores.NewInternalf(err, errores.ErrDatabaseInternal)
		}
		if n == "" {
			return nil, errores.NewNotFoundf(nil, "error identificador nulo")
		}
		return n, nil
	default:
		err := fmt.Errorf("%s, tipo de dato incorrecot", sf.Type().Name())
		return nil, errores.NewInternalf(err, errores.ErrDatabaseInternal)
	}
}
