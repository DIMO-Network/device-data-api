// Code generated by SQLBoiler 4.13.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// VehicleSignalsTrackingProperty is an object representing the database table.
type VehicleSignalsTrackingProperty struct {
	ID          string      `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name        string      `boil:"name" json:"name" toml:"name" yaml:"name"`
	Description null.String `boil:"description" json:"description,omitempty" toml:"description" yaml:"description,omitempty"`
	CreatedAt   time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt   time.Time   `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *vehicleSignalsTrackingPropertyR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L vehicleSignalsTrackingPropertyL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var VehicleSignalsTrackingPropertyColumns = struct {
	ID          string
	Name        string
	Description string
	CreatedAt   string
	UpdatedAt   string
}{
	ID:          "id",
	Name:        "name",
	Description: "description",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

var VehicleSignalsTrackingPropertyTableColumns = struct {
	ID          string
	Name        string
	Description string
	CreatedAt   string
	UpdatedAt   string
}{
	ID:          "vehicle_signals_tracking_properties.id",
	Name:        "vehicle_signals_tracking_properties.name",
	Description: "vehicle_signals_tracking_properties.description",
	CreatedAt:   "vehicle_signals_tracking_properties.created_at",
	UpdatedAt:   "vehicle_signals_tracking_properties.updated_at",
}

// Generated where

var VehicleSignalsTrackingPropertyWhere = struct {
	ID          whereHelperstring
	Name        whereHelperstring
	Description whereHelpernull_String
	CreatedAt   whereHelpertime_Time
	UpdatedAt   whereHelpertime_Time
}{
	ID:          whereHelperstring{field: "\"device_data_api\".\"vehicle_signals_tracking_properties\".\"id\""},
	Name:        whereHelperstring{field: "\"device_data_api\".\"vehicle_signals_tracking_properties\".\"name\""},
	Description: whereHelpernull_String{field: "\"device_data_api\".\"vehicle_signals_tracking_properties\".\"description\""},
	CreatedAt:   whereHelpertime_Time{field: "\"device_data_api\".\"vehicle_signals_tracking_properties\".\"created_at\""},
	UpdatedAt:   whereHelpertime_Time{field: "\"device_data_api\".\"vehicle_signals_tracking_properties\".\"updated_at\""},
}

// VehicleSignalsTrackingPropertyRels is where relationship names are stored.
var VehicleSignalsTrackingPropertyRels = struct {
}{}

// vehicleSignalsTrackingPropertyR is where relationships are stored.
type vehicleSignalsTrackingPropertyR struct {
}

// NewStruct creates a new relationship struct
func (*vehicleSignalsTrackingPropertyR) NewStruct() *vehicleSignalsTrackingPropertyR {
	return &vehicleSignalsTrackingPropertyR{}
}

// vehicleSignalsTrackingPropertyL is where Load methods for each relationship are stored.
type vehicleSignalsTrackingPropertyL struct{}

var (
	vehicleSignalsTrackingPropertyAllColumns            = []string{"id", "name", "description", "created_at", "updated_at"}
	vehicleSignalsTrackingPropertyColumnsWithoutDefault = []string{"id", "name"}
	vehicleSignalsTrackingPropertyColumnsWithDefault    = []string{"description", "created_at", "updated_at"}
	vehicleSignalsTrackingPropertyPrimaryKeyColumns     = []string{"id"}
	vehicleSignalsTrackingPropertyGeneratedColumns      = []string{}
)

type (
	// VehicleSignalsTrackingPropertySlice is an alias for a slice of pointers to VehicleSignalsTrackingProperty.
	// This should almost always be used instead of []VehicleSignalsTrackingProperty.
	VehicleSignalsTrackingPropertySlice []*VehicleSignalsTrackingProperty
	// VehicleSignalsTrackingPropertyHook is the signature for custom VehicleSignalsTrackingProperty hook methods
	VehicleSignalsTrackingPropertyHook func(context.Context, boil.ContextExecutor, *VehicleSignalsTrackingProperty) error

	vehicleSignalsTrackingPropertyQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	vehicleSignalsTrackingPropertyType                 = reflect.TypeOf(&VehicleSignalsTrackingProperty{})
	vehicleSignalsTrackingPropertyMapping              = queries.MakeStructMapping(vehicleSignalsTrackingPropertyType)
	vehicleSignalsTrackingPropertyPrimaryKeyMapping, _ = queries.BindMapping(vehicleSignalsTrackingPropertyType, vehicleSignalsTrackingPropertyMapping, vehicleSignalsTrackingPropertyPrimaryKeyColumns)
	vehicleSignalsTrackingPropertyInsertCacheMut       sync.RWMutex
	vehicleSignalsTrackingPropertyInsertCache          = make(map[string]insertCache)
	vehicleSignalsTrackingPropertyUpdateCacheMut       sync.RWMutex
	vehicleSignalsTrackingPropertyUpdateCache          = make(map[string]updateCache)
	vehicleSignalsTrackingPropertyUpsertCacheMut       sync.RWMutex
	vehicleSignalsTrackingPropertyUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var vehicleSignalsTrackingPropertyAfterSelectHooks []VehicleSignalsTrackingPropertyHook

var vehicleSignalsTrackingPropertyBeforeInsertHooks []VehicleSignalsTrackingPropertyHook
var vehicleSignalsTrackingPropertyAfterInsertHooks []VehicleSignalsTrackingPropertyHook

var vehicleSignalsTrackingPropertyBeforeUpdateHooks []VehicleSignalsTrackingPropertyHook
var vehicleSignalsTrackingPropertyAfterUpdateHooks []VehicleSignalsTrackingPropertyHook

var vehicleSignalsTrackingPropertyBeforeDeleteHooks []VehicleSignalsTrackingPropertyHook
var vehicleSignalsTrackingPropertyAfterDeleteHooks []VehicleSignalsTrackingPropertyHook

var vehicleSignalsTrackingPropertyBeforeUpsertHooks []VehicleSignalsTrackingPropertyHook
var vehicleSignalsTrackingPropertyAfterUpsertHooks []VehicleSignalsTrackingPropertyHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *VehicleSignalsTrackingProperty) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleSignalsTrackingPropertyAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *VehicleSignalsTrackingProperty) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleSignalsTrackingPropertyBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *VehicleSignalsTrackingProperty) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleSignalsTrackingPropertyAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *VehicleSignalsTrackingProperty) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleSignalsTrackingPropertyBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *VehicleSignalsTrackingProperty) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleSignalsTrackingPropertyAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *VehicleSignalsTrackingProperty) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleSignalsTrackingPropertyBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *VehicleSignalsTrackingProperty) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleSignalsTrackingPropertyAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *VehicleSignalsTrackingProperty) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleSignalsTrackingPropertyBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *VehicleSignalsTrackingProperty) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleSignalsTrackingPropertyAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddVehicleSignalsTrackingPropertyHook registers your hook function for all future operations.
func AddVehicleSignalsTrackingPropertyHook(hookPoint boil.HookPoint, vehicleSignalsTrackingPropertyHook VehicleSignalsTrackingPropertyHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		vehicleSignalsTrackingPropertyAfterSelectHooks = append(vehicleSignalsTrackingPropertyAfterSelectHooks, vehicleSignalsTrackingPropertyHook)
	case boil.BeforeInsertHook:
		vehicleSignalsTrackingPropertyBeforeInsertHooks = append(vehicleSignalsTrackingPropertyBeforeInsertHooks, vehicleSignalsTrackingPropertyHook)
	case boil.AfterInsertHook:
		vehicleSignalsTrackingPropertyAfterInsertHooks = append(vehicleSignalsTrackingPropertyAfterInsertHooks, vehicleSignalsTrackingPropertyHook)
	case boil.BeforeUpdateHook:
		vehicleSignalsTrackingPropertyBeforeUpdateHooks = append(vehicleSignalsTrackingPropertyBeforeUpdateHooks, vehicleSignalsTrackingPropertyHook)
	case boil.AfterUpdateHook:
		vehicleSignalsTrackingPropertyAfterUpdateHooks = append(vehicleSignalsTrackingPropertyAfterUpdateHooks, vehicleSignalsTrackingPropertyHook)
	case boil.BeforeDeleteHook:
		vehicleSignalsTrackingPropertyBeforeDeleteHooks = append(vehicleSignalsTrackingPropertyBeforeDeleteHooks, vehicleSignalsTrackingPropertyHook)
	case boil.AfterDeleteHook:
		vehicleSignalsTrackingPropertyAfterDeleteHooks = append(vehicleSignalsTrackingPropertyAfterDeleteHooks, vehicleSignalsTrackingPropertyHook)
	case boil.BeforeUpsertHook:
		vehicleSignalsTrackingPropertyBeforeUpsertHooks = append(vehicleSignalsTrackingPropertyBeforeUpsertHooks, vehicleSignalsTrackingPropertyHook)
	case boil.AfterUpsertHook:
		vehicleSignalsTrackingPropertyAfterUpsertHooks = append(vehicleSignalsTrackingPropertyAfterUpsertHooks, vehicleSignalsTrackingPropertyHook)
	}
}

// One returns a single vehicleSignalsTrackingProperty record from the query.
func (q vehicleSignalsTrackingPropertyQuery) One(ctx context.Context, exec boil.ContextExecutor) (*VehicleSignalsTrackingProperty, error) {
	o := &VehicleSignalsTrackingProperty{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for vehicle_signals_tracking_properties")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all VehicleSignalsTrackingProperty records from the query.
func (q vehicleSignalsTrackingPropertyQuery) All(ctx context.Context, exec boil.ContextExecutor) (VehicleSignalsTrackingPropertySlice, error) {
	var o []*VehicleSignalsTrackingProperty

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to VehicleSignalsTrackingProperty slice")
	}

	if len(vehicleSignalsTrackingPropertyAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all VehicleSignalsTrackingProperty records in the query.
func (q vehicleSignalsTrackingPropertyQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count vehicle_signals_tracking_properties rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q vehicleSignalsTrackingPropertyQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if vehicle_signals_tracking_properties exists")
	}

	return count > 0, nil
}

// VehicleSignalsTrackingProperties retrieves all the records using an executor.
func VehicleSignalsTrackingProperties(mods ...qm.QueryMod) vehicleSignalsTrackingPropertyQuery {
	mods = append(mods, qm.From("\"device_data_api\".\"vehicle_signals_tracking_properties\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"device_data_api\".\"vehicle_signals_tracking_properties\".*"})
	}

	return vehicleSignalsTrackingPropertyQuery{q}
}

// FindVehicleSignalsTrackingProperty retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindVehicleSignalsTrackingProperty(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*VehicleSignalsTrackingProperty, error) {
	vehicleSignalsTrackingPropertyObj := &VehicleSignalsTrackingProperty{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"device_data_api\".\"vehicle_signals_tracking_properties\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, vehicleSignalsTrackingPropertyObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from vehicle_signals_tracking_properties")
	}

	if err = vehicleSignalsTrackingPropertyObj.doAfterSelectHooks(ctx, exec); err != nil {
		return vehicleSignalsTrackingPropertyObj, err
	}

	return vehicleSignalsTrackingPropertyObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *VehicleSignalsTrackingProperty) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no vehicle_signals_tracking_properties provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		if o.UpdatedAt.IsZero() {
			o.UpdatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(vehicleSignalsTrackingPropertyColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	vehicleSignalsTrackingPropertyInsertCacheMut.RLock()
	cache, cached := vehicleSignalsTrackingPropertyInsertCache[key]
	vehicleSignalsTrackingPropertyInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			vehicleSignalsTrackingPropertyAllColumns,
			vehicleSignalsTrackingPropertyColumnsWithDefault,
			vehicleSignalsTrackingPropertyColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(vehicleSignalsTrackingPropertyType, vehicleSignalsTrackingPropertyMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(vehicleSignalsTrackingPropertyType, vehicleSignalsTrackingPropertyMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"device_data_api\".\"vehicle_signals_tracking_properties\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"device_data_api\".\"vehicle_signals_tracking_properties\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into vehicle_signals_tracking_properties")
	}

	if !cached {
		vehicleSignalsTrackingPropertyInsertCacheMut.Lock()
		vehicleSignalsTrackingPropertyInsertCache[key] = cache
		vehicleSignalsTrackingPropertyInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the VehicleSignalsTrackingProperty.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *VehicleSignalsTrackingProperty) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	vehicleSignalsTrackingPropertyUpdateCacheMut.RLock()
	cache, cached := vehicleSignalsTrackingPropertyUpdateCache[key]
	vehicleSignalsTrackingPropertyUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			vehicleSignalsTrackingPropertyAllColumns,
			vehicleSignalsTrackingPropertyPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update vehicle_signals_tracking_properties, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"device_data_api\".\"vehicle_signals_tracking_properties\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, vehicleSignalsTrackingPropertyPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(vehicleSignalsTrackingPropertyType, vehicleSignalsTrackingPropertyMapping, append(wl, vehicleSignalsTrackingPropertyPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update vehicle_signals_tracking_properties row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for vehicle_signals_tracking_properties")
	}

	if !cached {
		vehicleSignalsTrackingPropertyUpdateCacheMut.Lock()
		vehicleSignalsTrackingPropertyUpdateCache[key] = cache
		vehicleSignalsTrackingPropertyUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q vehicleSignalsTrackingPropertyQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for vehicle_signals_tracking_properties")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for vehicle_signals_tracking_properties")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o VehicleSignalsTrackingPropertySlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), vehicleSignalsTrackingPropertyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"device_data_api\".\"vehicle_signals_tracking_properties\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, vehicleSignalsTrackingPropertyPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in vehicleSignalsTrackingProperty slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all vehicleSignalsTrackingProperty")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *VehicleSignalsTrackingProperty) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no vehicle_signals_tracking_properties provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		o.UpdatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(vehicleSignalsTrackingPropertyColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	vehicleSignalsTrackingPropertyUpsertCacheMut.RLock()
	cache, cached := vehicleSignalsTrackingPropertyUpsertCache[key]
	vehicleSignalsTrackingPropertyUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			vehicleSignalsTrackingPropertyAllColumns,
			vehicleSignalsTrackingPropertyColumnsWithDefault,
			vehicleSignalsTrackingPropertyColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			vehicleSignalsTrackingPropertyAllColumns,
			vehicleSignalsTrackingPropertyPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert vehicle_signals_tracking_properties, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(vehicleSignalsTrackingPropertyPrimaryKeyColumns))
			copy(conflict, vehicleSignalsTrackingPropertyPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"device_data_api\".\"vehicle_signals_tracking_properties\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(vehicleSignalsTrackingPropertyType, vehicleSignalsTrackingPropertyMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(vehicleSignalsTrackingPropertyType, vehicleSignalsTrackingPropertyMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert vehicle_signals_tracking_properties")
	}

	if !cached {
		vehicleSignalsTrackingPropertyUpsertCacheMut.Lock()
		vehicleSignalsTrackingPropertyUpsertCache[key] = cache
		vehicleSignalsTrackingPropertyUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single VehicleSignalsTrackingProperty record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *VehicleSignalsTrackingProperty) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no VehicleSignalsTrackingProperty provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), vehicleSignalsTrackingPropertyPrimaryKeyMapping)
	sql := "DELETE FROM \"device_data_api\".\"vehicle_signals_tracking_properties\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from vehicle_signals_tracking_properties")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for vehicle_signals_tracking_properties")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q vehicleSignalsTrackingPropertyQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no vehicleSignalsTrackingPropertyQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from vehicle_signals_tracking_properties")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for vehicle_signals_tracking_properties")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o VehicleSignalsTrackingPropertySlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(vehicleSignalsTrackingPropertyBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), vehicleSignalsTrackingPropertyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"device_data_api\".\"vehicle_signals_tracking_properties\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, vehicleSignalsTrackingPropertyPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from vehicleSignalsTrackingProperty slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for vehicle_signals_tracking_properties")
	}

	if len(vehicleSignalsTrackingPropertyAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *VehicleSignalsTrackingProperty) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindVehicleSignalsTrackingProperty(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *VehicleSignalsTrackingPropertySlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := VehicleSignalsTrackingPropertySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), vehicleSignalsTrackingPropertyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"device_data_api\".\"vehicle_signals_tracking_properties\".* FROM \"device_data_api\".\"vehicle_signals_tracking_properties\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, vehicleSignalsTrackingPropertyPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in VehicleSignalsTrackingPropertySlice")
	}

	*o = slice

	return nil
}

// VehicleSignalsTrackingPropertyExists checks if the VehicleSignalsTrackingProperty row exists.
func VehicleSignalsTrackingPropertyExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"device_data_api\".\"vehicle_signals_tracking_properties\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if vehicle_signals_tracking_properties exists")
	}

	return exists, nil
}
