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

///////////////////////////////////
type allowActions struct {
	actions map[string]struct{}
}

func WithAllowActions(actions ...string) allowActions {
	aa := allowActions{}
	for _, action := range actions {
		aa.actions[action] = struct{}{}
	}
	return aa
}
func (a allowActions) apply(h *HandlerMan) {
	for action, str := range a.actions {
		h.allowActions[action] = str
	}
}
