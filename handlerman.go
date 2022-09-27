package handlerman

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/ksaucedo002/kcheck"
	"github.com/labstack/echo/v4"
	"github.com/user0608/goones/answer"
	"github.com/user0608/goones/errs"
	"gorm.io/gorm"
)

const (
	ACTION_CREATE   = "CREATE"
	ACTION_DELETE   = "DELETE"
	ACTION_UPDATE   = "UPDATE"
	ACTION_FIND_ALL = "FIND_ALL"
	ACTION_FIND_BY  = "FIND_BY_IDENTIFIER"
)

type isIstring bool
type HandlerMan struct {
	fieldKey        fieldName
	allowActions    map[string]struct{}
	translateFields map[string]string
	filtrableFields map[string]isIstring

	createSelects []string
	updateSelects []string
	findSelects   []string

	createdMiddlewares []echo.MiddlewareFunc
	updateMiddlewares  []echo.MiddlewareFunc
	findMiddlewares    []echo.MiddlewareFunc
	deleteMidlewares   []echo.MiddlewareFunc

	//ignoreFiels  map[string]struct{}
	group   *echo.Group
	storage *storage
}

var once sync.Once
var connection *gorm.DB

func SetConn(db *gorm.DB) {
	once.Do(func() { connection = db })
}
func NewGroup(g *echo.Group) *HandlerMan {
	if connection == nil {
		log.Panic("gorm db connection not found, set connection before create new groups")
	}
	return &HandlerMan{
		fieldKey: fieldName{
			TableFieldName: "id",
			ModelFieldName: "ID",
			IsNumber:       true,
		},
		filtrableFields: make(map[string]isIstring),
		allowActions: map[string]struct{}{
			ACTION_CREATE: {}, ACTION_DELETE: {}, ACTION_UPDATE: {},
			ACTION_FIND_ALL: {}, ACTION_FIND_BY: {}},
		//ignoreFiels: make(map[string]struct{}),
		group: g,
		storage: &storage{
			conn: connection,
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
		h.group.GET("", h.findAll, h.findMiddlewares...)
	}
	if _, ok := h.allowActions[ACTION_FIND_BY]; ok {
		h.group.GET(fmt.Sprintf("/:%s", h.fieldKey.TableFieldName), h.findByIdentifier, h.findMiddlewares...)
	}
	if _, ok := h.allowActions[ACTION_CREATE]; ok {
		h.group.POST("", h.create, h.createdMiddlewares...)
	}
	if _, ok := h.allowActions[ACTION_UPDATE]; ok {
		h.group.PUT("", h.update, h.updateMiddlewares...)
	}
	if _, ok := h.allowActions[ACTION_DELETE]; ok {
		h.group.DELETE(fmt.Sprintf("/:%s", h.fieldKey.TableFieldName), h.delete, h.deleteMidlewares...)
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
				return answer.Err(c, errs.Bad("%s invalido", splits[0]))
			}
			fl, ok := h.filtrableFields[nameField]
			if !ok {
				return answer.Err(c, errs.BadReqf(
					fmt.Errorf("findAll: filed %s invalido", nameField),
					"%s invalido", splits[0]),
				)
			}
			if fl {
				value = splits[1]
			} else {
				num, err := strconv.Atoi(splits[1])
				if err != nil {
					return answer.Err(c, errs.Bad("%s invalido, debe ser numerico", splits[0]))
				}
				value = num
			}
			data, err := h.storage.findAllEntiesWithFilter(filter{fieldName: nameField, value: value}, h.findSelects)
			if err != nil {
				return answer.Err(c, err)
			}
			return answer.Ok(c, data)
		}

	}
	data, err := h.storage.findAllEnties(h.findSelects)
	if err != nil {
		return answer.Err(c, err)
	}
	return answer.Ok(c, data)
}

func (h *HandlerMan) findByIdentifier(c echo.Context) error {
	identifier := c.Param(h.fieldKey.TableFieldName)
	var key interface{}
	key = identifier
	var err error
	if h.fieldKey.IsNumber {
		key, err = strconv.Atoi(identifier)
		if err != nil {
			return answer.Err(c, errs.Notfoundf(nil, errs.ErrRecordNotFaund))
		}
	}

	data, serr := h.storage.findByIdentifier(h.fieldKey.TableFieldName, key, h.findSelects)
	if serr != nil {
		return answer.Err(c, serr)
	}
	return answer.Ok(c, data)
}

func (h *HandlerMan) create(c echo.Context) error {
	newObjet := reflect.New(h.storage.rType).Interface()
	if err := jsonBind(c, newObjet); err != nil {
		return answer.JsonErr(c)
	}
	if err := kcheck.Valid(newObjet); err != nil {
		return answer.Err(c, errs.Bad(err.Error()))
	}
	if err := h.storage.create(newObjet, h.createSelects); err != nil {
		return answer.Err(c, err)
	}
	return answer.Message(c, answer.SUCCESS)
}

func (h *HandlerMan) update(c echo.Context) error {
	newObjet := reflect.New(h.storage.rType).Interface()
	if err := jsonBind(c, newObjet); err != nil {
		return answer.JsonErr(c)
	}
	field := reflect.ValueOf(newObjet).Elem().FieldByName(h.fieldKey.ModelFieldName)
	if !field.IsValid() {
		err := fmt.Errorf("invalid %s no se encontro en la estrucutra %v", h.fieldKey.ModelFieldName, h.storage.rType)
		return answer.Err(c, errs.Internalf(err, errs.ErrDatabase))
	}
	if field.IsZero() {
		err := fmt.Errorf("%s no se encontro en la estrucutra %v", h.fieldKey.ModelFieldName, h.storage.rType)
		return answer.Err(c, errs.Internalf(err, errs.ErrDatabase))
	}
	pkvalue, err := getIndentifierValues(field)
	if err != nil {
		return answer.Err(c, err)
	}
	if err := kcheck.Valid(newObjet); err != nil {
		return answer.Err(c, errs.Bad(err.Error()))
	}
	if err := h.storage.update(h.fieldKey.TableFieldName, pkvalue, newObjet, h.updateSelects); err != nil {
		return answer.Err(c, err)
	}
	return answer.Message(c, answer.SUCCESS)
}

func (h *HandlerMan) delete(c echo.Context) error {
	identifier := c.Param(h.fieldKey.TableFieldName)
	var key interface{}
	key = identifier
	var err error
	if h.fieldKey.IsNumber {
		key, err = strconv.Atoi(identifier)
		if err != nil {
			return answer.Err(c, errs.Notfoundf(nil, errs.ErrRecordNotFaund))
		}
	}
	if err := h.storage.delete(h.fieldKey.TableFieldName, key); err != nil {
		return answer.Err(c, err)
	}
	return answer.Message(c, answer.SUCCESS)
}
