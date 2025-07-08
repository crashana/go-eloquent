package eloquent

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"
	"unicode"
)

// Model represents the base model interface
type Model interface {
	GetTable() string
	GetPrimaryKey() string
	GetConnection() string
	GetFillable() []string
	GetGuarded() []string
	GetHidden() []string
	GetVisible() []string
	GetCasts() map[string]string
	GetDates() []string
	GetTimestamps() bool
	GetCreatedAtColumn() string
	GetUpdatedAtColumn() string
	GetDeletedAtColumn() string

	// Query methods
	Save() error
	Delete() error
	ForceDelete() error
	Restore() error
	Fill(attributes map[string]interface{}) Model
	Update(attributes map[string]interface{}) error
	Fresh() (Model, error)
	Refresh() error

	// Attribute methods
	GetAttribute(key string) interface{}
	SetAttribute(key string, value interface{})
	GetOriginal(key string) interface{}
	GetDirty() map[string]interface{}
	IsDirty(key ...string) bool
	IsClean(key ...string) bool

	// Serialization
	ToMap() map[string]interface{}
	ToJSON() ([]byte, error)
}

// BaseModel provides the default implementation
type BaseModel struct {
	// Configuration
	table      string
	primaryKey string
	connection string
	fillable   []string
	guarded    []string
	hidden     []string
	visible    []string
	casts      map[string]string
	dates      []string
	timestamps bool
	createdAt  string
	updatedAt  string
	deletedAt  string

	// State
	attributes         map[string]interface{}
	original           map[string]interface{}
	exists             bool
	wasRecentlyCreated bool

	// Relationships
	relations map[string]interface{}
}

// ModelQueryBuilder wraps QueryBuilder and returns model instances
type ModelQueryBuilder struct {
	*QueryBuilder
	model Model
}

// TypedModelQueryBuilder wraps QueryBuilder and returns typed model instances
type TypedModelQueryBuilder[T Model] struct {
	*QueryBuilder
	model        Model
	modelFactory func() T
}

// NewModelQueryBuilder creates a new model query builder
func NewModelQueryBuilder(model Model) *ModelQueryBuilder {
	db := DB()
	if db == nil {
		panic("Database connection not initialized")
	}

	qb := NewQueryBuilder(db)
	qb.Table(model.GetTable())

	return &ModelQueryBuilder{
		QueryBuilder: qb,
		model:        model,
	}
}

// Get returns multiple model instances
func (mqb *ModelQueryBuilder) Get() ([]Model, error) {
	results, err := mqb.QueryBuilder.Get()
	if err != nil {
		return nil, err
	}

	var models []Model
	for _, result := range results {
		model := mqb.newModelInstance()
		mqb.fillModelFromMap(model, result)
		models = append(models, model)
	}

	return models, nil
}

// First returns the first model instance
func (mqb *ModelQueryBuilder) First() (Model, error) {
	result, err := mqb.QueryBuilder.First()
	if err != nil {
		return nil, err
	}

	model := mqb.newModelInstance()
	mqb.fillModelFromMap(model, result)
	return model, nil
}

// FirstOrFail returns the first model instance or fails
func (mqb *ModelQueryBuilder) FirstOrFail() (Model, error) {
	result, err := mqb.QueryBuilder.FirstOrFail()
	if err != nil {
		return nil, err
	}

	model := mqb.newModelInstance()
	mqb.fillModelFromMap(model, result)
	return model, nil
}

// Find finds a model by primary key
func (mqb *ModelQueryBuilder) Find(id interface{}) (Model, error) {
	result, err := mqb.QueryBuilder.Find(id)
	if err != nil {
		return nil, err
	}

	model := mqb.newModelInstance()
	mqb.fillModelFromMap(model, result)
	return model, nil
}

// FindOrFail finds a model by primary key or fails
func (mqb *ModelQueryBuilder) FindOrFail(id interface{}) (Model, error) {
	result, err := mqb.QueryBuilder.FindOrFail(id)
	if err != nil {
		return nil, err
	}

	model := mqb.newModelInstance()
	mqb.fillModelFromMap(model, result)
	return model, nil
}

// Where adds a where clause and returns ModelQueryBuilder
func (mqb *ModelQueryBuilder) Where(column string, args ...interface{}) *ModelQueryBuilder {
	mqb.QueryBuilder.Where(column, args...)
	return mqb
}

// OrWhere adds an OR where clause and returns ModelQueryBuilder
func (mqb *ModelQueryBuilder) OrWhere(column string, args ...interface{}) *ModelQueryBuilder {
	mqb.QueryBuilder.OrWhere(column, args...)
	return mqb
}

// WhereIn adds a where in clause and returns ModelQueryBuilder
func (mqb *ModelQueryBuilder) WhereIn(column string, values []interface{}) *ModelQueryBuilder {
	mqb.QueryBuilder.WhereIn(column, values)
	return mqb
}

// WhereNotIn adds a where not in clause and returns ModelQueryBuilder
func (mqb *ModelQueryBuilder) WhereNotIn(column string, values []interface{}) *ModelQueryBuilder {
	mqb.QueryBuilder.WhereNotIn(column, values)
	return mqb
}

// WhereNull adds a where null clause and returns ModelQueryBuilder
func (mqb *ModelQueryBuilder) WhereNull(column string) *ModelQueryBuilder {
	mqb.QueryBuilder.WhereNull(column)
	return mqb
}

// WhereNotNull adds a where not null clause and returns ModelQueryBuilder
func (mqb *ModelQueryBuilder) WhereNotNull(column string) *ModelQueryBuilder {
	mqb.QueryBuilder.WhereNotNull(column)
	return mqb
}

// OrderBy adds an order by clause and returns ModelQueryBuilder
func (mqb *ModelQueryBuilder) OrderBy(column, direction string) *ModelQueryBuilder {
	mqb.QueryBuilder.OrderBy(column, direction)
	return mqb
}

// OrderByDesc adds an order by desc clause and returns ModelQueryBuilder
func (mqb *ModelQueryBuilder) OrderByDesc(column string) *ModelQueryBuilder {
	mqb.QueryBuilder.OrderByDesc(column)
	return mqb
}

// Limit adds a limit clause and returns ModelQueryBuilder
func (mqb *ModelQueryBuilder) Limit(limit int) *ModelQueryBuilder {
	mqb.QueryBuilder.Limit(limit)
	return mqb
}

// Take adds a limit clause and returns ModelQueryBuilder
func (mqb *ModelQueryBuilder) Take(limit int) *ModelQueryBuilder {
	mqb.QueryBuilder.Take(limit)
	return mqb
}

// Offset adds an offset clause and returns ModelQueryBuilder
func (mqb *ModelQueryBuilder) Offset(offset int) *ModelQueryBuilder {
	mqb.QueryBuilder.Offset(offset)
	return mqb
}

// Skip adds an offset clause and returns ModelQueryBuilder
func (mqb *ModelQueryBuilder) Skip(offset int) *ModelQueryBuilder {
	mqb.QueryBuilder.Skip(offset)
	return mqb
}

// newModelInstance creates a new instance of the model
func (mqb *ModelQueryBuilder) newModelInstance() Model {
	modelType := reflect.TypeOf(mqb.model).Elem()
	newModel := reflect.New(modelType).Interface().(Model)

	// Initialize embedded BaseModel if it exists
	modelValue := reflect.ValueOf(newModel)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	// Look for embedded BaseModel field and initialize it
	for i := 0; i < modelValue.NumField(); i++ {
		field := modelValue.Field(i)
		if field.Type() == reflect.TypeOf((*BaseModel)(nil)) && field.CanSet() {
			field.Set(reflect.ValueOf(NewBaseModel()))
			break
		}
	}

	return newModel
}

// fillModelFromMap fills a model with data from a map
func (mqb *ModelQueryBuilder) fillModelFromMap(model Model, data map[string]interface{}) {
	// Use reflection to find the embedded BaseModel
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	// Look for embedded BaseModel field
	var baseModel *BaseModel
	for i := 0; i < modelValue.NumField(); i++ {
		field := modelValue.Field(i)
		if field.Type() == reflect.TypeOf((*BaseModel)(nil)) {
			baseModel = field.Interface().(*BaseModel)
			break
		}
	}

	if baseModel != nil {
		if baseModel.attributes == nil {
			baseModel.attributes = make(map[string]interface{})
		}
		if baseModel.original == nil {
			baseModel.original = make(map[string]interface{})
		}

		for key, value := range data {
			baseModel.attributes[key] = value
			baseModel.original[key] = value
		}

		baseModel.exists = true
		baseModel.wasRecentlyCreated = false

		// Copy table configuration from the template model
		if mqb.model != nil {
			baseModel.table = mqb.model.GetTable()
			baseModel.primaryKey = mqb.model.GetPrimaryKey()
			baseModel.fillable = mqb.model.GetFillable()
			baseModel.guarded = mqb.model.GetGuarded()
			baseModel.hidden = mqb.model.GetHidden()
			baseModel.visible = mqb.model.GetVisible()
			baseModel.casts = mqb.model.GetCasts()
			baseModel.dates = mqb.model.GetDates()
			baseModel.timestamps = mqb.model.GetTimestamps()
			baseModel.createdAt = mqb.model.GetCreatedAtColumn()
			baseModel.updatedAt = mqb.model.GetUpdatedAtColumn()
			baseModel.deletedAt = mqb.model.GetDeletedAtColumn()
		}
	}

	// Auto-sync attributes to struct fields
	mqb.autoSyncAttributes(model, data)
}

// autoSyncAttributes automatically syncs database attributes to struct fields
func (mqb *ModelQueryBuilder) autoSyncAttributes(model Model, data map[string]interface{}) {
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	modelType := modelValue.Type()

	// Iterate through all struct fields
	for i := 0; i < modelValue.NumField(); i++ {
		field := modelValue.Field(i)
		fieldType := modelType.Field(i)

		// Skip unexported fields and BaseModel
		if !field.CanSet() || fieldType.Type == reflect.TypeOf((*BaseModel)(nil)) {
			continue
		}

		// Get the database column name from the db tag, or use field name
		dbTag := fieldType.Tag.Get("db")
		if dbTag == "" {
			dbTag = toSnakeCase(fieldType.Name)
		}

		// Check if we have data for this field
		if value, exists := data[dbTag]; exists && value != nil {
			mqb.setFieldValue(field, value)
		}
	}
}

// setFieldValue sets a struct field value with proper type conversion
func (mqb *ModelQueryBuilder) setFieldValue(field reflect.Value, value interface{}) {
	if !field.CanSet() {
		return
	}

	valueType := reflect.TypeOf(value)
	fieldType := field.Type()

	// Handle different type conversions
	switch fieldType.Kind() {
	case reflect.String:
		if str, ok := value.(string); ok {
			field.SetString(str)
		}
	case reflect.Bool:
		if b, ok := value.(bool); ok {
			field.SetBool(b)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if i, ok := value.(int64); ok {
			field.SetInt(i)
		} else if i, ok := value.(int); ok {
			field.SetInt(int64(i))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if u, ok := value.(uint64); ok {
			field.SetUint(u)
		} else if u, ok := value.(uint); ok {
			field.SetUint(uint64(u))
		}
	case reflect.Float32, reflect.Float64:
		if f, ok := value.(float64); ok {
			field.SetFloat(f)
		} else if f, ok := value.(float32); ok {
			field.SetFloat(float64(f))
		}
	default:
		// Handle time.Time and other types
		if fieldType == reflect.TypeOf(time.Time{}) {
			if t, ok := value.(time.Time); ok {
				field.Set(reflect.ValueOf(t))
			}
		} else if valueType.AssignableTo(fieldType) {
			field.Set(reflect.ValueOf(value))
		}
	}
}

// Static query methods for BaseModel
func (m *BaseModel) Query() *ModelQueryBuilder {
	return NewModelQueryBuilder(m)
}

// Where creates a new query with a where clause
func (m *BaseModel) Where(column string, args ...interface{}) *ModelQueryBuilder {
	return m.Query().Where(column, args...)
}

// OrWhere creates a new query with an OR where clause
func (m *BaseModel) OrWhere(column string, args ...interface{}) *ModelQueryBuilder {
	return m.Query().OrWhere(column, args...)
}

// WhereIn creates a new query with a where in clause
func (m *BaseModel) WhereIn(column string, values []interface{}) *ModelQueryBuilder {
	return m.Query().WhereIn(column, values)
}

// All returns all records
func (m *BaseModel) All() ([]Model, error) {
	return m.Query().Get()
}

// First returns the first record
func (m *BaseModel) First() (Model, error) {
	return m.Query().First()
}

// Find finds a record by primary key
func (m *BaseModel) Find(id interface{}) (Model, error) {
	return m.Query().Find(id)
}

// NewBaseModel creates a new base model instance
func NewBaseModel() *BaseModel {
	return &BaseModel{
		primaryKey: "id",
		timestamps: true,
		createdAt:  "created_at",
		updatedAt:  "updated_at",
		deletedAt:  "deleted_at",
		attributes: make(map[string]interface{}),
		original:   make(map[string]interface{}),
		relations:  make(map[string]interface{}),
		casts:      make(map[string]string),
	}
}

// Table configuration methods
func (m *BaseModel) Table(table string) *BaseModel {
	m.table = table
	return m
}

func (m *BaseModel) PrimaryKey(key string) *BaseModel {
	m.primaryKey = key
	return m
}

func (m *BaseModel) Connection(conn string) *BaseModel {
	m.connection = conn
	return m
}

func (m *BaseModel) Fillable(fields ...string) *BaseModel {
	m.fillable = fields
	return m
}

func (m *BaseModel) Guarded(fields ...string) *BaseModel {
	m.guarded = fields
	return m
}

func (m *BaseModel) Hidden(fields ...string) *BaseModel {
	m.hidden = fields
	return m
}

func (m *BaseModel) Visible(fields ...string) *BaseModel {
	m.visible = fields
	return m
}

func (m *BaseModel) Casts(casts map[string]string) *BaseModel {
	m.casts = casts
	return m
}

func (m *BaseModel) Dates(dates ...string) *BaseModel {
	m.dates = dates
	return m
}

func (m *BaseModel) WithoutTimestamps() *BaseModel {
	m.timestamps = false
	return m
}

// Getter methods
func (m *BaseModel) GetTable() string {
	if m.table != "" {
		return m.table
	}
	// Auto-generate table name from struct name
	modelType := reflect.TypeOf(m)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	return toSnakeCase(modelType.Name()) + "s"
}

func (m *BaseModel) GetPrimaryKey() string {
	return m.primaryKey
}

func (m *BaseModel) GetConnection() string {
	return m.connection
}

func (m *BaseModel) GetFillable() []string {
	return m.fillable
}

func (m *BaseModel) GetGuarded() []string {
	return m.guarded
}

func (m *BaseModel) GetHidden() []string {
	return m.hidden
}

func (m *BaseModel) GetVisible() []string {
	return m.visible
}

func (m *BaseModel) GetCasts() map[string]string {
	return m.casts
}

func (m *BaseModel) GetDates() []string {
	return m.dates
}

func (m *BaseModel) GetTimestamps() bool {
	return m.timestamps
}

func (m *BaseModel) GetCreatedAtColumn() string {
	return m.createdAt
}

func (m *BaseModel) GetUpdatedAtColumn() string {
	return m.updatedAt
}

func (m *BaseModel) GetDeletedAtColumn() string {
	return m.deletedAt
}

// Attribute methods
func (m *BaseModel) GetAttribute(key string) interface{} {
	value, exists := m.attributes[key]
	if !exists {
		return nil
	}

	// Apply casts
	if castType, hasCast := m.casts[key]; hasCast {
		return m.castAttribute(key, value, castType)
	}

	return value
}

func (m *BaseModel) SetAttribute(key string, value interface{}) {
	m.attributes[key] = value
}

func (m *BaseModel) GetOriginal(key string) interface{} {
	return m.original[key]
}

func (m *BaseModel) GetDirty() map[string]interface{} {
	dirty := make(map[string]interface{})

	for key, value := range m.attributes {
		original, hasOriginal := m.original[key]
		if !hasOriginal || !m.valuesEqual(value, original) {
			dirty[key] = value
		}
	}

	return dirty
}

func (m *BaseModel) IsDirty(keys ...string) bool {
	dirty := m.GetDirty()

	if len(keys) == 0 {
		return len(dirty) > 0
	}

	for _, key := range keys {
		if _, isDirty := dirty[key]; isDirty {
			return true
		}
	}

	return false
}

func (m *BaseModel) IsClean(keys ...string) bool {
	return !m.IsDirty(keys...)
}

// Fill method
func (m *BaseModel) Fill(attributes map[string]interface{}) Model {
	for key, value := range attributes {
		if m.isFillable(key) {
			m.SetAttribute(key, value)
		}
	}
	return m
}

// Save method
func (m *BaseModel) Save() error {
	if m.exists {
		return m.performUpdate()
	}
	return m.performInsert()
}

// Delete methods
func (m *BaseModel) Delete() error {
	if m.usesSoftDeletes() {
		return m.runSoftDelete()
	}
	return m.performDelete()
}

func (m *BaseModel) ForceDelete() error {
	return m.performDelete()
}

func (m *BaseModel) Restore() error {
	if !m.usesSoftDeletes() {
		return fmt.Errorf("model does not use soft deletes")
	}
	return m.performRestore()
}

// Update method
func (m *BaseModel) Update(attributes map[string]interface{}) error {
	m.Fill(attributes)
	return m.performUpdate()
}

func (m *BaseModel) Fresh() (Model, error) {
	// Implementation would query fresh data from database
	return nil, fmt.Errorf("not implemented")
}

func (m *BaseModel) Refresh() error {
	// Implementation would refresh current model from database
	return fmt.Errorf("not implemented")
}

// Serialization methods
func (m *BaseModel) ToMap() map[string]interface{} {
	result := make(map[string]interface{})

	for key := range m.attributes {
		if !m.isHidden(key) {
			result[key] = m.GetAttribute(key)
		}
	}

	// Add relations
	for key, relation := range m.relations {
		if !m.isHidden(key) {
			result[key] = relation
		}
	}

	return result
}

func (m *BaseModel) ToJSON() ([]byte, error) {
	// Implementation would marshal to JSON
	return nil, fmt.Errorf("not implemented")
}

// Helper methods
func (m *BaseModel) isFillable(key string) bool {
	if len(m.fillable) > 0 {
		return m.contains(m.fillable, key)
	}

	if len(m.guarded) > 0 {
		return !m.contains(m.guarded, key)
	}

	return true
}

func (m *BaseModel) isHidden(key string) bool {
	if len(m.visible) > 0 {
		return !m.contains(m.visible, key)
	}

	return m.contains(m.hidden, key)
}

func (m *BaseModel) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (m *BaseModel) valuesEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func (m *BaseModel) usesSoftDeletes() bool {
	return m.deletedAt != ""
}

func (m *BaseModel) castAttribute(_ string, val interface{}, castType string) interface{} {
	switch castType {
	case "string":
		return fmt.Sprintf("%v", val)
	case "int":
		if v, ok := val.(int); ok {
			return v
		}
		return 0
	case "float":
		if v, ok := val.(float64); ok {
			return v
		}
		return 0.0
	case "bool":
		if v, ok := val.(bool); ok {
			return v
		}
		return false
	case "datetime":
		if v, ok := val.(time.Time); ok {
			return v
		}
		return time.Time{}
	}
	return val
}

// Database operation methods (to be implemented with actual DB connection)
func (m *BaseModel) performInsert() error {
	db := DB()
	if db == nil {
		return fmt.Errorf("database connection not initialized")
	}

	if m.timestamps {
		now := time.Now()
		m.SetAttribute(m.createdAt, now)
		m.SetAttribute(m.updatedAt, now)
	}

	// Generate ID for primary key if needed
	if m.GetAttribute(m.primaryKey) == nil {
		// For PostgreSQL, let the database generate the UUID
		db := DB()
		if db != nil && db.Driver == "postgres" {
			// Use PostgreSQL's gen_random_uuid() function
			var id string
			err := db.DB.QueryRow("SELECT gen_random_uuid()").Scan(&id)
			if err != nil {
				// Fallback to manual UUID generation
				m.SetAttribute(m.primaryKey, generateID())
			} else {
				m.SetAttribute(m.primaryKey, id)
			}
		} else {
			m.SetAttribute(m.primaryKey, generateID())
		}
	}

	// Build INSERT query
	var columns []string
	var values []interface{}
	var placeholders []string

	for key, value := range m.attributes {
		columns = append(columns, key)
		values = append(values, value)
		placeholders = append(placeholders, "?")
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		m.GetTable(),
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	// Convert ? to $1, $2, etc. for PostgreSQL
	if db.Driver == "postgres" {
		for i := 0; i < len(placeholders); i++ {
			query = strings.Replace(query, "?", fmt.Sprintf("$%d", i+1), 1)
		}
	}

	_, err := db.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}

	m.exists = true
	m.wasRecentlyCreated = true
	m.syncOriginal()
	return nil
}

func (m *BaseModel) performUpdate() error {
	db := DB()
	if db == nil {
		return fmt.Errorf("database connection not initialized")
	}

	if m.timestamps {
		m.SetAttribute(m.updatedAt, time.Now())
	}

	// Build UPDATE query
	var setParts []string
	var values []interface{}

	for key, value := range m.attributes {
		if key != m.primaryKey { // Don't update primary key
			setParts = append(setParts, fmt.Sprintf("%s = ?", key))
			values = append(values, value)
		}
	}

	// Add primary key value for WHERE clause
	primaryKeyValue := m.GetAttribute(m.primaryKey)
	values = append(values, primaryKeyValue)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?",
		m.GetTable(),
		strings.Join(setParts, ", "),
		m.primaryKey)

	// Convert ? to $1, $2, etc. for PostgreSQL
	if db.Driver == "postgres" {
		placeholderIndex := 1
		for strings.Contains(query, "?") {
			query = strings.Replace(query, "?", fmt.Sprintf("$%d", placeholderIndex), 1)
			placeholderIndex++
		}
	}

	_, err := db.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}

	m.syncOriginal()
	return nil
}

func (m *BaseModel) performDelete() error {
	db := DB()
	if db == nil {
		return fmt.Errorf("database connection not initialized")
	}

	primaryKeyValue := m.GetAttribute(m.primaryKey)
	if primaryKeyValue == nil {
		return fmt.Errorf("cannot delete record without primary key")
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", m.GetTable(), m.primaryKey)

	// Convert ? to $1 for PostgreSQL
	if db.Driver == "postgres" {
		query = strings.Replace(query, "?", "$1", 1)
	}

	_, err := db.Exec(query, primaryKeyValue)
	if err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}

	return nil
}

func (m *BaseModel) runSoftDelete() error {
	// Implementation would set deleted_at timestamp
	m.SetAttribute(m.deletedAt, time.Now())
	return m.performUpdate()
}

func (m *BaseModel) performRestore() error {
	// Implementation would set deleted_at to null
	m.SetAttribute(m.deletedAt, nil)
	return m.performUpdate()
}

func (m *BaseModel) syncOriginal() {
	m.original = make(map[string]interface{})
	for k, v := range m.attributes {
		m.original[k] = v
	}
}

// Helper utility functions

// toSnakeCase converts CamelCase to snake_case
func toSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

// generateID generates a UUID-like ID for PostgreSQL compatibility
func generateID() string {
	// Generate a UUID-like string
	b := make([]byte, 16)
	rand.Read(b)

	// Format as UUID: xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx
	return fmt.Sprintf("%x-%x-4%x-%x%x-%x",
		b[0:4],
		b[4:6],
		b[6:8],
		b[8:9],
		b[9:10],
		b[10:16])
}

// Static-like methods that work like Eloquent
// These create a new instance and return the query builder

// Where creates a new query with where clause (static-like)
func Where(model Model, column string, args ...interface{}) *ModelQueryBuilder {
	return NewModelQueryBuilder(model).Where(column, args...)
}

// First gets the first record (static-like)
func First(model Model) (Model, error) {
	return NewModelQueryBuilder(model).First()
}

// All gets all records (static-like)
func All(model Model) ([]Model, error) {
	return NewModelQueryBuilder(model).Get()
}

// Find finds by primary key (static-like)
func Find(model Model, id interface{}) (Model, error) {
	return NewModelQueryBuilder(model).Find(id)
}

// Create creates a new record (static-like)
func Create(model Model, attributes map[string]interface{}) (Model, error) {
	newModel := model
	if baseModel, ok := newModel.(*BaseModel); ok {
		baseModel.Fill(attributes)
		err := baseModel.Save()
		if err != nil {
			return nil, err
		}
		return newModel, nil
	}
	return nil, fmt.Errorf("model does not support Create")
}

// ModelStatic provides Eloquent-style static methods for any model
type ModelStatic[T Model] struct {
	modelFactory func() T
}

// NewModelStatic creates a new ModelStatic instance for any model type
func NewModelStatic[T Model](factory func() T) *ModelStatic[T] {
	return &ModelStatic[T]{
		modelFactory: factory,
	}
}

// Where creates a new query with where clause (static-like)
func (ms *ModelStatic[T]) Where(column string, args ...interface{}) *TypedModelQueryBuilder[T] {
	model := ms.modelFactory()
	qb := NewModelQueryBuilder(model).Where(column, args...)
	return &TypedModelQueryBuilder[T]{
		QueryBuilder: qb.QueryBuilder,
		model:        model,
		modelFactory: ms.modelFactory,
	}
}

// First gets the first record (static-like) - returns the typed model directly
func (ms *ModelStatic[T]) First() (T, error) {
	model := ms.modelFactory()
	result, err := NewModelQueryBuilder(model).First()
	if err != nil {
		var zero T
		return zero, err
	}
	return result.(T), nil
}

// All gets all records (static-like) - returns slice of typed models
func (ms *ModelStatic[T]) All() ([]T, error) {
	model := ms.modelFactory()
	results, err := NewModelQueryBuilder(model).Get()
	if err != nil {
		return nil, err
	}

	typedResults := make([]T, len(results))
	for i, result := range results {
		typedResults[i] = result.(T)
	}
	return typedResults, nil
}

// Find finds by primary key (static-like) - returns the typed model directly
func (ms *ModelStatic[T]) Find(id interface{}) (T, error) {
	model := ms.modelFactory()
	result, err := NewModelQueryBuilder(model).Find(id)
	if err != nil {
		var zero T
		return zero, err
	}
	return result.(T), nil
}

// Create creates a new record (static-like) - returns the typed model directly
func (ms *ModelStatic[T]) Create(attributes map[string]interface{}) (T, error) {
	model := ms.modelFactory()

	// Use reflection to find the embedded BaseModel
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	// Look for embedded BaseModel field
	var baseModel *BaseModel
	for i := 0; i < modelValue.NumField(); i++ {
		field := modelValue.Field(i)
		if field.Type() == reflect.TypeOf((*BaseModel)(nil)) {
			baseModel = field.Interface().(*BaseModel)
			break
		}
	}

	if baseModel != nil {
		baseModel.Fill(attributes)
		err := baseModel.Save()
		if err != nil {
			var zero T
			return zero, err
		}

		// Sync attributes back to struct fields after creation
		mqb := &ModelQueryBuilder{
			QueryBuilder: NewQueryBuilder(DB()),
			model:        model,
		}
		mqb.autoSyncAttributes(model, baseModel.attributes)

		return model, nil
	}

	var zero T
	return zero, fmt.Errorf("model does not support Create")
}

// Get gets all records (alias for All) - returns slice of typed models
func (ms *ModelStatic[T]) Get() ([]T, error) {
	return ms.All()
}

// Methods for TypedModelQueryBuilder

// First returns the first typed model instance
func (tmqb *TypedModelQueryBuilder[T]) First() (T, error) {
	result, err := tmqb.QueryBuilder.First()
	if err != nil {
		var zero T
		return zero, err
	}

	model := tmqb.modelFactory()
	mqb := &ModelQueryBuilder{
		QueryBuilder: tmqb.QueryBuilder,
		model:        model,
	}
	mqb.fillModelFromMap(model, result)
	return model, nil
}

// Get returns multiple typed model instances
func (tmqb *TypedModelQueryBuilder[T]) Get() ([]T, error) {
	results, err := tmqb.QueryBuilder.Get()
	if err != nil {
		return nil, err
	}

	var models []T
	for _, result := range results {
		model := tmqb.modelFactory()
		mqb := &ModelQueryBuilder{
			QueryBuilder: tmqb.QueryBuilder,
			model:        model,
		}
		mqb.fillModelFromMap(model, result)
		models = append(models, model)
	}

	return models, nil
}

// Find finds a typed model by primary key
func (tmqb *TypedModelQueryBuilder[T]) Find(id interface{}) (T, error) {
	result, err := tmqb.QueryBuilder.Find(id)
	if err != nil {
		var zero T
		return zero, err
	}

	model := tmqb.modelFactory()
	mqb := &ModelQueryBuilder{
		QueryBuilder: tmqb.QueryBuilder,
		model:        model,
	}
	mqb.fillModelFromMap(model, result)
	return model, nil
}

// Where adds a where clause and returns TypedModelQueryBuilder
func (tmqb *TypedModelQueryBuilder[T]) Where(column string, args ...interface{}) *TypedModelQueryBuilder[T] {
	tmqb.QueryBuilder.Where(column, args...)
	return tmqb
}

// OrWhere adds an OR where clause and returns TypedModelQueryBuilder
func (tmqb *TypedModelQueryBuilder[T]) OrWhere(column string, args ...interface{}) *TypedModelQueryBuilder[T] {
	tmqb.QueryBuilder.OrWhere(column, args...)
	return tmqb
}

// WhereIn adds a where in clause and returns TypedModelQueryBuilder
func (tmqb *TypedModelQueryBuilder[T]) WhereIn(column string, values []interface{}) *TypedModelQueryBuilder[T] {
	tmqb.QueryBuilder.WhereIn(column, values)
	return tmqb
}

// WhereNotIn adds a where not in clause and returns TypedModelQueryBuilder
func (tmqb *TypedModelQueryBuilder[T]) WhereNotIn(column string, values []interface{}) *TypedModelQueryBuilder[T] {
	tmqb.QueryBuilder.WhereNotIn(column, values)
	return tmqb
}

// WhereNull adds a where null clause and returns TypedModelQueryBuilder
func (tmqb *TypedModelQueryBuilder[T]) WhereNull(column string) *TypedModelQueryBuilder[T] {
	tmqb.QueryBuilder.WhereNull(column)
	return tmqb
}

// WhereNotNull adds a where not null clause and returns TypedModelQueryBuilder
func (tmqb *TypedModelQueryBuilder[T]) WhereNotNull(column string) *TypedModelQueryBuilder[T] {
	tmqb.QueryBuilder.WhereNotNull(column)
	return tmqb
}

// OrderBy adds an order by clause and returns TypedModelQueryBuilder
func (tmqb *TypedModelQueryBuilder[T]) OrderBy(column, direction string) *TypedModelQueryBuilder[T] {
	tmqb.QueryBuilder.OrderBy(column, direction)
	return tmqb
}

// OrderByDesc adds an order by desc clause and returns TypedModelQueryBuilder
func (tmqb *TypedModelQueryBuilder[T]) OrderByDesc(column string) *TypedModelQueryBuilder[T] {
	tmqb.QueryBuilder.OrderByDesc(column)
	return tmqb
}

// Limit adds a limit clause and returns TypedModelQueryBuilder
func (tmqb *TypedModelQueryBuilder[T]) Limit(limit int) *TypedModelQueryBuilder[T] {
	tmqb.QueryBuilder.Limit(limit)
	return tmqb
}

// Take adds a limit clause and returns TypedModelQueryBuilder
func (tmqb *TypedModelQueryBuilder[T]) Take(limit int) *TypedModelQueryBuilder[T] {
	tmqb.QueryBuilder.Take(limit)
	return tmqb
}

// Offset adds an offset clause and returns TypedModelQueryBuilder
func (tmqb *TypedModelQueryBuilder[T]) Offset(offset int) *TypedModelQueryBuilder[T] {
	tmqb.QueryBuilder.Offset(offset)
	return tmqb
}

// Skip adds an offset clause and returns TypedModelQueryBuilder
func (tmqb *TypedModelQueryBuilder[T]) Skip(offset int) *TypedModelQueryBuilder[T] {
	tmqb.QueryBuilder.Skip(offset)
	return tmqb
}
