package handlerman

import (
	"fmt"
	"reflect"

	"github.com/ksaucedo002/answer/errores"
	"gorm.io/gorm"
)

type storage struct {
	rType reflect.Type
	conn  *gorm.DB
}

func (s *storage) findAllEnties() (interface{}, error) {
	container := reflect.MakeSlice(reflect.SliceOf(s.rType), 0, 0).Interface()
	if err := s.conn.Find(&container).Error; err != nil {
		return nil, errores.NewInternalDBf(err)
	}
	return container, nil
}
func (s *storage) findByIdentifier(fileName string, value interface{}) (interface{}, error) {
	newObjet := reflect.New(s.rType).Interface()
	tx := s.conn.Limit(1)
	rs := tx.Where(fmt.Sprintf("%s = ?", fileName), value).Find(&newObjet)
	if rs.Error != nil {
		return nil, errores.NewInternalDBf(rs.Error)
	}
	if rs.RowsAffected == 0 {
		return nil, errores.NewNotFoundf(nil, errores.ErrRecordNotFaund)
	}
	return newObjet, nil
}
func (s *storage) create(i interface{}) error {
	if err := s.conn.Create(i).Error; err != nil {
		return errores.NewInternalDBf(err)
	}
	return nil
}
func (s *storage) update(pkname string, pkvalue interface{}, i interface{}) error {
	if err := s.conn.Where(fmt.Sprintf("%s = ?", pkname), pkvalue).Updates(i).Error; err != nil {
		return errores.NewInternalDBf(err)
	}
	return nil
}
func (s *storage) delete(pkname string, pkvalue interface{}) error {
	rs := s.conn.Delete(reflect.New(s.rType).Interface(), fmt.Sprintf("%s = ?", pkname), pkvalue)
	if rs.Error != nil {
		return errores.NewInternalDBf(rs.Error)
	}
	if rs.RowsAffected == 0 {
		return errores.NewBadRequestf(nil, errores.ErrRecordNotFaund)
	}
	return nil
}
