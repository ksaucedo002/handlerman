package handlerman

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/user0608/goones/errs"
)

func jsonBind(c echo.Context, payload interface{}) error {
	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err != nil {
		return errs.BadReqf(nil, "body documento invalido")
	}
	return nil
}
func search(collection []string, s string) bool {
	for _, item := range collection {
		if item == s {
			return true
		}
	}
	return false
}
func cameCaseToSnake(s string) string {
	for _, reStr := range []string{`([A-Z]+)([A-Z][a-z])`, `([a-z\d])([A-Z])`} {
		re := regexp.MustCompile(reStr)
		s = re.ReplaceAllString(s, "${1}_${2}")
	}
	return strings.ToLower(s)
}
func getMapJsonFieldNameWithModelFieldName(i interface{}, ignore ...string) map[string]string {
	var responseMap map[string]string = make(map[string]string)
	rType := reflect.TypeOf(i)
	if rType == nil {
		return responseMap
	}
	if rType.Kind() == reflect.Ptr {
		if rType.Elem().Kind() == reflect.Struct {
			rType = rType.Elem()
		} else {
			return responseMap
		}
	}
	for i := 0; i < rType.NumField(); i++ {
		rsf := rType.Field(i)
		jsonvalue := rsf.Tag.Get("json")
		if jsonvalue == "" {
			continue
		}
		jsonvalue = strings.Split(jsonvalue, ",")[0]
		if !search(ignore, jsonvalue) {
			responseMap[jsonvalue] = cameCaseToSnake(rsf.Name)
		}
	}
	return responseMap
}
func getIndentifierValues(sf reflect.Value) (interface{}, error) {
	switch sf.Type().Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		n, ok := (sf.Interface()).(int)
		if !ok {
			err := fmt.Errorf("%s, error assertion int", sf.Type().Name())
			return nil, errs.Internalf(err, errs.ErrDatabase)
		}
		if n == 0 {
			return nil, errs.Bad("error identificador nulo")
		}
		return n, nil
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		n, ok := (sf.Interface()).(uint)
		if !ok {
			err := fmt.Errorf("%s, error assertion uint", sf.Type().Name())
			return nil, errs.Internalf(err, errs.ErrDatabase)
		}
		if n == 0 {
			return nil, errs.Bad("error identificador nulo")
		}
		return n, nil
	case reflect.String:
		n, ok := (sf.Interface()).(string)
		if !ok {
			err := fmt.Errorf("%s, error assertion string", sf.Type().Name())
			return nil, errs.Internalf(err, errs.ErrDatabase)
		}
		if n == "" {
			return nil, errs.Bad("error identificador nulo")
		}
		return n, nil
	default:
		err := fmt.Errorf("%s, tipo de dato incorrecot", sf.Type().Name())
		return nil, errs.Internalf(err, errs.ErrDatabase)
	}
}
