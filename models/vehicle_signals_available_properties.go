// Code generated by SQLBoiler 4.16.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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
	"github.com/volatiletech/sqlboiler/v4/types"
	"github.com/volatiletech/strmangle"
)

// VehicleSignalsAvailableProperty is an object representing the database table.
type VehicleSignalsAvailableProperty struct {
	ID             string            `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name           string            `boil:"name" json:"name" toml:"name" yaml:"name"`
	CreatedAt      time.Time         `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt      time.Time         `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	PowerTrainType types.StringArray `boil:"power_train_type" json:"power_train_type,omitempty" toml:"power_train_type" yaml:"power_train_type,omitempty"`
	ValidMinLength null.Int          `boil:"valid_min_length" json:"valid_min_length,omitempty" toml:"valid_min_length" yaml:"valid_min_length,omitempty"`

	R *vehicleSignalsAvailablePropertyR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L vehicleSignalsAvailablePropertyL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var VehicleSignalsAvailablePropertyColumns = struct {
	ID             string
	Name           string
	CreatedAt      string
	UpdatedAt      string
	PowerTrainType string
	ValidMinLength string
}{
	ID:             "id",
	Name:           "name",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	PowerTrainType: "power_train_type",
	ValidMinLength: "valid_min_length",
}

var VehicleSignalsAvailablePropertyTableColumns = struct {
	ID             string
	Name           string
	CreatedAt      string
	UpdatedAt      string
	PowerTrainType string
	ValidMinLength string
}{
	ID:             "vehicle_signals_available_properties.id",
	Name:           "vehicle_signals_available_properties.name",
	CreatedAt:      "vehicle_signals_available_properties.created_at",
	UpdatedAt:      "vehicle_signals_available_properties.updated_at",
	PowerTrainType: "vehicle_signals_available_properties.power_train_type",
	ValidMinLength: "vehicle_signals_available_properties.valid_min_length",
}

// Generated where

type whereHelpertypes_StringArray struct{ field string }

func (w whereHelpertypes_StringArray) EQ(x types.StringArray) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpertypes_StringArray) NEQ(x types.StringArray) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpertypes_StringArray) LT(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertypes_StringArray) LTE(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertypes_StringArray) GT(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertypes_StringArray) GTE(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

func (w whereHelpertypes_StringArray) IsNull() qm.QueryMod { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpertypes_StringArray) IsNotNull() qm.QueryMod {
	return qmhelper.WhereIsNotNull(w.field)
}

type whereHelpernull_Int struct{ field string }

func (w whereHelpernull_Int) EQ(x null.Int) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Int) NEQ(x null.Int) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Int) LT(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Int) LTE(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Int) GT(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Int) GTE(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}
func (w whereHelpernull_Int) IN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelpernull_Int) NIN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

func (w whereHelpernull_Int) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Int) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }

var VehicleSignalsAvailablePropertyWhere = struct {
	ID             whereHelperstring
	Name           whereHelperstring
	CreatedAt      whereHelpertime_Time
	UpdatedAt      whereHelpertime_Time
	PowerTrainType whereHelpertypes_StringArray
	ValidMinLength whereHelpernull_Int
}{
	ID:             whereHelperstring{field: "\"device_data_api\".\"vehicle_signals_available_properties\".\"id\""},
	Name:           whereHelperstring{field: "\"device_data_api\".\"vehicle_signals_available_properties\".\"name\""},
	CreatedAt:      whereHelpertime_Time{field: "\"device_data_api\".\"vehicle_signals_available_properties\".\"created_at\""},
	UpdatedAt:      whereHelpertime_Time{field: "\"device_data_api\".\"vehicle_signals_available_properties\".\"updated_at\""},
	PowerTrainType: whereHelpertypes_StringArray{field: "\"device_data_api\".\"vehicle_signals_available_properties\".\"power_train_type\""},
	ValidMinLength: whereHelpernull_Int{field: "\"device_data_api\".\"vehicle_signals_available_properties\".\"valid_min_length\""},
}

// VehicleSignalsAvailablePropertyRels is where relationship names are stored.
var VehicleSignalsAvailablePropertyRels = struct {
}{}

// vehicleSignalsAvailablePropertyR is where relationships are stored.
type vehicleSignalsAvailablePropertyR struct {
}

// NewStruct creates a new relationship struct
func (*vehicleSignalsAvailablePropertyR) NewStruct() *vehicleSignalsAvailablePropertyR {
	return &vehicleSignalsAvailablePropertyR{}
}

// vehicleSignalsAvailablePropertyL is where Load methods for each relationship are stored.
type vehicleSignalsAvailablePropertyL struct{}

var (
	vehicleSignalsAvailablePropertyAllColumns            = []string{"id", "name", "created_at", "updated_at", "power_train_type", "valid_min_length"}
	vehicleSignalsAvailablePropertyColumnsWithoutDefault = []string{"id", "name"}
	vehicleSignalsAvailablePropertyColumnsWithDefault    = []string{"created_at", "updated_at", "power_train_type", "valid_min_length"}
	vehicleSignalsAvailablePropertyPrimaryKeyColumns     = []string{"id"}
	vehicleSignalsAvailablePropertyGeneratedColumns      = []string{}
)

type (
	// VehicleSignalsAvailablePropertySlice is an alias for a slice of pointers to VehicleSignalsAvailableProperty.
	// This should almost always be used instead of []VehicleSignalsAvailableProperty.
	VehicleSignalsAvailablePropertySlice []*VehicleSignalsAvailableProperty
	// VehicleSignalsAvailablePropertyHook is the signature for custom VehicleSignalsAvailableProperty hook methods
	VehicleSignalsAvailablePropertyHook func(context.Context, boil.ContextExecutor, *VehicleSignalsAvailableProperty) error

	vehicleSignalsAvailablePropertyQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	vehicleSignalsAvailablePropertyType                 = reflect.TypeOf(&VehicleSignalsAvailableProperty{})
	vehicleSignalsAvailablePropertyMapping              = queries.MakeStructMapping(vehicleSignalsAvailablePropertyType)
	vehicleSignalsAvailablePropertyPrimaryKeyMapping, _ = queries.BindMapping(vehicleSignalsAvailablePropertyType, vehicleSignalsAvailablePropertyMapping, vehicleSignalsAvailablePropertyPrimaryKeyColumns)
	vehicleSignalsAvailablePropertyInsertCacheMut       sync.RWMutex
	vehicleSignalsAvailablePropertyInsertCache          = make(map[string]insertCache)
	vehicleSignalsAvailablePropertyUpdateCacheMut       sync.RWMutex
	vehicleSignalsAvailablePropertyUpdateCache          = make(map[string]updateCache)
	vehicleSignalsAvailablePropertyUpsertCacheMut       sync.RWMutex
	vehicleSignalsAvailablePropertyUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var vehicleSignalsAvailablePropertyAfterSelectMu sync.Mutex
var vehicleSignalsAvailablePropertyAfterSelectHooks []VehicleSignalsAvailablePropertyHook

var vehicleSignalsAvailablePropertyBeforeInsertMu sync.Mutex
var vehicleSignalsAvailablePropertyBeforeInsertHooks []VehicleSignalsAvailablePropertyHook
var vehicleSignalsAvailablePropertyAfterInsertMu sync.Mutex
var vehicleSignalsAvailablePropertyAfterInsertHooks []VehicleSignalsAvailablePropertyHook

var vehicleSignalsAvailablePropertyBeforeUpdateMu sync.Mutex
var vehicleSignalsAvailablePropertyBeforeUpdateHooks []VehicleSignalsAvailablePropertyHook
var vehicleSignalsAvailablePropertyAfterUpdateMu sync.Mutex
var vehicleSignalsAvailablePropertyAfterUpdateHooks []VehicleSignalsAvailablePropertyHook

var vehicleSignalsAvailablePropertyBeforeDeleteMu sync.Mutex
var vehicleSignalsAvailablePropertyBeforeDeleteHooks []VehicleSignalsAvailablePropertyHook
var vehicleSignalsAvailablePropertyAfterDeleteMu sync.Mutex
var vehicleSignalsAvailablePropertyAfterDeleteHooks []VehicleSignalsAvailablePropertyHook

var vehicleSignalsAvailablePropertyBeforeUpsertMu sync.Mutex
var vehicleSignalsAvailablePropertyBeforeUpsertHooks []VehicleSignalsAvailablePropertyHook
var vehicleSignalsAvailablePropertyAfterUpsertMu sync.Mutex
var vehicleSignalsAvailablePropertyAfterUpsertHooks []VehicleSignalsAvailablePropertyHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *VehicleSignalsAvailableProperty) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleSignalsAvailablePropertyAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *VehicleSignalsAvailableProperty) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleSignalsAvailablePropertyBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *VehicleSignalsAvailableProperty) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleSignalsAvailablePropertyAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *VehicleSignalsAvailableProperty) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleSignalsAvailablePropertyBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *VehicleSignalsAvailableProperty) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleSignalsAvailablePropertyAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *VehicleSignalsAvailableProperty) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleSignalsAvailablePropertyBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *VehicleSignalsAvailableProperty) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleSignalsAvailablePropertyAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *VehicleSignalsAvailableProperty) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleSignalsAvailablePropertyBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *VehicleSignalsAvailableProperty) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleSignalsAvailablePropertyAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddVehicleSignalsAvailablePropertyHook registers your hook function for all future operations.
func AddVehicleSignalsAvailablePropertyHook(hookPoint boil.HookPoint, vehicleSignalsAvailablePropertyHook VehicleSignalsAvailablePropertyHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		vehicleSignalsAvailablePropertyAfterSelectMu.Lock()
		vehicleSignalsAvailablePropertyAfterSelectHooks = append(vehicleSignalsAvailablePropertyAfterSelectHooks, vehicleSignalsAvailablePropertyHook)
		vehicleSignalsAvailablePropertyAfterSelectMu.Unlock()
	case boil.BeforeInsertHook:
		vehicleSignalsAvailablePropertyBeforeInsertMu.Lock()
		vehicleSignalsAvailablePropertyBeforeInsertHooks = append(vehicleSignalsAvailablePropertyBeforeInsertHooks, vehicleSignalsAvailablePropertyHook)
		vehicleSignalsAvailablePropertyBeforeInsertMu.Unlock()
	case boil.AfterInsertHook:
		vehicleSignalsAvailablePropertyAfterInsertMu.Lock()
		vehicleSignalsAvailablePropertyAfterInsertHooks = append(vehicleSignalsAvailablePropertyAfterInsertHooks, vehicleSignalsAvailablePropertyHook)
		vehicleSignalsAvailablePropertyAfterInsertMu.Unlock()
	case boil.BeforeUpdateHook:
		vehicleSignalsAvailablePropertyBeforeUpdateMu.Lock()
		vehicleSignalsAvailablePropertyBeforeUpdateHooks = append(vehicleSignalsAvailablePropertyBeforeUpdateHooks, vehicleSignalsAvailablePropertyHook)
		vehicleSignalsAvailablePropertyBeforeUpdateMu.Unlock()
	case boil.AfterUpdateHook:
		vehicleSignalsAvailablePropertyAfterUpdateMu.Lock()
		vehicleSignalsAvailablePropertyAfterUpdateHooks = append(vehicleSignalsAvailablePropertyAfterUpdateHooks, vehicleSignalsAvailablePropertyHook)
		vehicleSignalsAvailablePropertyAfterUpdateMu.Unlock()
	case boil.BeforeDeleteHook:
		vehicleSignalsAvailablePropertyBeforeDeleteMu.Lock()
		vehicleSignalsAvailablePropertyBeforeDeleteHooks = append(vehicleSignalsAvailablePropertyBeforeDeleteHooks, vehicleSignalsAvailablePropertyHook)
		vehicleSignalsAvailablePropertyBeforeDeleteMu.Unlock()
	case boil.AfterDeleteHook:
		vehicleSignalsAvailablePropertyAfterDeleteMu.Lock()
		vehicleSignalsAvailablePropertyAfterDeleteHooks = append(vehicleSignalsAvailablePropertyAfterDeleteHooks, vehicleSignalsAvailablePropertyHook)
		vehicleSignalsAvailablePropertyAfterDeleteMu.Unlock()
	case boil.BeforeUpsertHook:
		vehicleSignalsAvailablePropertyBeforeUpsertMu.Lock()
		vehicleSignalsAvailablePropertyBeforeUpsertHooks = append(vehicleSignalsAvailablePropertyBeforeUpsertHooks, vehicleSignalsAvailablePropertyHook)
		vehicleSignalsAvailablePropertyBeforeUpsertMu.Unlock()
	case boil.AfterUpsertHook:
		vehicleSignalsAvailablePropertyAfterUpsertMu.Lock()
		vehicleSignalsAvailablePropertyAfterUpsertHooks = append(vehicleSignalsAvailablePropertyAfterUpsertHooks, vehicleSignalsAvailablePropertyHook)
		vehicleSignalsAvailablePropertyAfterUpsertMu.Unlock()
	}
}

// One returns a single vehicleSignalsAvailableProperty record from the query.
func (q vehicleSignalsAvailablePropertyQuery) One(ctx context.Context, exec boil.ContextExecutor) (*VehicleSignalsAvailableProperty, error) {
	o := &VehicleSignalsAvailableProperty{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for vehicle_signals_available_properties")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all VehicleSignalsAvailableProperty records from the query.
func (q vehicleSignalsAvailablePropertyQuery) All(ctx context.Context, exec boil.ContextExecutor) (VehicleSignalsAvailablePropertySlice, error) {
	var o []*VehicleSignalsAvailableProperty

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to VehicleSignalsAvailableProperty slice")
	}

	if len(vehicleSignalsAvailablePropertyAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all VehicleSignalsAvailableProperty records in the query.
func (q vehicleSignalsAvailablePropertyQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count vehicle_signals_available_properties rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q vehicleSignalsAvailablePropertyQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if vehicle_signals_available_properties exists")
	}

	return count > 0, nil
}

// VehicleSignalsAvailableProperties retrieves all the records using an executor.
func VehicleSignalsAvailableProperties(mods ...qm.QueryMod) vehicleSignalsAvailablePropertyQuery {
	mods = append(mods, qm.From("\"device_data_api\".\"vehicle_signals_available_properties\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"device_data_api\".\"vehicle_signals_available_properties\".*"})
	}

	return vehicleSignalsAvailablePropertyQuery{q}
}

// FindVehicleSignalsAvailableProperty retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindVehicleSignalsAvailableProperty(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*VehicleSignalsAvailableProperty, error) {
	vehicleSignalsAvailablePropertyObj := &VehicleSignalsAvailableProperty{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"device_data_api\".\"vehicle_signals_available_properties\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, vehicleSignalsAvailablePropertyObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from vehicle_signals_available_properties")
	}

	if err = vehicleSignalsAvailablePropertyObj.doAfterSelectHooks(ctx, exec); err != nil {
		return vehicleSignalsAvailablePropertyObj, err
	}

	return vehicleSignalsAvailablePropertyObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *VehicleSignalsAvailableProperty) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no vehicle_signals_available_properties provided for insertion")
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

	nzDefaults := queries.NonZeroDefaultSet(vehicleSignalsAvailablePropertyColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	vehicleSignalsAvailablePropertyInsertCacheMut.RLock()
	cache, cached := vehicleSignalsAvailablePropertyInsertCache[key]
	vehicleSignalsAvailablePropertyInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			vehicleSignalsAvailablePropertyAllColumns,
			vehicleSignalsAvailablePropertyColumnsWithDefault,
			vehicleSignalsAvailablePropertyColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(vehicleSignalsAvailablePropertyType, vehicleSignalsAvailablePropertyMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(vehicleSignalsAvailablePropertyType, vehicleSignalsAvailablePropertyMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"device_data_api\".\"vehicle_signals_available_properties\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"device_data_api\".\"vehicle_signals_available_properties\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into vehicle_signals_available_properties")
	}

	if !cached {
		vehicleSignalsAvailablePropertyInsertCacheMut.Lock()
		vehicleSignalsAvailablePropertyInsertCache[key] = cache
		vehicleSignalsAvailablePropertyInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the VehicleSignalsAvailableProperty.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *VehicleSignalsAvailableProperty) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	vehicleSignalsAvailablePropertyUpdateCacheMut.RLock()
	cache, cached := vehicleSignalsAvailablePropertyUpdateCache[key]
	vehicleSignalsAvailablePropertyUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			vehicleSignalsAvailablePropertyAllColumns,
			vehicleSignalsAvailablePropertyPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update vehicle_signals_available_properties, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"device_data_api\".\"vehicle_signals_available_properties\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, vehicleSignalsAvailablePropertyPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(vehicleSignalsAvailablePropertyType, vehicleSignalsAvailablePropertyMapping, append(wl, vehicleSignalsAvailablePropertyPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update vehicle_signals_available_properties row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for vehicle_signals_available_properties")
	}

	if !cached {
		vehicleSignalsAvailablePropertyUpdateCacheMut.Lock()
		vehicleSignalsAvailablePropertyUpdateCache[key] = cache
		vehicleSignalsAvailablePropertyUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q vehicleSignalsAvailablePropertyQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for vehicle_signals_available_properties")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for vehicle_signals_available_properties")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o VehicleSignalsAvailablePropertySlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), vehicleSignalsAvailablePropertyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"device_data_api\".\"vehicle_signals_available_properties\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, vehicleSignalsAvailablePropertyPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in vehicleSignalsAvailableProperty slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all vehicleSignalsAvailableProperty")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *VehicleSignalsAvailableProperty) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("models: no vehicle_signals_available_properties provided for upsert")
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

	nzDefaults := queries.NonZeroDefaultSet(vehicleSignalsAvailablePropertyColumnsWithDefault, o)

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

	vehicleSignalsAvailablePropertyUpsertCacheMut.RLock()
	cache, cached := vehicleSignalsAvailablePropertyUpsertCache[key]
	vehicleSignalsAvailablePropertyUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			vehicleSignalsAvailablePropertyAllColumns,
			vehicleSignalsAvailablePropertyColumnsWithDefault,
			vehicleSignalsAvailablePropertyColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			vehicleSignalsAvailablePropertyAllColumns,
			vehicleSignalsAvailablePropertyPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert vehicle_signals_available_properties, could not build update column list")
		}

		ret := strmangle.SetComplement(vehicleSignalsAvailablePropertyAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(vehicleSignalsAvailablePropertyPrimaryKeyColumns) == 0 {
				return errors.New("models: unable to upsert vehicle_signals_available_properties, could not build conflict column list")
			}

			conflict = make([]string, len(vehicleSignalsAvailablePropertyPrimaryKeyColumns))
			copy(conflict, vehicleSignalsAvailablePropertyPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"device_data_api\".\"vehicle_signals_available_properties\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(vehicleSignalsAvailablePropertyType, vehicleSignalsAvailablePropertyMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(vehicleSignalsAvailablePropertyType, vehicleSignalsAvailablePropertyMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert vehicle_signals_available_properties")
	}

	if !cached {
		vehicleSignalsAvailablePropertyUpsertCacheMut.Lock()
		vehicleSignalsAvailablePropertyUpsertCache[key] = cache
		vehicleSignalsAvailablePropertyUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single VehicleSignalsAvailableProperty record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *VehicleSignalsAvailableProperty) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no VehicleSignalsAvailableProperty provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), vehicleSignalsAvailablePropertyPrimaryKeyMapping)
	sql := "DELETE FROM \"device_data_api\".\"vehicle_signals_available_properties\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from vehicle_signals_available_properties")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for vehicle_signals_available_properties")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q vehicleSignalsAvailablePropertyQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no vehicleSignalsAvailablePropertyQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from vehicle_signals_available_properties")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for vehicle_signals_available_properties")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o VehicleSignalsAvailablePropertySlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(vehicleSignalsAvailablePropertyBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), vehicleSignalsAvailablePropertyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"device_data_api\".\"vehicle_signals_available_properties\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, vehicleSignalsAvailablePropertyPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from vehicleSignalsAvailableProperty slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for vehicle_signals_available_properties")
	}

	if len(vehicleSignalsAvailablePropertyAfterDeleteHooks) != 0 {
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
func (o *VehicleSignalsAvailableProperty) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindVehicleSignalsAvailableProperty(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *VehicleSignalsAvailablePropertySlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := VehicleSignalsAvailablePropertySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), vehicleSignalsAvailablePropertyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"device_data_api\".\"vehicle_signals_available_properties\".* FROM \"device_data_api\".\"vehicle_signals_available_properties\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, vehicleSignalsAvailablePropertyPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in VehicleSignalsAvailablePropertySlice")
	}

	*o = slice

	return nil
}

// VehicleSignalsAvailablePropertyExists checks if the VehicleSignalsAvailableProperty row exists.
func VehicleSignalsAvailablePropertyExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"device_data_api\".\"vehicle_signals_available_properties\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if vehicle_signals_available_properties exists")
	}

	return exists, nil
}

// Exists checks if the VehicleSignalsAvailableProperty row exists.
func (o *VehicleSignalsAvailableProperty) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return VehicleSignalsAvailablePropertyExists(ctx, exec, o.ID)
}
