package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
)

func BindQueryContext(payload string) (*QueryContext, error) {
	qc := &QueryContext{}

	if err := json.Unmarshal([]byte(payload), qc); err != nil {
		return nil, err
	}

	return qc, qc.Validate()
}

type QueryContext struct {
	Wheres        []*where        `json:"and_filters"`
	OrWheres      []*orWhere      `json:"or_filters"`
	WhereIns      []*whereIn      `json:"in_filters"`
	WhereNulls    []*whereNull    `json:"null_filters"`
	WhereNotNulls []*whereNotNull `json:"not_null_filters"`
	WhereNotIns   []*whereNotIn   `json:"not_in_filters"`
	Groups        []*group        `json:"groups"`
	Orders        []*order        `json:"orders"`
	Includes      []string        `json:"includes"`
	Limit         int             `json:"limit"`
	Offset        int             `json:"offset"`
}

func (qc *QueryContext) Where(field string, operator string, value interface{}) *QueryContext {
	if len(qc.Wheres) == 0 {
		qc.Wheres = []*where{}
	}

	qc.Wheres = append(qc.Wheres, &where{Field: field, Operator: operator, Value: value})

	return qc
}

func (qc *QueryContext) OrWhere(field string, operator string, value interface{}) *QueryContext {
	if len(qc.OrWheres) == 0 {
		qc.OrWheres = []*orWhere{}
	}

	qc.OrWheres = append(qc.OrWheres, &orWhere{Field: field, Operator: operator, Value: value})

	return qc
}

func (qc *QueryContext) WhereNull(field string) *QueryContext {
	if len(qc.WhereNulls) == 0 {
		qc.WhereNulls = []*whereNull{}
	}

	qc.WhereNulls = append(qc.WhereNulls, &whereNull{Field: field})

	return qc
}

func (qc *QueryContext) WhereNotNull(field string) *QueryContext {
	if len(qc.WhereNotNulls) == 0 {
		qc.WhereNotNulls = []*whereNotNull{}
	}

	qc.WhereNotNulls = append(qc.WhereNotNulls, &whereNotNull{Field: field})

	return qc
}

func (qc *QueryContext) WhereIn(field string, values []interface{}) *QueryContext {
	if len(qc.WhereIns) == 0 {
		qc.WhereIns = []*whereIn{}
	}

	qc.WhereIns = append(qc.WhereIns, &whereIn{Field: field, Values: values})

	return qc
}

func (qc *QueryContext) WhereNotIn(field string, values []interface{}) *QueryContext {
	if len(qc.WhereNotIns) == 0 {
		qc.WhereNotIns = []*whereNotIn{}
	}

	qc.WhereNotIns = append(qc.WhereNotIns, &whereNotIn{Field: field, Values: values})

	return qc
}

func (qc *QueryContext) OrderBy(field, direction string) *QueryContext {
	if len(qc.Orders) == 0 {
		qc.Orders = []*order{}
	}

	qc.Orders = append(qc.Orders, &order{Field: field, Direction: direction})

	return qc
}

func (qc *QueryContext) Size(limit int) *QueryContext {
	qc.Limit = limit

	return qc
}

func (qc *QueryContext) Page(offset int) *QueryContext {
	qc.Offset = offset

	return qc
}

func (qc *QueryContext) Load(includes ...string) *QueryContext {
	qc.Includes = includes

	return qc
}

func (qc *QueryContext) Validate() error {
	for _, where := range qc.Wheres {
		if err := where.validate(); err != nil {
			return err
		}
	}

	for _, orWhere := range qc.OrWheres {
		if err := orWhere.validate(); err != nil {
			return err
		}
	}

	for _, whereIn := range qc.WhereIns {
		if err := whereIn.validate(); err != nil {
			return err
		}
	}

	for _, whereNull := range qc.WhereNulls {
		if err := whereNull.validate(); err != nil {
			return err
		}
	}

	for _, whereNotNull := range qc.WhereNotNulls {
		if err := whereNotNull.validate(); err != nil {
			return err
		}
	}

	for _, whereNotIn := range qc.WhereNotIns {
		if err := whereNotIn.validate(); err != nil {
			return err
		}
	}

	for _, group := range qc.Groups {
		if err := group.validate(); err != nil {
			return err
		}
	}

	for _, order := range qc.Orders {
		if err := order.validate(); err != nil {
			return err
		}
	}

	return nil
}

type where struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

func (w *where) validate() error {
	if len(w.Field) == 0 {
		return errors.New("field cannot be empty for where clause")
	}

	if !inStringSlice(w.Operator, []string{"=", ">", "<", "=>", "<=", "!=", "like"}) {
		return errors.New("invalid operator for where clause")
	}

	if w.Value == nil {
		return errors.New("value cannot be nil for where clause")
	}

	return nil
}

func (w *where) Prepare() (interface{}, interface{}) {
	value := w.Value

	if w.Operator == "like" {
		value = "%" + fmt.Sprint(w.Value) + "%"
	}

	return fmt.Sprintf("%v %v ?", w.Field, w.Operator), value
}

type orWhere struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

func (ow *orWhere) validate() error {
	if len(ow.Field) == 0 {
		return errors.New("field cannot be empty for orWhere clause")
	}

	if !inStringSlice(ow.Operator, []string{"=", ">", "<", "=>", "<=", "!=", "like"}) {
		return errors.New("invalid operator for orWhere clause")
	}

	if ow.Value == nil {
		return errors.New("value cannot be nil for orWhere clause")
	}

	return nil
}

func (ow *orWhere) Prepare() (interface{}, interface{}) {
	value := ow.Field

	if ow.Operator == "like" {
		value = "%" + fmt.Sprint(ow.Value) + "%"
	}

	return fmt.Sprintf("%v %v ?", ow.Field, ow.Operator), value
}

type whereNull struct {
	Field string `json:"field"`
}

func (w *whereNull) validate() error {
	if len(w.Field) == 0 {
		return errors.New("field cannot be empty for whereNull clause")
	}

	return nil
}

func (w *whereNull) Prepare() interface{} {
	return fmt.Sprintf("%v IS NULL", w.Field)
}

type whereNotNull struct {
	Field string `json:"field"`
}

func (w *whereNotNull) validate() error {
	if len(w.Field) == 0 {
		return errors.New("field cannot be empty for whereNotNull clause")
	}

	return nil
}

func (w *whereNotNull) Prepare() interface{} {
	return fmt.Sprintf("%v IS NOT NULL", w.Field)
}

type whereIn struct {
	Field  string        `json:"field"`
	Values []interface{} `json:"values"`
}

func (w *whereIn) validate() error {
	if len(w.Field) == 0 {
		return errors.New("field cannot be empty for whereIn clause")
	}

	if len(w.Values) == 0 {
		return errors.New("values cannot be empty for whereIn clause")
	}

	return nil
}

func (w *whereIn) Prepare() (interface{}, interface{}) {
	return w.Field, w.Values
}

type whereNotIn struct {
	Field  string        `json:"field"`
	Values []interface{} `json:"values"`
}

func (w *whereNotIn) validate() error {
	if len(w.Field) == 0 {
		return errors.New("field cannot be empty for whereNotIn clause")
	}

	if len(w.Values) == 0 {
		return errors.New("values cannot be empty for whereNotIn clause")
	}

	return nil
}

func (w *whereNotIn) Prepare() (interface{}, interface{}) {
	return w.Field, w.Values
}

type group struct {
	Field string `json:"field"`
}

func (g *group) validate() error {
	if len(g.Field) == 0 {
		return errors.New("field cannot be empty for group clause")
	}

	return nil
}

type order struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

func (o *order) validate() error {
	if len(o.Field) == 0 {
		return errors.New("field cannot be empty for order clause")
	}

	if !inStringSlice(o.Direction, []string{"asc", "desc"}) {
		return errors.New("invalid sort direction specified for order clause")
	}

	return nil
}
