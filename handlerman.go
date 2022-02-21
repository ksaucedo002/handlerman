package handlerman

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/ksaucedo002/answer"
	"github.com/ksaucedo002/answer/errores"
	"github.com/ksaucedo002/kcheck"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const (
	ACTION_CREATE   = "CREATE"
	ACTION_DELETE   = "DELETE"
	ACTION_UPDATE   = "UPDATE"
	ACTION_FIND_ALL = "FIND_ALL"
	ACTION_FIND_BY  = "FIND_BY_IDENTIFIER"
	ACTION_PATCH    = "PATCH"
)

type isIstring bool
type HandlerMan struct {
	fieldKey        fieldName
	allowActions    map[string]struct{}
	translateFields map[string]string
	filtrableFields map[string]isIstring
	//ignoreFiels  map[string]struct{}
	group   *echo.Group
	storage *storage
}

func NewHandlerMan(g *echo.Group, conn *gorm.DB) *HandlerMan {
	return &HandlerMan{
		fieldKey: fieldName{
			TableFieldName: "id",
			ModelFieldName: "ID",
			IsNumber:       true,
		},
		filtrableFields: make(map[string]isIstring),
		allowActions: map[string]struct{}{
			ACTION_CREATE: {}, ACTION_DELETE: {}, ACTION_UPDATE: {},
			ACTION_FIND_ALL: {}, ACTION_FIND_BY: {}, ACTION_PATCH: {}},
		//ignoreFiels: make(map[string]struct{}),
		group: g,
		storage: &storage{
			conn: conn,
		},
	}
}

func (h *HandlerMan) Start(i interface{}, options ...options) error {
	rType := reflect.TypeOf(i)
	if rType.Kind() == reflect.Ptr {
		if rType.Elem().Kind() != reflect.Struct {
			return fmt.Errorf("error, se esperabas una estructura")
		}
	} else if rType.Kind() != reflect.Struct {
		return fmt.Errorf("error, se esperabas una estructura")
	}
	h.storage.rType = rType
	h.translateFields = getMapJsonFieldNameWithModelFieldName(i)
	for _, op := range options {
		op.apply(h)
	}
	h.router()
	return nil
}
func (h *HandlerMan) router() {
	if _, ok := h.allowActions[ACTION_FIND_ALL]; ok {
		h.group.GET("", h.findAll)
	}
	if _, ok := h.allowActions[ACTION_FIND_BY]; ok {
		h.group.GET(fmt.Sprintf("/:%s", h.fieldKey.TableFieldName), h.findByIdentifier)
	}
	if _, ok := h.allowActions[ACTION_CREATE]; ok {
		h.group.POST("", h.create)
	}
	if _, ok := h.allowActions[ACTION_UPDATE]; ok {
		h.group.PUT("", h.update)
	}
	if _, ok := h.allowActions[ACTION_DELETE]; ok {
		h.group.DELETE(fmt.Sprintf("/:%s", h.fieldKey.TableFieldName), h.delete)
	}
}

func (h *HandlerMan) findAll(c echo.Context) error {
	filterValue := c.QueryParam("filter")
	if filterValue != "" {
		splits := strings.Split(filterValue, ",")
		if len(splits) == 2 {
			nameField, ok := h.translateFields[splits[0]]
			var value interface{}
			if !ok {
				return answer.ErrorResponse(c, errores.NewBadRequestf(nil, "%s invalido", splits[0]))
			}
			fl, ok := h.filtrableFields[nameField]
			if !ok {
				return answer.ErrorResponse(c, errores.NewBadRequestf(
					fmt.Errorf("findAll: filed %s invalido", nameField),
					"%s invalido", splits[0]),
				)
			}
			if fl {
				value = splits[1]
			} else {
				num, err := strconv.Atoi(splits[1])
				if err != nil {
					return answer.ErrorResponse(c, errores.NewBadRequestf(nil, "%s invalido, debe ser numerico", splits[0]))
				}
				value = num
			}
			data, err := h.storage.findAllEntiesWithFilter(filter{fieldName: nameField, value: value})
			if err != nil {
				return answer.ErrorResponse(c, err)
			}
			return answer.OK(c, data)
		}

	}
	data, err := h.storage.findAllEnties()
	if err != nil {
		return answer.ErrorResponse(c, err)
	}
	return answer.OK(c, data)
}
func (h *HandlerMan) findByIdentifier(c echo.Context) error {
	identifier := c.Param(h.fieldKey.TableFieldName)
	var key interface{}
	key = identifier
	var err error
	if h.fieldKey.IsNumber {
		key, err = strconv.Atoi(identifier)
		if err != nil {
			return answer.ErrorResponse(c, errores.NewNotFoundf(nil, errores.ErrRecordNotFaund))
		}
	}

	data, serr := h.storage.findByIdentifier(h.fieldKey.TableFieldName, key)
	if serr != nil {
		return answer.ErrorResponse(c, serr)
	}
	return answer.OK(c, data)
}
func (h *HandlerMan) create(c echo.Context) error {
	newObjet := reflect.New(h.storage.rType).Interface()
	if err := jsonBind(c, newObjet); err != nil {
		return answer.JSONErrorResponse(c)
	}
	if err := kcheck.Valid(newObjet); err != nil {
		return answer.ErrorResponse(c, errores.NewBadRequestf(nil, err.Error()))
	}
	if err := h.storage.create(newObjet); err != nil {
		return answer.ErrorResponse(c, err)
	}
	return answer.Message(c, answer.SUCCESS_OPERATION)
}
func (h *HandlerMan) update(c echo.Context) error {
	newObjet := reflect.New(h.storage.rType).Interface()
	if err := jsonBind(c, newObjet); err != nil {
		return answer.JSONErrorResponse(c)
	}
	field := reflect.ValueOf(newObjet).Elem().FieldByName(h.fieldKey.ModelFieldName)
	if !field.IsValid() {
		err := fmt.Errorf("invalid %s no se encontro en la estrucutra %v", h.fieldKey.ModelFieldName, h.storage.rType)
		return answer.ErrorResponse(c, errores.NewInternalf(err, errores.ErrDatabaseInternal))
	}
	if field.IsZero() {
		err := fmt.Errorf("%s no se encontro en la estrucutra %v", h.fieldKey.ModelFieldName, h.storage.rType)
		return answer.ErrorResponse(c, errores.NewInternalf(err, errores.ErrDatabaseInternal))
	}
	pkvalue, err := getIndentifierValues(field)
	if err != nil {
		return answer.ErrorResponse(c, err)
	}
	if err := kcheck.Valid(newObjet); err != nil {
		return answer.ErrorResponse(c, errores.NewBadRequestf(nil, err.Error()))
	}
	if err := h.storage.update(h.fieldKey.TableFieldName, pkvalue, newObjet); err != nil {
		return answer.ErrorResponse(c, err)
	}
	return answer.Message(c, answer.SUCCESS_OPERATION)
}

func (h *HandlerMan) delete(c echo.Context) error {
	identifier := c.Param(h.fieldKey.TableFieldName)
	var key interface{}
	key = identifier
	var err error
	if h.fieldKey.IsNumber {
		key, err = strconv.Atoi(identifier)
		if err != nil {
			return answer.ErrorResponse(c, errores.NewNotFoundf(nil, errores.ErrRecordNotFaund))
		}
	}
	if err := h.storage.delete(h.fieldKey.TableFieldName, key); err != nil {
		return answer.ErrorResponse(c, err)
	}
	return answer.Message(c, answer.SUCCESS_OPERATION)
}
