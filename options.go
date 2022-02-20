package handlerman

type options interface {
	apply(*HandlerMan)
}

/////////////////////////////////
type fieldName struct {
	TableFieldName string
	ModelFieldName string
	IsNumber       bool
}

func WithKeyFieldName(tableFieldName, modelFieldName string, isNumber bool) fieldName {
	return fieldName{TableFieldName: tableFieldName, ModelFieldName: modelFieldName, IsNumber: isNumber}
}
func (fn fieldName) apply(h *HandlerMan) {
	h.fieldKey.TableFieldName = fn.TableFieldName
	h.fieldKey.ModelFieldName = fn.ModelFieldName
	h.fieldKey.IsNumber = fn.IsNumber
}

/////////////////////////////////
type allowActions struct {
	actions []string
}

func WithAllowActions(actions ...string) allowActions {
	return allowActions{actions: actions}
}
func (a allowActions) apply(h *HandlerMan) {
	for _, action := range a.actions {
		h.allowActions[action] = struct{}{}
	}
}

/////////////////////////////////
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
