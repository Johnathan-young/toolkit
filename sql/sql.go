package sql

import (
	"fmt"
	"reflect"

	"github.com/doug-martin/goqu/v9"
	"gorm.io/gorm"
)

// QueryOp ...
type QueryOp string

// Op defined
const (
	EQ         QueryOp = "eq"
	NEQ        QueryOp = "neq"
	BETWEEN    QueryOp = "between"
	IN         QueryOp = "in"
	NOTIN      QueryOp = "notIn"
	GT         QueryOp = "gt"
	GTE        QueryOp = "gte"
	LT         QueryOp = "lt"
	LTE        QueryOp = "lte"
	LIKE       QueryOp = "like"
	ILIKE      QueryOp = "iLike"
	NOTLIKE    QueryOp = "notLike"
	IS         QueryOp = "is"
	NOTBETWEEN QueryOp = "notBetween"
)

// ArgsFilter
type ArgsFilter struct {
	ex        goqu.Ex
	filterMap Condition
	exOr      goqu.ExOr
}

// QueryFilter ...
func QueryFilter(filterMap Condition) *ArgsFilter {
	return &ArgsFilter{
		ex:        make(goqu.Ex),
		filterMap: filterMap,
		exOr:      make(goqu.ExOr),
	}
}

// Update ...
func (f *ArgsFilter) Update(field string,
	filterField string) *ArgsFilter {
	if f.filterMap == nil {
		return f
	}

	value, ok := f.filterMap[filterField]
	if ok {
		f.ex[field] = value
	}
	return f
}

// Set ...
func (f *ArgsFilter) Set(field string, value interface{}) *ArgsFilter {
	f.ex[field] = value
	return f
}

// Ex ...
func (f *ArgsFilter) Ex() map[string]interface{} {
	return f.ex
}

// Where ...
func (f *ArgsFilter) Where(field string, op QueryOp,
	filterField string) *ArgsFilter {
	if f.filterMap == nil {
		return f
	}

	value, ok := f.filterMap[filterField]
	if ok {
		if opEx := f.execOp(op, value); opEx != nil {
			f.ex[field] = opEx
		}

	}
	return f
}

// Or ...
func (f *ArgsFilter) Or(field string, op QueryOp, filterField string) *ArgsFilter {
	if f.filterMap == nil {
		return f
	}

	value, ok := f.filterMap[filterField]
	if ok {
		if opEx := f.execOp(op, value); opEx != nil {
			f.exOr[field] = opEx
		}
	}
	return f
}

func (f *ArgsFilter) execOp(op QueryOp, value interface{}) goqu.Op {
	var opEx goqu.Op
	typ := reflect.TypeOf(value)
	switch op {
	case EQ, NEQ, GT, LT, LTE, GTE, IS:
		opEx = goqu.Op{string(op): value}
	case BETWEEN, NOTBETWEEN:
		if typ == nil {
			return opEx
		}
		kind := typ.Kind()
		v := reflect.ValueOf(value)
		if kind == reflect.Slice && v.Len() == 2 {
			opEx = goqu.Op{
				string(op): goqu.Range(v.Index(0).Interface(), v.Index(1).Interface()),
			}
		}
	case IN, NOTIN:
		if typ == nil {
			return opEx
		}
		kind := typ.Kind()
		v := reflect.ValueOf(value)
		if kind == reflect.Slice && v.Len() > 0 {
			opEx = goqu.Op{string(op): value}
		}

	case LIKE, NOTLIKE, ILIKE:
		opEx = goqu.Op{string(op): "%" + fmt.Sprintf("%v", value) + "%"}
	default:
	}

	return opEx
}

// End ...
func (f *ArgsFilter) End() []goqu.Expression {
	var ex []goqu.Expression
	if len(f.ex) > 0 {
		ex = append(ex, f.ex)
	}
	if len(f.exOr) > 0 {
		ex = append(ex, f.exOr)
	}
	return ex
}

// QueryFirst ..
func (db *DataBase) QueryFirst(query *goqu.SelectDataset, scaner *gorm.DB,
	outRows interface{}, selectEx ...interface{}) error {

	selectQuery := query

	if selectEx != nil {
		selectQuery = query.Select(selectEx...)
	}

	sql, args, err := selectQuery.Prepared(true).ToSQL()
	if err != nil {
		return err
	}

	// use gorm to scan
	tx := scaner.Raw(sql, args...).Limit(1).Find(outRows)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
