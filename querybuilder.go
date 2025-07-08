package eloquent

import (
	"fmt"
	"strings"
)

// QueryBuilder provides fluent query building interface
type QueryBuilder struct {
	connection  *Connection
	table       string
	wheres      []WhereClause
	orders      []OrderClause
	joins       []JoinClause
	groups      []string
	havings     []HavingClause
	limitValue  *int
	offsetValue *int
	columns     []string
	distinct    bool

	// For relations
	eagerLoad map[string]func(*QueryBuilder)
}

// WhereClause represents a where condition
type WhereClause struct {
	Column   string
	Operator string
	Value    interface{}
	Boolean  string        // "and" or "or"
	Type     string        // "basic", "in", "null", "between", "exists", "raw"
	Values   []interface{} // for IN clauses
}

// OrderClause represents an order by clause
type OrderClause struct {
	Column    string
	Direction string
}

// JoinClause represents a join
type JoinClause struct {
	Table    string
	First    string
	Operator string
	Second   string
	Type     string // "inner", "left", "right", "cross"
}

// HavingClause represents a having condition
type HavingClause struct {
	Column   string
	Operator string
	Value    interface{}
	Boolean  string
}

// NewQueryBuilder creates a new query builder
func NewQueryBuilder(connection *Connection) *QueryBuilder {
	return &QueryBuilder{
		connection: connection,
		eagerLoad:  make(map[string]func(*QueryBuilder)),
		columns:    []string{"*"},
	}
}

// Table sets the table name
func (qb *QueryBuilder) Table(table string) *QueryBuilder {
	qb.table = table
	return qb
}

// Select specifies columns to select
func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	qb.columns = columns
	return qb
}

// Distinct adds distinct clause
func (qb *QueryBuilder) Distinct() *QueryBuilder {
	qb.distinct = true
	return qb
}

// Where adds a basic where clause
func (qb *QueryBuilder) Where(column string, args ...interface{}) *QueryBuilder {
	return qb.addWhere(column, "and", args...)
}

// OrWhere adds an OR where clause
func (qb *QueryBuilder) OrWhere(column string, args ...interface{}) *QueryBuilder {
	return qb.addWhere(column, "or", args...)
}

// WhereIn adds a where in clause
func (qb *QueryBuilder) WhereIn(column string, values []interface{}) *QueryBuilder {
	qb.wheres = append(qb.wheres, WhereClause{
		Column:  column,
		Type:    "in",
		Values:  values,
		Boolean: "and",
	})
	return qb
}

// WhereNotIn adds a where not in clause
func (qb *QueryBuilder) WhereNotIn(column string, values []interface{}) *QueryBuilder {
	qb.wheres = append(qb.wheres, WhereClause{
		Column:   column,
		Operator: "not in",
		Type:     "in",
		Values:   values,
		Boolean:  "and",
	})
	return qb
}

// WhereNull adds a where null clause
func (qb *QueryBuilder) WhereNull(column string) *QueryBuilder {
	qb.wheres = append(qb.wheres, WhereClause{
		Column:  column,
		Type:    "null",
		Boolean: "and",
	})
	return qb
}

// WhereNotNull adds a where not null clause
func (qb *QueryBuilder) WhereNotNull(column string) *QueryBuilder {
	qb.wheres = append(qb.wheres, WhereClause{
		Column:   column,
		Operator: "not null",
		Type:     "null",
		Boolean:  "and",
	})
	return qb
}

// WhereBetween adds a where between clause
func (qb *QueryBuilder) WhereBetween(column string, min, max interface{}) *QueryBuilder {
	qb.wheres = append(qb.wheres, WhereClause{
		Column:  column,
		Type:    "between",
		Values:  []interface{}{min, max},
		Boolean: "and",
	})
	return qb
}

// WhereDate adds a where date clause
func (qb *QueryBuilder) WhereDate(column string, operator string, value interface{}) *QueryBuilder {
	return qb.Where(fmt.Sprintf("DATE(%s)", column), operator, value)
}

// WhereTime adds a where time clause
func (qb *QueryBuilder) WhereTime(column string, operator string, value interface{}) *QueryBuilder {
	return qb.Where(fmt.Sprintf("TIME(%s)", column), operator, value)
}

// WhereYear adds a where year clause
func (qb *QueryBuilder) WhereYear(column string, operator string, value interface{}) *QueryBuilder {
	return qb.Where(fmt.Sprintf("YEAR(%s)", column), operator, value)
}

// WhereMonth adds a where month clause
func (qb *QueryBuilder) WhereMonth(column string, operator string, value interface{}) *QueryBuilder {
	return qb.Where(fmt.Sprintf("MONTH(%s)", column), operator, value)
}

// WhereDay adds a where day clause
func (qb *QueryBuilder) WhereDay(column string, operator string, value interface{}) *QueryBuilder {
	return qb.Where(fmt.Sprintf("DAY(%s)", column), operator, value)
}

// Join adds an inner join
func (qb *QueryBuilder) Join(table, first, operator, second string) *QueryBuilder {
	qb.joins = append(qb.joins, JoinClause{
		Table:    table,
		First:    first,
		Operator: operator,
		Second:   second,
		Type:     "inner",
	})
	return qb
}

// LeftJoin adds a left join
func (qb *QueryBuilder) LeftJoin(table, first, operator, second string) *QueryBuilder {
	qb.joins = append(qb.joins, JoinClause{
		Table:    table,
		First:    first,
		Operator: operator,
		Second:   second,
		Type:     "left",
	})
	return qb
}

// RightJoin adds a right join
func (qb *QueryBuilder) RightJoin(table, first, operator, second string) *QueryBuilder {
	qb.joins = append(qb.joins, JoinClause{
		Table:    table,
		First:    first,
		Operator: operator,
		Second:   second,
		Type:     "right",
	})
	return qb
}

// CrossJoin adds a cross join
func (qb *QueryBuilder) CrossJoin(table string) *QueryBuilder {
	qb.joins = append(qb.joins, JoinClause{
		Table: table,
		Type:  "cross",
	})
	return qb
}

// OrderBy adds an order by clause
func (qb *QueryBuilder) OrderBy(column, direction string) *QueryBuilder {
	if direction == "" {
		direction = "asc"
	}
	qb.orders = append(qb.orders, OrderClause{
		Column:    column,
		Direction: strings.ToLower(direction),
	})
	return qb
}

// OrderByDesc adds a descending order by clause
func (qb *QueryBuilder) OrderByDesc(column string) *QueryBuilder {
	return qb.OrderBy(column, "desc")
}

// Latest orders by created_at desc
func (qb *QueryBuilder) Latest(column ...string) *QueryBuilder {
	col := "created_at"
	if len(column) > 0 {
		col = column[0]
	}
	return qb.OrderByDesc(col)
}

// Oldest orders by created_at asc
func (qb *QueryBuilder) Oldest(column ...string) *QueryBuilder {
	col := "created_at"
	if len(column) > 0 {
		col = column[0]
	}
	return qb.OrderBy(col, "asc")
}

// GroupBy adds group by columns
func (qb *QueryBuilder) GroupBy(columns ...string) *QueryBuilder {
	qb.groups = append(qb.groups, columns...)
	return qb
}

// Having adds a having clause
func (qb *QueryBuilder) Having(column, operator string, value interface{}) *QueryBuilder {
	qb.havings = append(qb.havings, HavingClause{
		Column:   column,
		Operator: operator,
		Value:    value,
		Boolean:  "and",
	})
	return qb
}

// OrHaving adds an OR having clause
func (qb *QueryBuilder) OrHaving(column, operator string, value interface{}) *QueryBuilder {
	qb.havings = append(qb.havings, HavingClause{
		Column:   column,
		Operator: operator,
		Value:    value,
		Boolean:  "or",
	})
	return qb
}

// Limit sets the limit
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limitValue = &limit
	return qb
}

// Take is an alias for Limit
func (qb *QueryBuilder) Take(limit int) *QueryBuilder {
	return qb.Limit(limit)
}

// Offset sets the offset
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offsetValue = &offset
	return qb
}

// Skip is an alias for Offset
func (qb *QueryBuilder) Skip(offset int) *QueryBuilder {
	return qb.Offset(offset)
}

// With adds eager loading
func (qb *QueryBuilder) With(relations ...string) *QueryBuilder {
	for _, relation := range relations {
		qb.eagerLoad[relation] = nil
	}
	return qb
}

// WithCallback adds eager loading with callback
func (qb *QueryBuilder) WithCallback(relation string, callback func(*QueryBuilder)) *QueryBuilder {
	qb.eagerLoad[relation] = callback
	return qb
}

// Scopes
func (qb *QueryBuilder) When(condition bool, callback func(*QueryBuilder)) *QueryBuilder {
	if condition {
		callback(qb)
	}
	return qb
}

func (qb *QueryBuilder) Unless(condition bool, callback func(*QueryBuilder)) *QueryBuilder {
	if !condition {
		callback(qb)
	}
	return qb
}

// Execution methods

// Get retrieves all records
func (qb *QueryBuilder) Get() ([]map[string]interface{}, error) {
	sql, args := qb.ToSQL()
	return qb.connection.Select(sql, args...)
}

// First retrieves the first record
func (qb *QueryBuilder) First() (map[string]interface{}, error) {
	qb.Limit(1)
	results, err := qb.Get()
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no records found")
	}
	return results[0], nil
}

// FirstOrFail retrieves the first record or fails
func (qb *QueryBuilder) FirstOrFail() (map[string]interface{}, error) {
	result, err := qb.First()
	if err != nil {
		return nil, fmt.Errorf("model not found")
	}
	return result, nil
}

// Find finds a record by primary key
func (qb *QueryBuilder) Find(id interface{}) (map[string]interface{}, error) {
	return qb.Where("id", id).First()
}

// FindOrFail finds a record by primary key or fails
func (qb *QueryBuilder) FindOrFail(id interface{}) (map[string]interface{}, error) {
	result, err := qb.Find(id)
	if err != nil {
		return nil, fmt.Errorf("model not found")
	}
	return result, nil
}

// Count returns the count of records
func (qb *QueryBuilder) Count(columns ...string) (int64, error) {
	column := "*"
	if len(columns) > 0 {
		column = columns[0]
	}

	countQB := qb.clone()
	countQB.columns = []string{fmt.Sprintf("COUNT(%s) as count", column)}
	countQB.orders = nil
	countQB.limitValue = nil
	countQB.offsetValue = nil

	result, err := countQB.First()
	if err != nil {
		return 0, err
	}

	if count, ok := result["count"].(int64); ok {
		return count, nil
	}

	return 0, fmt.Errorf("invalid count result")
}

// Exists checks if any records exist
func (qb *QueryBuilder) Exists() (bool, error) {
	count, err := qb.Count()
	return count > 0, err
}

// DoesntExist checks if no records exist
func (qb *QueryBuilder) DoesntExist() (bool, error) {
	exists, err := qb.Exists()
	return !exists, err
}

// Paginate returns paginated results
func (qb *QueryBuilder) Paginate(page, perPage int) (*PaginationResult, error) {
	total, err := qb.Count()
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * perPage
	results, err := qb.Offset(offset).Limit(perPage).Get()
	if err != nil {
		return nil, err
	}

	return &PaginationResult{
		Data:        results,
		Total:       total,
		PerPage:     int64(perPage),
		CurrentPage: int64(page),
		LastPage:    (total + int64(perPage) - 1) / int64(perPage),
		From:        int64(offset + 1),
		To:          int64(offset + len(results)),
	}, nil
}

// PaginationResult holds pagination data
type PaginationResult struct {
	Data        []map[string]interface{} `json:"data"`
	Total       int64                    `json:"total"`
	PerPage     int64                    `json:"per_page"`
	CurrentPage int64                    `json:"current_page"`
	LastPage    int64                    `json:"last_page"`
	From        int64                    `json:"from"`
	To          int64                    `json:"to"`
}

// Aggregate methods
func (qb *QueryBuilder) Sum(column string) (float64, error) {
	sumQB := qb.clone()
	sumQB.columns = []string{fmt.Sprintf("SUM(%s) as sum", column)}

	result, err := sumQB.First()
	if err != nil {
		return 0, err
	}

	if sum, ok := result["sum"].(float64); ok {
		return sum, nil
	}

	return 0, nil
}

func (qb *QueryBuilder) Avg(column string) (float64, error) {
	avgQB := qb.clone()
	avgQB.columns = []string{fmt.Sprintf("AVG(%s) as avg", column)}

	result, err := avgQB.First()
	if err != nil {
		return 0, err
	}

	if avg, ok := result["avg"].(float64); ok {
		return avg, nil
	}

	return 0, nil
}

func (qb *QueryBuilder) Max(column string) (interface{}, error) {
	maxQB := qb.clone()
	maxQB.columns = []string{fmt.Sprintf("MAX(%s) as max", column)}

	result, err := maxQB.First()
	if err != nil {
		return nil, err
	}

	return result["max"], nil
}

func (qb *QueryBuilder) Min(column string) (interface{}, error) {
	minQB := qb.clone()
	minQB.columns = []string{fmt.Sprintf("MIN(%s) as min", column)}

	result, err := minQB.First()
	if err != nil {
		return nil, err
	}

	return result["min"], nil
}

// Helper methods
func (qb *QueryBuilder) addWhere(column, boolean string, args ...interface{}) *QueryBuilder {
	var operator string = "="
	var value interface{}

	switch len(args) {
	case 1:
		value = args[0]
	case 2:
		operator = args[0].(string)
		value = args[1]
	default:
		panic("Invalid number of arguments for where clause")
	}

	qb.wheres = append(qb.wheres, WhereClause{
		Column:   column,
		Operator: operator,
		Value:    value,
		Boolean:  boolean,
		Type:     "basic",
	})

	return qb
}

func (qb *QueryBuilder) clone() *QueryBuilder {
	clone := &QueryBuilder{
		connection: qb.connection,
		table:      qb.table,
		wheres:     make([]WhereClause, len(qb.wheres)),
		orders:     make([]OrderClause, len(qb.orders)),
		joins:      make([]JoinClause, len(qb.joins)),
		groups:     make([]string, len(qb.groups)),
		havings:    make([]HavingClause, len(qb.havings)),
		columns:    make([]string, len(qb.columns)),
		distinct:   qb.distinct,
		eagerLoad:  make(map[string]func(*QueryBuilder)),
	}

	copy(clone.wheres, qb.wheres)
	copy(clone.orders, qb.orders)
	copy(clone.joins, qb.joins)
	copy(clone.groups, qb.groups)
	copy(clone.havings, qb.havings)
	copy(clone.columns, qb.columns)

	if qb.limitValue != nil {
		val := *qb.limitValue
		clone.limitValue = &val
	}

	if qb.offsetValue != nil {
		val := *qb.offsetValue
		clone.offsetValue = &val
	}

	for k, v := range qb.eagerLoad {
		clone.eagerLoad[k] = v
	}

	return clone
}

// ToSQL converts the query to SQL
func (qb *QueryBuilder) ToSQL() (string, []interface{}) {
	var sql strings.Builder
	var args []interface{}

	// SELECT clause
	sql.WriteString("SELECT ")
	if qb.distinct {
		sql.WriteString("DISTINCT ")
	}
	sql.WriteString(strings.Join(qb.columns, ", "))

	// FROM clause
	sql.WriteString(" FROM ")
	sql.WriteString(qb.table)

	// JOIN clauses
	for _, join := range qb.joins {
		sql.WriteString(" ")
		sql.WriteString(strings.ToUpper(join.Type))
		sql.WriteString(" JOIN ")
		sql.WriteString(join.Table)
		if join.Type != "cross" {
			sql.WriteString(" ON ")
			sql.WriteString(join.First)
			sql.WriteString(" ")
			sql.WriteString(join.Operator)
			sql.WriteString(" ")
			sql.WriteString(join.Second)
		}
	}

	// WHERE clauses
	if len(qb.wheres) > 0 {
		sql.WriteString(" WHERE ")
		for i, where := range qb.wheres {
			if i > 0 {
				sql.WriteString(" ")
				sql.WriteString(strings.ToUpper(where.Boolean))
				sql.WriteString(" ")
			}

			switch where.Type {
			case "basic":
				sql.WriteString(where.Column)
				sql.WriteString(" ")
				sql.WriteString(where.Operator)
				sql.WriteString(" ?")
				args = append(args, where.Value)
			case "in":
				sql.WriteString(where.Column)
				if where.Operator == "not in" {
					sql.WriteString(" NOT IN (")
				} else {
					sql.WriteString(" IN (")
				}
				placeholders := make([]string, len(where.Values))
				for j, val := range where.Values {
					placeholders[j] = "?"
					args = append(args, val)
				}
				sql.WriteString(strings.Join(placeholders, ", "))
				sql.WriteString(")")
			case "null":
				sql.WriteString(where.Column)
				if where.Operator == "not null" {
					sql.WriteString(" IS NOT NULL")
				} else {
					sql.WriteString(" IS NULL")
				}
			case "between":
				sql.WriteString(where.Column)
				sql.WriteString(" BETWEEN ? AND ?")
				args = append(args, where.Values[0], where.Values[1])
			}
		}
	}

	// GROUP BY clause
	if len(qb.groups) > 0 {
		sql.WriteString(" GROUP BY ")
		sql.WriteString(strings.Join(qb.groups, ", "))
	}

	// HAVING clauses
	if len(qb.havings) > 0 {
		sql.WriteString(" HAVING ")
		for i, having := range qb.havings {
			if i > 0 {
				sql.WriteString(" ")
				sql.WriteString(strings.ToUpper(having.Boolean))
				sql.WriteString(" ")
			}
			sql.WriteString(having.Column)
			sql.WriteString(" ")
			sql.WriteString(having.Operator)
			sql.WriteString(" ?")
			args = append(args, having.Value)
		}
	}

	// ORDER BY clause
	if len(qb.orders) > 0 {
		sql.WriteString(" ORDER BY ")
		orderClauses := make([]string, len(qb.orders))
		for i, order := range qb.orders {
			orderClauses[i] = order.Column + " " + strings.ToUpper(order.Direction)
		}
		sql.WriteString(strings.Join(orderClauses, ", "))
	}

	// LIMIT clause
	if qb.limitValue != nil {
		sql.WriteString(" LIMIT ?")
		args = append(args, *qb.limitValue)
	}

	// OFFSET clause
	if qb.offsetValue != nil {
		sql.WriteString(" OFFSET ?")
		args = append(args, *qb.offsetValue)
	}

	return sql.String(), args
}
