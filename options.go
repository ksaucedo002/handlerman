package handlerman

import "github.com/labstack/echo/v4"

type options interface {
	apply(*HandlerMan)
}

// ///////////////////////////////
type fieldName struct {
	TableFieldName string
	ModelFieldName string
	IsNumber       bool
}

// define atributos del model
func WithPrimary(modelfieldNamd string, isnumeric bool) fieldName {
	return fieldName{
		TableFieldName: cameCaseToSnake(modelfieldNamd),
		ModelFieldName: modelfieldNamd,
		IsNumber:       isnumeric,
	}
}
func WithPrimaryT(tableFieldName, modelFieldName string, isNumber bool) fieldName {
	return fieldName{
		TableFieldName: tableFieldName,
		ModelFieldName: modelFieldName,
		IsNumber:       isNumber,
	}
}

func (fn fieldName) apply(h *HandlerMan) {
	h.fieldKey.TableFieldName = fn.TableFieldName
	h.fieldKey.ModelFieldName = fn.ModelFieldName
	h.fieldKey.IsNumber = fn.IsNumber
}

// ///////////////////////////////
type allowActions struct {
	actions []string
}

func WithAcctions(actions ...string) allowActions {
	return allowActions{actions: actions}
}
func (a allowActions) apply(h *HandlerMan) {
	for _, action := range a.actions {
		h.allowActions[action] = struct{}{}
	}
}

// ///////////////////////////////
type FilterOption struct {
	FieldTableName string
	IsString       bool
}
type filtrable struct {
	filtrable []FilterOption
}

func WithFieldFilter(tableFields ...FilterOption) filtrable {
	return filtrable{filtrable: tableFields}
}

func (f filtrable) apply(h *HandlerMan) {
	for _, fil := range f.filtrable {
		h.filtrableFields[fil.FieldTableName] = isIstring(fil.IsString)
	}
}

// /
type createSelect []string

func WithCreateSelect(fields ...string) createSelect {
	return fields
}

func (selects createSelect) apply(h *HandlerMan) {
	for _, s := range selects {
		h.createSelects = append(h.createSelects, s)
	}
}

type updateSelect []string

func WithUpdateSelect(fields ...string) updateSelect {
	return fields
}

func (selects updateSelect) apply(h *HandlerMan) {
	for _, s := range selects {
		h.updateSelects = append(h.updateSelects, s)
	}
}

type findSelect []string

func WithFindSelect(fields ...string) findSelect {
	return fields
}

func (selects findSelect) apply(h *HandlerMan) {
	for _, s := range selects {
		h.findSelects = append(h.updateSelects, s)
	}
}

// ///////////
type createMiddlewares []echo.MiddlewareFunc

func WithCreateMiddl(fields ...echo.MiddlewareFunc) createMiddlewares {
	return fields
}

func (selects createMiddlewares) apply(h *HandlerMan) {
	for _, s := range selects {
		h.createdMiddlewares = append(h.createdMiddlewares, s)
	}
}

type updateMiddlewares []echo.MiddlewareFunc

func WithUpdateMiddl(fields ...echo.MiddlewareFunc) updateMiddlewares {
	return fields
}

func (selects updateMiddlewares) apply(h *HandlerMan) {
	for _, s := range selects {
		h.updateMiddlewares = append(h.updateMiddlewares, s)
	}
}

type findMiddlewares []echo.MiddlewareFunc

func WithFindMiddl(fields ...echo.MiddlewareFunc) findMiddlewares {
	return fields
}

func (selects findMiddlewares) apply(h *HandlerMan) {
	for _, s := range selects {
		h.findMiddlewares = append(h.findMiddlewares, s)
	}
}

type deleteMiddleware []echo.MiddlewareFunc

func WithDeleteMiddl(fields ...echo.MiddlewareFunc) deleteMiddleware {
	return fields
}

func (selects deleteMiddleware) apply(h *HandlerMan) {
	for _, s := range selects {
		h.deleteMidlewares = append(h.deleteMidlewares, s)
	}
}
