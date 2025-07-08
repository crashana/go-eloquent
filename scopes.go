package eloquent

import (
	"fmt"
	"strings"
	"time"
)

// Scope represents a query scope function
type Scope func(*QueryBuilder)

// GlobalScope represents a global query scope
type GlobalScope interface {
	Apply(*QueryBuilder, Model)
}

// ScopeRegistry manages query scopes
type ScopeRegistry struct {
	scopes map[string]Scope
	global []GlobalScope
}

// NewScopeRegistry creates a new scope registry
func NewScopeRegistry() *ScopeRegistry {
	return &ScopeRegistry{
		scopes: make(map[string]Scope),
		global: make([]GlobalScope, 0),
	}
}

// Register registers a named scope
func (sr *ScopeRegistry) Register(name string, scope Scope) {
	sr.scopes[name] = scope
}

// RegisterGlobal registers a global scope
func (sr *ScopeRegistry) RegisterGlobal(scope GlobalScope) {
	sr.global = append(sr.global, scope)
}

// Apply applies a named scope to a query builder
func (sr *ScopeRegistry) Apply(name string, qb *QueryBuilder) error {
	if scope, exists := sr.scopes[name]; exists {
		scope(qb)
		return nil
	}
	return fmt.Errorf("scope '%s' not found", name)
}

// ApplyGlobal applies all global scopes to a query builder
func (sr *ScopeRegistry) ApplyGlobal(qb *QueryBuilder, model Model) {
	for _, scope := range sr.global {
		scope.Apply(qb, model)
	}
}

// Common scopes

// ActiveScope filters out soft-deleted records
type ActiveScope struct{}

func (s ActiveScope) Apply(qb *QueryBuilder, model Model) {
	if model.GetDeletedAtColumn() != "" {
		qb.WhereNull(model.GetDeletedAtColumn())
	}
}

// PublishedScope filters for published records
func PublishedScope() Scope {
	return func(qb *QueryBuilder) {
		qb.Where("published", true)
	}
}

// RecentScope filters for recent records
func RecentScope(days int) Scope {
	return func(qb *QueryBuilder) {
		date := time.Now().AddDate(0, 0, -days)
		qb.Where("created_at", ">=", date)
	}
}

// PopularScope orders by popularity
func PopularScope() Scope {
	return func(qb *QueryBuilder) {
		qb.OrderByDesc("views")
	}
}

// SearchScope adds search functionality
func SearchScope(query string, columns ...string) Scope {
	return func(qb *QueryBuilder) {
		if query == "" {
			return
		}

		searchTerm := "%" + strings.ToLower(query) + "%"

		if len(columns) == 0 {
			columns = []string{"name", "title", "description"}
		}

		// Add OR conditions for each column
		for i, column := range columns {
			if i == 0 {
				qb.Where(fmt.Sprintf("LOWER(%s)", column), "LIKE", searchTerm)
			} else {
				qb.OrWhere(fmt.Sprintf("LOWER(%s)", column), "LIKE", searchTerm)
			}
		}
	}
}

// WhereStatusScope filters by status
func WhereStatusScope(status string) Scope {
	return func(qb *QueryBuilder) {
		qb.Where("status", status)
	}
}

// WhereCategoryScope filters by category
func WhereCategoryScope(categoryId interface{}) Scope {
	return func(qb *QueryBuilder) {
		qb.Where("category_id", categoryId)
	}
}

// WhereUserScope filters by user
func WhereUserScope(userId interface{}) Scope {
	return func(qb *QueryBuilder) {
		qb.Where("user_id", userId)
	}
}

// BetweenDatesScope filters records between two dates
func BetweenDatesScope(start, end time.Time, column ...string) Scope {
	return func(qb *QueryBuilder) {
		col := "created_at"
		if len(column) > 0 {
			col = column[0]
		}
		qb.WhereBetween(col, start, end)
	}
}

// WithinDaysScope filters records within specified days
func WithinDaysScope(days int, column ...string) Scope {
	return func(qb *QueryBuilder) {
		col := "created_at"
		if len(column) > 0 {
			col = column[0]
		}
		date := time.Now().AddDate(0, 0, -days)
		qb.Where(col, ">=", date)
	}
}

// LimitScope limits the number of results
func LimitScope(limit int) Scope {
	return func(qb *QueryBuilder) {
		qb.Limit(limit)
	}
}

// OffsetScope sets the offset for results
func OffsetScope(offset int) Scope {
	return func(qb *QueryBuilder) {
		qb.Offset(offset)
	}
}

// OrderScope adds ordering to the query
func OrderScope(column, direction string) Scope {
	return func(qb *QueryBuilder) {
		qb.OrderBy(column, direction)
	}
}

// GroupScope adds grouping to the query
func GroupScope(columns ...string) Scope {
	return func(qb *QueryBuilder) {
		qb.GroupBy(columns...)
	}
}

// HavingScope adds having conditions to the query
func HavingScope(column, operator string, value interface{}) Scope {
	return func(qb *QueryBuilder) {
		qb.Having(column, operator, value)
	}
}

// JoinScope adds joins to the query
func JoinScope(table, first, operator, second string) Scope {
	return func(qb *QueryBuilder) {
		qb.Join(table, first, operator, second)
	}
}

// LeftJoinScope adds left joins to the query
func LeftJoinScope(table, first, operator, second string) Scope {
	return func(qb *QueryBuilder) {
		qb.LeftJoin(table, first, operator, second)
	}
}

// SelectScope specifies columns to select
func SelectScope(columns ...string) Scope {
	return func(qb *QueryBuilder) {
		qb.Select(columns...)
	}
}

// DistinctScope adds distinct to the query
func DistinctScope() Scope {
	return func(qb *QueryBuilder) {
		qb.Distinct()
	}
}

// Conditional scopes

// WhenScope applies a scope conditionally
func WhenScope(condition bool, scope Scope) Scope {
	return func(qb *QueryBuilder) {
		if condition {
			scope(qb)
		}
	}
}

// UnlessScope applies a scope unless condition is true
func UnlessScope(condition bool, scope Scope) Scope {
	return func(qb *QueryBuilder) {
		if !condition {
			scope(qb)
		}
	}
}

// Complex scopes

// PaginateScope adds pagination to the query
func PaginateScope(page, perPage int) Scope {
	return func(qb *QueryBuilder) {
		offset := (page - 1) * perPage
		qb.Offset(offset).Limit(perPage)
	}
}

// FilterScope applies multiple filters based on a map
func FilterScope(filters map[string]interface{}) Scope {
	return func(qb *QueryBuilder) {
		for column, value := range filters {
			if value != nil && value != "" {
				qb.Where(column, value)
			}
		}
	}
}

// DateRangeScope filters by date range
func DateRangeScope(startDate, endDate *time.Time, column ...string) Scope {
	return func(qb *QueryBuilder) {
		col := "created_at"
		if len(column) > 0 {
			col = column[0]
		}

		if startDate != nil {
			qb.Where(col, ">=", *startDate)
		}

		if endDate != nil {
			qb.Where(col, "<=", *endDate)
		}
	}
}

// AdvancedSearchScope provides advanced search functionality
func AdvancedSearchScope(searchParams map[string]interface{}) Scope {
	return func(qb *QueryBuilder) {
		for field, value := range searchParams {
			switch v := value.(type) {
			case string:
				if v != "" {
					qb.Where(field, "LIKE", "%"+v+"%")
				}
			case []interface{}:
				if len(v) > 0 {
					qb.WhereIn(field, v)
				}
			case map[string]interface{}:
				// Handle range queries
				if min, hasMin := v["min"]; hasMin {
					qb.Where(field, ">=", min)
				}
				if max, hasMax := v["max"]; hasMax {
					qb.Where(field, "<=", max)
				}
			default:
				qb.Where(field, value)
			}
		}
	}
}

// SoftDeleteScope handles soft deletes
type SoftDeleteScope struct {
	includeDeleted bool
	onlyDeleted    bool
}

func (s SoftDeleteScope) Apply(qb *QueryBuilder, model Model) {
	deletedAtColumn := model.GetDeletedAtColumn()
	if deletedAtColumn == "" {
		return
	}

	if s.onlyDeleted {
		qb.WhereNotNull(deletedAtColumn)
	} else if !s.includeDeleted {
		qb.WhereNull(deletedAtColumn)
	}
	// If includeDeleted is true, don't add any conditions
}

// WithTrashedScope includes soft-deleted records
func WithTrashedScope() GlobalScope {
	return SoftDeleteScope{includeDeleted: true}
}

// OnlyTrashedScope returns only soft-deleted records
func OnlyTrashedScope() GlobalScope {
	return SoftDeleteScope{onlyDeleted: true}
}

// Model scope methods

// ApplyScope applies a scope to a model's query
func ApplyScope(qb *QueryBuilder, scope Scope) *QueryBuilder {
	scope(qb)
	return qb
}

// ApplyScopes applies multiple scopes to a model's query
func ApplyScopes(qb *QueryBuilder, scopes ...Scope) *QueryBuilder {
	for _, scope := range scopes {
		scope(qb)
	}
	return qb
}

// ChainScopes chains multiple scopes together
func ChainScopes(scopes ...Scope) Scope {
	return func(qb *QueryBuilder) {
		for _, scope := range scopes {
			scope(qb)
		}
	}
}

// Dynamic scope creation

// CreateScope creates a dynamic scope from a function
func CreateScope(name string, fn func(*QueryBuilder, ...interface{})) func(...interface{}) Scope {
	return func(args ...interface{}) Scope {
		return func(qb *QueryBuilder) {
			fn(qb, args...)
		}
	}
}

// Example usage of dynamic scopes:
// var WhereColumnScope = CreateScope("whereColumn", func(qb *QueryBuilder, args ...interface{}) {
//     if len(args) >= 2 {
//         qb.Where(args[0].(string), args[1])
//     }
// })
//
// Usage: query.Apply(WhereColumnScope("name", "John"))

// Scope utilities

// ScopeExists checks if a scope exists in the registry
func (sr *ScopeRegistry) ScopeExists(name string) bool {
	_, exists := sr.scopes[name]
	return exists
}

// ListScopes returns all registered scope names
func (sr *ScopeRegistry) ListScopes() []string {
	var names []string
	for name := range sr.scopes {
		names = append(names, name)
	}
	return names
}

// RemoveScope removes a scope from the registry
func (sr *ScopeRegistry) RemoveScope(name string) {
	delete(sr.scopes, name)
}

// ClearScopes removes all scopes from the registry
func (sr *ScopeRegistry) ClearScopes() {
	sr.scopes = make(map[string]Scope)
}

// Global scope registry
var globalScopeRegistry = NewScopeRegistry()

// RegisterGlobalScope registers a global scope
func RegisterGlobalScope(scope GlobalScope) {
	globalScopeRegistry.RegisterGlobal(scope)
}

// RegisterScope registers a named scope globally
func RegisterScope(name string, scope Scope) {
	globalScopeRegistry.Register(name, scope)
}

// GetGlobalScopeRegistry returns the global scope registry
func GetGlobalScopeRegistry() *ScopeRegistry {
	return globalScopeRegistry
}
