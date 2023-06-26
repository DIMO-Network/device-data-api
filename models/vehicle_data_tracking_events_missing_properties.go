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

// VehicleDataTrackingEventsMissingProperty is an object representing the database table.
type VehicleDataTrackingEventsMissingProperty struct {
	IntegrationID string      `boil:"integration_id" json:"integration_id" toml:"integration_id" yaml:"integration_id"`
	DeviceMakeID  string      `boil:"device_make_id" json:"device_make_id" toml:"device_make_id" yaml:"device_make_id"`
	PropertyID    string      `boil:"property_id" json:"property_id" toml:"property_id" yaml:"property_id"`
	Model         string      `boil:"model" json:"model" toml:"model" yaml:"model"`
	Year          int         `boil:"year" json:"year" toml:"year" yaml:"year"`
	Description   null.String `boil:"description" json:"description,omitempty" toml:"description" yaml:"description,omitempty"`
	Count         int         `boil:"count" json:"count" toml:"count" yaml:"count"`
	CreatedAt     time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *vehicleDataTrackingEventsMissingPropertyR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L vehicleDataTrackingEventsMissingPropertyL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var VehicleDataTrackingEventsMissingPropertyColumns = struct {
	IntegrationID string
	DeviceMakeID  string
	PropertyID    string
	Model         string
	Year          string
	Description   string
	Count         string
	CreatedAt     string
}{
	IntegrationID: "integration_id",
	DeviceMakeID:  "device_make_id",
	PropertyID:    "property_id",
	Model:         "model",
	Year:          "year",
	Description:   "description",
	Count:         "count",
	CreatedAt:     "created_at",
}

var VehicleDataTrackingEventsMissingPropertyTableColumns = struct {
	IntegrationID string
	DeviceMakeID  string
	PropertyID    string
	Model         string
	Year          string
	Description   string
	Count         string
	CreatedAt     string
}{
	IntegrationID: "vehicle_data_tracking_events_missing_properties.integration_id",
	DeviceMakeID:  "vehicle_data_tracking_events_missing_properties.device_make_id",
	PropertyID:    "vehicle_data_tracking_events_missing_properties.property_id",
	Model:         "vehicle_data_tracking_events_missing_properties.model",
	Year:          "vehicle_data_tracking_events_missing_properties.year",
	Description:   "vehicle_data_tracking_events_missing_properties.description",
	Count:         "vehicle_data_tracking_events_missing_properties.count",
	CreatedAt:     "vehicle_data_tracking_events_missing_properties.created_at",
}

// Generated where

type whereHelperint struct{ field string }

func (w whereHelperint) EQ(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint) NEQ(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint) LT(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint) LTE(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint) GT(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint) GTE(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint) IN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperint) NIN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

var VehicleDataTrackingEventsMissingPropertyWhere = struct {
	IntegrationID whereHelperstring
	DeviceMakeID  whereHelperstring
	PropertyID    whereHelperstring
	Model         whereHelperstring
	Year          whereHelperint
	Description   whereHelpernull_String
	Count         whereHelperint
	CreatedAt     whereHelpertime_Time
}{
	IntegrationID: whereHelperstring{field: "\"device_data_api\".\"vehicle_data_tracking_events_missing_properties\".\"integration_id\""},
	DeviceMakeID:  whereHelperstring{field: "\"device_data_api\".\"vehicle_data_tracking_events_missing_properties\".\"device_make_id\""},
	PropertyID:    whereHelperstring{field: "\"device_data_api\".\"vehicle_data_tracking_events_missing_properties\".\"property_id\""},
	Model:         whereHelperstring{field: "\"device_data_api\".\"vehicle_data_tracking_events_missing_properties\".\"model\""},
	Year:          whereHelperint{field: "\"device_data_api\".\"vehicle_data_tracking_events_missing_properties\".\"year\""},
	Description:   whereHelpernull_String{field: "\"device_data_api\".\"vehicle_data_tracking_events_missing_properties\".\"description\""},
	Count:         whereHelperint{field: "\"device_data_api\".\"vehicle_data_tracking_events_missing_properties\".\"count\""},
	CreatedAt:     whereHelpertime_Time{field: "\"device_data_api\".\"vehicle_data_tracking_events_missing_properties\".\"created_at\""},
}

// VehicleDataTrackingEventsMissingPropertyRels is where relationship names are stored.
var VehicleDataTrackingEventsMissingPropertyRels = struct {
}{}

// vehicleDataTrackingEventsMissingPropertyR is where relationships are stored.
type vehicleDataTrackingEventsMissingPropertyR struct {
}

// NewStruct creates a new relationship struct
func (*vehicleDataTrackingEventsMissingPropertyR) NewStruct() *vehicleDataTrackingEventsMissingPropertyR {
	return &vehicleDataTrackingEventsMissingPropertyR{}
}

// vehicleDataTrackingEventsMissingPropertyL is where Load methods for each relationship are stored.
type vehicleDataTrackingEventsMissingPropertyL struct{}

var (
	vehicleDataTrackingEventsMissingPropertyAllColumns            = []string{"integration_id", "device_make_id", "property_id", "model", "year", "description", "count", "created_at"}
	vehicleDataTrackingEventsMissingPropertyColumnsWithoutDefault = []string{"integration_id", "device_make_id", "property_id", "model", "year", "count"}
	vehicleDataTrackingEventsMissingPropertyColumnsWithDefault    = []string{"description", "created_at"}
	vehicleDataTrackingEventsMissingPropertyPrimaryKeyColumns     = []string{"integration_id", "device_make_id", "property_id"}
	vehicleDataTrackingEventsMissingPropertyGeneratedColumns      = []string{}
)

type (
	// VehicleDataTrackingEventsMissingPropertySlice is an alias for a slice of pointers to VehicleDataTrackingEventsMissingProperty.
	// This should almost always be used instead of []VehicleDataTrackingEventsMissingProperty.
	VehicleDataTrackingEventsMissingPropertySlice []*VehicleDataTrackingEventsMissingProperty
	// VehicleDataTrackingEventsMissingPropertyHook is the signature for custom VehicleDataTrackingEventsMissingProperty hook methods
	VehicleDataTrackingEventsMissingPropertyHook func(context.Context, boil.ContextExecutor, *VehicleDataTrackingEventsMissingProperty) error

	vehicleDataTrackingEventsMissingPropertyQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	vehicleDataTrackingEventsMissingPropertyType                 = reflect.TypeOf(&VehicleDataTrackingEventsMissingProperty{})
	vehicleDataTrackingEventsMissingPropertyMapping              = queries.MakeStructMapping(vehicleDataTrackingEventsMissingPropertyType)
	vehicleDataTrackingEventsMissingPropertyPrimaryKeyMapping, _ = queries.BindMapping(vehicleDataTrackingEventsMissingPropertyType, vehicleDataTrackingEventsMissingPropertyMapping, vehicleDataTrackingEventsMissingPropertyPrimaryKeyColumns)
	vehicleDataTrackingEventsMissingPropertyInsertCacheMut       sync.RWMutex
	vehicleDataTrackingEventsMissingPropertyInsertCache          = make(map[string]insertCache)
	vehicleDataTrackingEventsMissingPropertyUpdateCacheMut       sync.RWMutex
	vehicleDataTrackingEventsMissingPropertyUpdateCache          = make(map[string]updateCache)
	vehicleDataTrackingEventsMissingPropertyUpsertCacheMut       sync.RWMutex
	vehicleDataTrackingEventsMissingPropertyUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var vehicleDataTrackingEventsMissingPropertyAfterSelectHooks []VehicleDataTrackingEventsMissingPropertyHook

var vehicleDataTrackingEventsMissingPropertyBeforeInsertHooks []VehicleDataTrackingEventsMissingPropertyHook
var vehicleDataTrackingEventsMissingPropertyAfterInsertHooks []VehicleDataTrackingEventsMissingPropertyHook

var vehicleDataTrackingEventsMissingPropertyBeforeUpdateHooks []VehicleDataTrackingEventsMissingPropertyHook
var vehicleDataTrackingEventsMissingPropertyAfterUpdateHooks []VehicleDataTrackingEventsMissingPropertyHook

var vehicleDataTrackingEventsMissingPropertyBeforeDeleteHooks []VehicleDataTrackingEventsMissingPropertyHook
var vehicleDataTrackingEventsMissingPropertyAfterDeleteHooks []VehicleDataTrackingEventsMissingPropertyHook

var vehicleDataTrackingEventsMissingPropertyBeforeUpsertHooks []VehicleDataTrackingEventsMissingPropertyHook
var vehicleDataTrackingEventsMissingPropertyAfterUpsertHooks []VehicleDataTrackingEventsMissingPropertyHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *VehicleDataTrackingEventsMissingProperty) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleDataTrackingEventsMissingPropertyAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *VehicleDataTrackingEventsMissingProperty) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleDataTrackingEventsMissingPropertyBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *VehicleDataTrackingEventsMissingProperty) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleDataTrackingEventsMissingPropertyAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *VehicleDataTrackingEventsMissingProperty) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleDataTrackingEventsMissingPropertyBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *VehicleDataTrackingEventsMissingProperty) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleDataTrackingEventsMissingPropertyAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *VehicleDataTrackingEventsMissingProperty) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleDataTrackingEventsMissingPropertyBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *VehicleDataTrackingEventsMissingProperty) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleDataTrackingEventsMissingPropertyAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *VehicleDataTrackingEventsMissingProperty) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleDataTrackingEventsMissingPropertyBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *VehicleDataTrackingEventsMissingProperty) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range vehicleDataTrackingEventsMissingPropertyAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddVehicleDataTrackingEventsMissingPropertyHook registers your hook function for all future operations.
func AddVehicleDataTrackingEventsMissingPropertyHook(hookPoint boil.HookPoint, vehicleDataTrackingEventsMissingPropertyHook VehicleDataTrackingEventsMissingPropertyHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		vehicleDataTrackingEventsMissingPropertyAfterSelectHooks = append(vehicleDataTrackingEventsMissingPropertyAfterSelectHooks, vehicleDataTrackingEventsMissingPropertyHook)
	case boil.BeforeInsertHook:
		vehicleDataTrackingEventsMissingPropertyBeforeInsertHooks = append(vehicleDataTrackingEventsMissingPropertyBeforeInsertHooks, vehicleDataTrackingEventsMissingPropertyHook)
	case boil.AfterInsertHook:
		vehicleDataTrackingEventsMissingPropertyAfterInsertHooks = append(vehicleDataTrackingEventsMissingPropertyAfterInsertHooks, vehicleDataTrackingEventsMissingPropertyHook)
	case boil.BeforeUpdateHook:
		vehicleDataTrackingEventsMissingPropertyBeforeUpdateHooks = append(vehicleDataTrackingEventsMissingPropertyBeforeUpdateHooks, vehicleDataTrackingEventsMissingPropertyHook)
	case boil.AfterUpdateHook:
		vehicleDataTrackingEventsMissingPropertyAfterUpdateHooks = append(vehicleDataTrackingEventsMissingPropertyAfterUpdateHooks, vehicleDataTrackingEventsMissingPropertyHook)
	case boil.BeforeDeleteHook:
		vehicleDataTrackingEventsMissingPropertyBeforeDeleteHooks = append(vehicleDataTrackingEventsMissingPropertyBeforeDeleteHooks, vehicleDataTrackingEventsMissingPropertyHook)
	case boil.AfterDeleteHook:
		vehicleDataTrackingEventsMissingPropertyAfterDeleteHooks = append(vehicleDataTrackingEventsMissingPropertyAfterDeleteHooks, vehicleDataTrackingEventsMissingPropertyHook)
	case boil.BeforeUpsertHook:
		vehicleDataTrackingEventsMissingPropertyBeforeUpsertHooks = append(vehicleDataTrackingEventsMissingPropertyBeforeUpsertHooks, vehicleDataTrackingEventsMissingPropertyHook)
	case boil.AfterUpsertHook:
		vehicleDataTrackingEventsMissingPropertyAfterUpsertHooks = append(vehicleDataTrackingEventsMissingPropertyAfterUpsertHooks, vehicleDataTrackingEventsMissingPropertyHook)
	}
}

// One returns a single vehicleDataTrackingEventsMissingProperty record from the query.
func (q vehicleDataTrackingEventsMissingPropertyQuery) One(ctx context.Context, exec boil.ContextExecutor) (*VehicleDataTrackingEventsMissingProperty, error) {
	o := &VehicleDataTrackingEventsMissingProperty{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for vehicle_data_tracking_events_missing_properties")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all VehicleDataTrackingEventsMissingProperty records from the query.
func (q vehicleDataTrackingEventsMissingPropertyQuery) All(ctx context.Context, exec boil.ContextExecutor) (VehicleDataTrackingEventsMissingPropertySlice, error) {
	var o []*VehicleDataTrackingEventsMissingProperty

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to VehicleDataTrackingEventsMissingProperty slice")
	}

	if len(vehicleDataTrackingEventsMissingPropertyAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all VehicleDataTrackingEventsMissingProperty records in the query.
func (q vehicleDataTrackingEventsMissingPropertyQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count vehicle_data_tracking_events_missing_properties rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q vehicleDataTrackingEventsMissingPropertyQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if vehicle_data_tracking_events_missing_properties exists")
	}

	return count > 0, nil
}

// VehicleDataTrackingEventsMissingProperties retrieves all the records using an executor.
func VehicleDataTrackingEventsMissingProperties(mods ...qm.QueryMod) vehicleDataTrackingEventsMissingPropertyQuery {
	mods = append(mods, qm.From("\"device_data_api\".\"vehicle_data_tracking_events_missing_properties\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"device_data_api\".\"vehicle_data_tracking_events_missing_properties\".*"})
	}

	return vehicleDataTrackingEventsMissingPropertyQuery{q}
}

// FindVehicleDataTrackingEventsMissingProperty retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindVehicleDataTrackingEventsMissingProperty(ctx context.Context, exec boil.ContextExecutor, integrationID string, deviceMakeID string, propertyID string, selectCols ...string) (*VehicleDataTrackingEventsMissingProperty, error) {
	vehicleDataTrackingEventsMissingPropertyObj := &VehicleDataTrackingEventsMissingProperty{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"device_data_api\".\"vehicle_data_tracking_events_missing_properties\" where \"integration_id\"=$1 AND \"device_make_id\"=$2 AND \"property_id\"=$3", sel,
	)

	q := queries.Raw(query, integrationID, deviceMakeID, propertyID)

	err := q.Bind(ctx, exec, vehicleDataTrackingEventsMissingPropertyObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from vehicle_data_tracking_events_missing_properties")
	}

	if err = vehicleDataTrackingEventsMissingPropertyObj.doAfterSelectHooks(ctx, exec); err != nil {
		return vehicleDataTrackingEventsMissingPropertyObj, err
	}

	return vehicleDataTrackingEventsMissingPropertyObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *VehicleDataTrackingEventsMissingProperty) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no vehicle_data_tracking_events_missing_properties provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(vehicleDataTrackingEventsMissingPropertyColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	vehicleDataTrackingEventsMissingPropertyInsertCacheMut.RLock()
	cache, cached := vehicleDataTrackingEventsMissingPropertyInsertCache[key]
	vehicleDataTrackingEventsMissingPropertyInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			vehicleDataTrackingEventsMissingPropertyAllColumns,
			vehicleDataTrackingEventsMissingPropertyColumnsWithDefault,
			vehicleDataTrackingEventsMissingPropertyColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(vehicleDataTrackingEventsMissingPropertyType, vehicleDataTrackingEventsMissingPropertyMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(vehicleDataTrackingEventsMissingPropertyType, vehicleDataTrackingEventsMissingPropertyMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"device_data_api\".\"vehicle_data_tracking_events_missing_properties\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"device_data_api\".\"vehicle_data_tracking_events_missing_properties\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into vehicle_data_tracking_events_missing_properties")
	}

	if !cached {
		vehicleDataTrackingEventsMissingPropertyInsertCacheMut.Lock()
		vehicleDataTrackingEventsMissingPropertyInsertCache[key] = cache
		vehicleDataTrackingEventsMissingPropertyInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the VehicleDataTrackingEventsMissingProperty.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *VehicleDataTrackingEventsMissingProperty) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	vehicleDataTrackingEventsMissingPropertyUpdateCacheMut.RLock()
	cache, cached := vehicleDataTrackingEventsMissingPropertyUpdateCache[key]
	vehicleDataTrackingEventsMissingPropertyUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			vehicleDataTrackingEventsMissingPropertyAllColumns,
			vehicleDataTrackingEventsMissingPropertyPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update vehicle_data_tracking_events_missing_properties, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"device_data_api\".\"vehicle_data_tracking_events_missing_properties\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, vehicleDataTrackingEventsMissingPropertyPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(vehicleDataTrackingEventsMissingPropertyType, vehicleDataTrackingEventsMissingPropertyMapping, append(wl, vehicleDataTrackingEventsMissingPropertyPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update vehicle_data_tracking_events_missing_properties row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for vehicle_data_tracking_events_missing_properties")
	}

	if !cached {
		vehicleDataTrackingEventsMissingPropertyUpdateCacheMut.Lock()
		vehicleDataTrackingEventsMissingPropertyUpdateCache[key] = cache
		vehicleDataTrackingEventsMissingPropertyUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q vehicleDataTrackingEventsMissingPropertyQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for vehicle_data_tracking_events_missing_properties")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for vehicle_data_tracking_events_missing_properties")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o VehicleDataTrackingEventsMissingPropertySlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), vehicleDataTrackingEventsMissingPropertyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"device_data_api\".\"vehicle_data_tracking_events_missing_properties\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, vehicleDataTrackingEventsMissingPropertyPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in vehicleDataTrackingEventsMissingProperty slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all vehicleDataTrackingEventsMissingProperty")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *VehicleDataTrackingEventsMissingProperty) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no vehicle_data_tracking_events_missing_properties provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(vehicleDataTrackingEventsMissingPropertyColumnsWithDefault, o)

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

	vehicleDataTrackingEventsMissingPropertyUpsertCacheMut.RLock()
	cache, cached := vehicleDataTrackingEventsMissingPropertyUpsertCache[key]
	vehicleDataTrackingEventsMissingPropertyUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			vehicleDataTrackingEventsMissingPropertyAllColumns,
			vehicleDataTrackingEventsMissingPropertyColumnsWithDefault,
			vehicleDataTrackingEventsMissingPropertyColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			vehicleDataTrackingEventsMissingPropertyAllColumns,
			vehicleDataTrackingEventsMissingPropertyPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert vehicle_data_tracking_events_missing_properties, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(vehicleDataTrackingEventsMissingPropertyPrimaryKeyColumns))
			copy(conflict, vehicleDataTrackingEventsMissingPropertyPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"device_data_api\".\"vehicle_data_tracking_events_missing_properties\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(vehicleDataTrackingEventsMissingPropertyType, vehicleDataTrackingEventsMissingPropertyMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(vehicleDataTrackingEventsMissingPropertyType, vehicleDataTrackingEventsMissingPropertyMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert vehicle_data_tracking_events_missing_properties")
	}

	if !cached {
		vehicleDataTrackingEventsMissingPropertyUpsertCacheMut.Lock()
		vehicleDataTrackingEventsMissingPropertyUpsertCache[key] = cache
		vehicleDataTrackingEventsMissingPropertyUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single VehicleDataTrackingEventsMissingProperty record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *VehicleDataTrackingEventsMissingProperty) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no VehicleDataTrackingEventsMissingProperty provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), vehicleDataTrackingEventsMissingPropertyPrimaryKeyMapping)
	sql := "DELETE FROM \"device_data_api\".\"vehicle_data_tracking_events_missing_properties\" WHERE \"integration_id\"=$1 AND \"device_make_id\"=$2 AND \"property_id\"=$3"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from vehicle_data_tracking_events_missing_properties")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for vehicle_data_tracking_events_missing_properties")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q vehicleDataTrackingEventsMissingPropertyQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no vehicleDataTrackingEventsMissingPropertyQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from vehicle_data_tracking_events_missing_properties")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for vehicle_data_tracking_events_missing_properties")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o VehicleDataTrackingEventsMissingPropertySlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(vehicleDataTrackingEventsMissingPropertyBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), vehicleDataTrackingEventsMissingPropertyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"device_data_api\".\"vehicle_data_tracking_events_missing_properties\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, vehicleDataTrackingEventsMissingPropertyPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from vehicleDataTrackingEventsMissingProperty slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for vehicle_data_tracking_events_missing_properties")
	}

	if len(vehicleDataTrackingEventsMissingPropertyAfterDeleteHooks) != 0 {
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
func (o *VehicleDataTrackingEventsMissingProperty) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindVehicleDataTrackingEventsMissingProperty(ctx, exec, o.IntegrationID, o.DeviceMakeID, o.PropertyID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *VehicleDataTrackingEventsMissingPropertySlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := VehicleDataTrackingEventsMissingPropertySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), vehicleDataTrackingEventsMissingPropertyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"device_data_api\".\"vehicle_data_tracking_events_missing_properties\".* FROM \"device_data_api\".\"vehicle_data_tracking_events_missing_properties\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, vehicleDataTrackingEventsMissingPropertyPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in VehicleDataTrackingEventsMissingPropertySlice")
	}

	*o = slice

	return nil
}

// VehicleDataTrackingEventsMissingPropertyExists checks if the VehicleDataTrackingEventsMissingProperty row exists.
func VehicleDataTrackingEventsMissingPropertyExists(ctx context.Context, exec boil.ContextExecutor, integrationID string, deviceMakeID string, propertyID string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"device_data_api\".\"vehicle_data_tracking_events_missing_properties\" where \"integration_id\"=$1 AND \"device_make_id\"=$2 AND \"property_id\"=$3 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, integrationID, deviceMakeID, propertyID)
	}
	row := exec.QueryRowContext(ctx, sql, integrationID, deviceMakeID, propertyID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if vehicle_data_tracking_events_missing_properties exists")
	}

	return exists, nil
}
