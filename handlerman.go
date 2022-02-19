package handlerman

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/ksaucedo002/answer"
	"github.com/ksaucedo002/answer/errores"
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

type HandlerMan struct {
	fieldKey     fieldName
	allowActions map[string]struct{}
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
	if field.IsZero() {
		err := fmt.Errorf("%s no se encontro en la estrucutra %v", h.fieldKey.ModelFieldName, h.storage.rType)
		return answer.ErrorResponse(c, errores.NewInternalf(err, errores.ErrDatabaseInternal))
	}
	pkvalue, err := h.getIndentifierValues(field)
	if err != nil {
		return answer.ErrorResponse(c, err)
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
