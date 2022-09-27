package handlerman

import (
	"fmt"
	"reflect"

	"github.com/user0608/goones/errs"
	"gorm.io/gorm"
)

type filter struct {
	fieldName string
	value     interface{}
}
type storage struct {
	rType reflect.Type
	conn  *gorm.DB
}

func (s *storage) findAllEntiesWithFilter(flt filter, selects []string) (interface{}, error) {
	container := reflect.MakeSlice(reflect.SliceOf(s.rType), 0, 0).Interface()
	tx := selectField(s.conn, selects)
	if err := tx.Find(&container, fmt.Sprintf("%s = ?", flt.fieldName), flt.value).Error; err != nil {
		return nil, errs.Pgf(err)
	}
	return container, nil
}
func (s *storage) findAllEnties(selects []string) (interface{}, error) {
	container := reflect.MakeSlice(reflect.SliceOf(s.rType), 0, 0).Interface()
	tx := selectField(s.conn, selects)
	if err := tx.Find(&container).Error; err != nil {
		return nil, errs.Pgf(err)
	}
	return container, nil
}
func (s *storage) findByIdentifier(fileName string, value interface{}, selects []string) (interface{}, error) {
	newObjet := reflect.New(s.rType).Interface()
	tx := selectField(s.conn.Limit(1), selects)
	rs := tx.Where(fmt.Sprintf("%s = ?", fileName), value).Find(&newObjet)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}
	if rs.RowsAffected == 0 {
		return nil, errs.Notfoundf(nil, errs.ErrRecordNotFaund)
	}
	return newObjet, nil
}
func (s *storage) create(i interface{}, selects []string) error {
	tx := selectField(s.conn, selects)
	if err := tx.Create(i).Error; err != nil {
		return errs.Pgf(err)
	}
	return nil
}
func (s *storage) update(pkname string, pkvalue interface{}, i interface{}, selects []string) error {
	tx := selectField(s.conn.Where(fmt.Sprintf("%s = ?", pkname), pkvalue), selects)
	if err := tx.Updates(i).Error; err != nil {
		return errs.Pgf(err)
	}
	return nil
}
func (s *storage) delete(pkname string, pkvalue interface{}) error {
	rs := s.conn.Delete(reflect.New(s.rType).Interface(), fmt.Sprintf("%s = ?", pkname), pkvalue)
	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}
	if rs.RowsAffected == 0 {
		return errs.Bad(errs.ErrRecordNotFaund)
	}
	return nil
}

func selectField(tx *gorm.DB, selects []string) *gorm.DB {
	if len(selects) == 0 {
		return tx
	}
	for _, s := range selects {
		tx = tx.Select(s)
	}
	return tx
}
