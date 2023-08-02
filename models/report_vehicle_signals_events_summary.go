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
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// ReportVehicleSignalsEventsSummary is an object representing the database table.
type ReportVehicleSignalsEventsSummary struct {
	DateID             string    `boil:"date_id" json:"date_id" toml:"date_id" yaml:"date_id"`
	IntegrationID      string    `boil:"integration_id" json:"integration_id" toml:"integration_id" yaml:"integration_id"`
	PowerTrainType     string    `boil:"power_train_type" json:"power_train_type" toml:"power_train_type" yaml:"power_train_type"`
	Count              int       `boil:"count" json:"count" toml:"count" yaml:"count"`
	CreatedAt          time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	DeviceDefinitionID string    `boil:"device_definition_id" json:"device_definition_id" toml:"device_definition_id" yaml:"device_definition_id"`

	R *reportVehicleSignalsEventsSummaryR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L reportVehicleSignalsEventsSummaryL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ReportVehicleSignalsEventsSummaryColumns = struct {
	DateID             string
	IntegrationID      string
	PowerTrainType     string
	Count              string
	CreatedAt          string
	DeviceDefinitionID string
}{
	DateID:             "date_id",
	IntegrationID:      "integration_id",
	PowerTrainType:     "power_train_type",
	Count:              "count",
	CreatedAt:          "created_at",
	DeviceDefinitionID: "device_definition_id",
}

var ReportVehicleSignalsEventsSummaryTableColumns = struct {
	DateID             string
	IntegrationID      string
	PowerTrainType     string
	Count              string
	CreatedAt          string
	DeviceDefinitionID string
}{
	DateID:             "report_vehicle_signals_events_summary.date_id",
	IntegrationID:      "report_vehicle_signals_events_summary.integration_id",
	PowerTrainType:     "report_vehicle_signals_events_summary.power_train_type",
	Count:              "report_vehicle_signals_events_summary.count",
	CreatedAt:          "report_vehicle_signals_events_summary.created_at",
	DeviceDefinitionID: "report_vehicle_signals_events_summary.device_definition_id",
}

// Generated where

var ReportVehicleSignalsEventsSummaryWhere = struct {
	DateID             whereHelperstring
	IntegrationID      whereHelperstring
	PowerTrainType     whereHelperstring
	Count              whereHelperint
	CreatedAt          whereHelpertime_Time
	DeviceDefinitionID whereHelperstring
}{
	DateID:             whereHelperstring{field: "\"device_data_api\".\"report_vehicle_signals_events_summary\".\"date_id\""},
	IntegrationID:      whereHelperstring{field: "\"device_data_api\".\"report_vehicle_signals_events_summary\".\"integration_id\""},
	PowerTrainType:     whereHelperstring{field: "\"device_data_api\".\"report_vehicle_signals_events_summary\".\"power_train_type\""},
	Count:              whereHelperint{field: "\"device_data_api\".\"report_vehicle_signals_events_summary\".\"count\""},
	CreatedAt:          whereHelpertime_Time{field: "\"device_data_api\".\"report_vehicle_signals_events_summary\".\"created_at\""},
	DeviceDefinitionID: whereHelperstring{field: "\"device_data_api\".\"report_vehicle_signals_events_summary\".\"device_definition_id\""},
}

// ReportVehicleSignalsEventsSummaryRels is where relationship names are stored.
var ReportVehicleSignalsEventsSummaryRels = struct {
}{}

// reportVehicleSignalsEventsSummaryR is where relationships are stored.
type reportVehicleSignalsEventsSummaryR struct {
}

// NewStruct creates a new relationship struct
func (*reportVehicleSignalsEventsSummaryR) NewStruct() *reportVehicleSignalsEventsSummaryR {
	return &reportVehicleSignalsEventsSummaryR{}
}

// reportVehicleSignalsEventsSummaryL is where Load methods for each relationship are stored.
type reportVehicleSignalsEventsSummaryL struct{}

var (
	reportVehicleSignalsEventsSummaryAllColumns            = []string{"date_id", "integration_id", "power_train_type", "count", "created_at", "device_definition_id"}
	reportVehicleSignalsEventsSummaryColumnsWithoutDefault = []string{"date_id", "integration_id", "power_train_type", "count"}
	reportVehicleSignalsEventsSummaryColumnsWithDefault    = []string{"created_at", "device_definition_id"}
	reportVehicleSignalsEventsSummaryPrimaryKeyColumns     = []string{"date_id", "integration_id", "power_train_type", "device_definition_id"}
	reportVehicleSignalsEventsSummaryGeneratedColumns      = []string{}
)

type (
	// ReportVehicleSignalsEventsSummarySlice is an alias for a slice of pointers to ReportVehicleSignalsEventsSummary.
	// This should almost always be used instead of []ReportVehicleSignalsEventsSummary.
	ReportVehicleSignalsEventsSummarySlice []*ReportVehicleSignalsEventsSummary
	// ReportVehicleSignalsEventsSummaryHook is the signature for custom ReportVehicleSignalsEventsSummary hook methods
	ReportVehicleSignalsEventsSummaryHook func(context.Context, boil.ContextExecutor, *ReportVehicleSignalsEventsSummary) error

	reportVehicleSignalsEventsSummaryQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	reportVehicleSignalsEventsSummaryType                 = reflect.TypeOf(&ReportVehicleSignalsEventsSummary{})
	reportVehicleSignalsEventsSummaryMapping              = queries.MakeStructMapping(reportVehicleSignalsEventsSummaryType)
	reportVehicleSignalsEventsSummaryPrimaryKeyMapping, _ = queries.BindMapping(reportVehicleSignalsEventsSummaryType, reportVehicleSignalsEventsSummaryMapping, reportVehicleSignalsEventsSummaryPrimaryKeyColumns)
	reportVehicleSignalsEventsSummaryInsertCacheMut       sync.RWMutex
	reportVehicleSignalsEventsSummaryInsertCache          = make(map[string]insertCache)
	reportVehicleSignalsEventsSummaryUpdateCacheMut       sync.RWMutex
	reportVehicleSignalsEventsSummaryUpdateCache          = make(map[string]updateCache)
	reportVehicleSignalsEventsSummaryUpsertCacheMut       sync.RWMutex
	reportVehicleSignalsEventsSummaryUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var reportVehicleSignalsEventsSummaryAfterSelectHooks []ReportVehicleSignalsEventsSummaryHook

var reportVehicleSignalsEventsSummaryBeforeInsertHooks []ReportVehicleSignalsEventsSummaryHook
var reportVehicleSignalsEventsSummaryAfterInsertHooks []ReportVehicleSignalsEventsSummaryHook

var reportVehicleSignalsEventsSummaryBeforeUpdateHooks []ReportVehicleSignalsEventsSummaryHook
var reportVehicleSignalsEventsSummaryAfterUpdateHooks []ReportVehicleSignalsEventsSummaryHook

var reportVehicleSignalsEventsSummaryBeforeDeleteHooks []ReportVehicleSignalsEventsSummaryHook
var reportVehicleSignalsEventsSummaryAfterDeleteHooks []ReportVehicleSignalsEventsSummaryHook

var reportVehicleSignalsEventsSummaryBeforeUpsertHooks []ReportVehicleSignalsEventsSummaryHook
var reportVehicleSignalsEventsSummaryAfterUpsertHooks []ReportVehicleSignalsEventsSummaryHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *ReportVehicleSignalsEventsSummary) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range reportVehicleSignalsEventsSummaryAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *ReportVehicleSignalsEventsSummary) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range reportVehicleSignalsEventsSummaryBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *ReportVehicleSignalsEventsSummary) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range reportVehicleSignalsEventsSummaryAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *ReportVehicleSignalsEventsSummary) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range reportVehicleSignalsEventsSummaryBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *ReportVehicleSignalsEventsSummary) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range reportVehicleSignalsEventsSummaryAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *ReportVehicleSignalsEventsSummary) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range reportVehicleSignalsEventsSummaryBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *ReportVehicleSignalsEventsSummary) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range reportVehicleSignalsEventsSummaryAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *ReportVehicleSignalsEventsSummary) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range reportVehicleSignalsEventsSummaryBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *ReportVehicleSignalsEventsSummary) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range reportVehicleSignalsEventsSummaryAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddReportVehicleSignalsEventsSummaryHook registers your hook function for all future operations.
func AddReportVehicleSignalsEventsSummaryHook(hookPoint boil.HookPoint, reportVehicleSignalsEventsSummaryHook ReportVehicleSignalsEventsSummaryHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		reportVehicleSignalsEventsSummaryAfterSelectHooks = append(reportVehicleSignalsEventsSummaryAfterSelectHooks, reportVehicleSignalsEventsSummaryHook)
	case boil.BeforeInsertHook:
		reportVehicleSignalsEventsSummaryBeforeInsertHooks = append(reportVehicleSignalsEventsSummaryBeforeInsertHooks, reportVehicleSignalsEventsSummaryHook)
	case boil.AfterInsertHook:
		reportVehicleSignalsEventsSummaryAfterInsertHooks = append(reportVehicleSignalsEventsSummaryAfterInsertHooks, reportVehicleSignalsEventsSummaryHook)
	case boil.BeforeUpdateHook:
		reportVehicleSignalsEventsSummaryBeforeUpdateHooks = append(reportVehicleSignalsEventsSummaryBeforeUpdateHooks, reportVehicleSignalsEventsSummaryHook)
	case boil.AfterUpdateHook:
		reportVehicleSignalsEventsSummaryAfterUpdateHooks = append(reportVehicleSignalsEventsSummaryAfterUpdateHooks, reportVehicleSignalsEventsSummaryHook)
	case boil.BeforeDeleteHook:
		reportVehicleSignalsEventsSummaryBeforeDeleteHooks = append(reportVehicleSignalsEventsSummaryBeforeDeleteHooks, reportVehicleSignalsEventsSummaryHook)
	case boil.AfterDeleteHook:
		reportVehicleSignalsEventsSummaryAfterDeleteHooks = append(reportVehicleSignalsEventsSummaryAfterDeleteHooks, reportVehicleSignalsEventsSummaryHook)
	case boil.BeforeUpsertHook:
		reportVehicleSignalsEventsSummaryBeforeUpsertHooks = append(reportVehicleSignalsEventsSummaryBeforeUpsertHooks, reportVehicleSignalsEventsSummaryHook)
	case boil.AfterUpsertHook:
		reportVehicleSignalsEventsSummaryAfterUpsertHooks = append(reportVehicleSignalsEventsSummaryAfterUpsertHooks, reportVehicleSignalsEventsSummaryHook)
	}
}

// One returns a single reportVehicleSignalsEventsSummary record from the query.
func (q reportVehicleSignalsEventsSummaryQuery) One(ctx context.Context, exec boil.ContextExecutor) (*ReportVehicleSignalsEventsSummary, error) {
	o := &ReportVehicleSignalsEventsSummary{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for report_vehicle_signals_events_summary")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all ReportVehicleSignalsEventsSummary records from the query.
func (q reportVehicleSignalsEventsSummaryQuery) All(ctx context.Context, exec boil.ContextExecutor) (ReportVehicleSignalsEventsSummarySlice, error) {
	var o []*ReportVehicleSignalsEventsSummary

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to ReportVehicleSignalsEventsSummary slice")
	}

	if len(reportVehicleSignalsEventsSummaryAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all ReportVehicleSignalsEventsSummary records in the query.
func (q reportVehicleSignalsEventsSummaryQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count report_vehicle_signals_events_summary rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q reportVehicleSignalsEventsSummaryQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if report_vehicle_signals_events_summary exists")
	}

	return count > 0, nil
}

// ReportVehicleSignalsEventsSummaries retrieves all the records using an executor.
func ReportVehicleSignalsEventsSummaries(mods ...qm.QueryMod) reportVehicleSignalsEventsSummaryQuery {
	mods = append(mods, qm.From("\"device_data_api\".\"report_vehicle_signals_events_summary\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"device_data_api\".\"report_vehicle_signals_events_summary\".*"})
	}

	return reportVehicleSignalsEventsSummaryQuery{q}
}

// FindReportVehicleSignalsEventsSummary retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindReportVehicleSignalsEventsSummary(ctx context.Context, exec boil.ContextExecutor, dateID string, integrationID string, powerTrainType string, deviceDefinitionID string, selectCols ...string) (*ReportVehicleSignalsEventsSummary, error) {
	reportVehicleSignalsEventsSummaryObj := &ReportVehicleSignalsEventsSummary{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"device_data_api\".\"report_vehicle_signals_events_summary\" where \"date_id\"=$1 AND \"integration_id\"=$2 AND \"power_train_type\"=$3 AND \"device_definition_id\"=$4", sel,
	)

	q := queries.Raw(query, dateID, integrationID, powerTrainType, deviceDefinitionID)

	err := q.Bind(ctx, exec, reportVehicleSignalsEventsSummaryObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from report_vehicle_signals_events_summary")
	}

	if err = reportVehicleSignalsEventsSummaryObj.doAfterSelectHooks(ctx, exec); err != nil {
		return reportVehicleSignalsEventsSummaryObj, err
	}

	return reportVehicleSignalsEventsSummaryObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *ReportVehicleSignalsEventsSummary) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no report_vehicle_signals_events_summary provided for insertion")
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

	nzDefaults := queries.NonZeroDefaultSet(reportVehicleSignalsEventsSummaryColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	reportVehicleSignalsEventsSummaryInsertCacheMut.RLock()
	cache, cached := reportVehicleSignalsEventsSummaryInsertCache[key]
	reportVehicleSignalsEventsSummaryInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			reportVehicleSignalsEventsSummaryAllColumns,
			reportVehicleSignalsEventsSummaryColumnsWithDefault,
			reportVehicleSignalsEventsSummaryColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(reportVehicleSignalsEventsSummaryType, reportVehicleSignalsEventsSummaryMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(reportVehicleSignalsEventsSummaryType, reportVehicleSignalsEventsSummaryMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"device_data_api\".\"report_vehicle_signals_events_summary\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"device_data_api\".\"report_vehicle_signals_events_summary\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into report_vehicle_signals_events_summary")
	}

	if !cached {
		reportVehicleSignalsEventsSummaryInsertCacheMut.Lock()
		reportVehicleSignalsEventsSummaryInsertCache[key] = cache
		reportVehicleSignalsEventsSummaryInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the ReportVehicleSignalsEventsSummary.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *ReportVehicleSignalsEventsSummary) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	reportVehicleSignalsEventsSummaryUpdateCacheMut.RLock()
	cache, cached := reportVehicleSignalsEventsSummaryUpdateCache[key]
	reportVehicleSignalsEventsSummaryUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			reportVehicleSignalsEventsSummaryAllColumns,
			reportVehicleSignalsEventsSummaryPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update report_vehicle_signals_events_summary, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"device_data_api\".\"report_vehicle_signals_events_summary\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, reportVehicleSignalsEventsSummaryPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(reportVehicleSignalsEventsSummaryType, reportVehicleSignalsEventsSummaryMapping, append(wl, reportVehicleSignalsEventsSummaryPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update report_vehicle_signals_events_summary row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for report_vehicle_signals_events_summary")
	}

	if !cached {
		reportVehicleSignalsEventsSummaryUpdateCacheMut.Lock()
		reportVehicleSignalsEventsSummaryUpdateCache[key] = cache
		reportVehicleSignalsEventsSummaryUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q reportVehicleSignalsEventsSummaryQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for report_vehicle_signals_events_summary")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for report_vehicle_signals_events_summary")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ReportVehicleSignalsEventsSummarySlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), reportVehicleSignalsEventsSummaryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"device_data_api\".\"report_vehicle_signals_events_summary\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, reportVehicleSignalsEventsSummaryPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in reportVehicleSignalsEventsSummary slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all reportVehicleSignalsEventsSummary")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *ReportVehicleSignalsEventsSummary) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no report_vehicle_signals_events_summary provided for upsert")
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

	nzDefaults := queries.NonZeroDefaultSet(reportVehicleSignalsEventsSummaryColumnsWithDefault, o)

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

	reportVehicleSignalsEventsSummaryUpsertCacheMut.RLock()
	cache, cached := reportVehicleSignalsEventsSummaryUpsertCache[key]
	reportVehicleSignalsEventsSummaryUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			reportVehicleSignalsEventsSummaryAllColumns,
			reportVehicleSignalsEventsSummaryColumnsWithDefault,
			reportVehicleSignalsEventsSummaryColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			reportVehicleSignalsEventsSummaryAllColumns,
			reportVehicleSignalsEventsSummaryPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert report_vehicle_signals_events_summary, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(reportVehicleSignalsEventsSummaryPrimaryKeyColumns))
			copy(conflict, reportVehicleSignalsEventsSummaryPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"device_data_api\".\"report_vehicle_signals_events_summary\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(reportVehicleSignalsEventsSummaryType, reportVehicleSignalsEventsSummaryMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(reportVehicleSignalsEventsSummaryType, reportVehicleSignalsEventsSummaryMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert report_vehicle_signals_events_summary")
	}

	if !cached {
		reportVehicleSignalsEventsSummaryUpsertCacheMut.Lock()
		reportVehicleSignalsEventsSummaryUpsertCache[key] = cache
		reportVehicleSignalsEventsSummaryUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single ReportVehicleSignalsEventsSummary record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *ReportVehicleSignalsEventsSummary) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no ReportVehicleSignalsEventsSummary provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), reportVehicleSignalsEventsSummaryPrimaryKeyMapping)
	sql := "DELETE FROM \"device_data_api\".\"report_vehicle_signals_events_summary\" WHERE \"date_id\"=$1 AND \"integration_id\"=$2 AND \"power_train_type\"=$3 AND \"device_definition_id\"=$4"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from report_vehicle_signals_events_summary")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for report_vehicle_signals_events_summary")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q reportVehicleSignalsEventsSummaryQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no reportVehicleSignalsEventsSummaryQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from report_vehicle_signals_events_summary")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for report_vehicle_signals_events_summary")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ReportVehicleSignalsEventsSummarySlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(reportVehicleSignalsEventsSummaryBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), reportVehicleSignalsEventsSummaryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"device_data_api\".\"report_vehicle_signals_events_summary\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, reportVehicleSignalsEventsSummaryPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from reportVehicleSignalsEventsSummary slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for report_vehicle_signals_events_summary")
	}

	if len(reportVehicleSignalsEventsSummaryAfterDeleteHooks) != 0 {
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
func (o *ReportVehicleSignalsEventsSummary) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindReportVehicleSignalsEventsSummary(ctx, exec, o.DateID, o.IntegrationID, o.PowerTrainType, o.DeviceDefinitionID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ReportVehicleSignalsEventsSummarySlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ReportVehicleSignalsEventsSummarySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), reportVehicleSignalsEventsSummaryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"device_data_api\".\"report_vehicle_signals_events_summary\".* FROM \"device_data_api\".\"report_vehicle_signals_events_summary\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, reportVehicleSignalsEventsSummaryPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ReportVehicleSignalsEventsSummarySlice")
	}

	*o = slice

	return nil
}

// ReportVehicleSignalsEventsSummaryExists checks if the ReportVehicleSignalsEventsSummary row exists.
func ReportVehicleSignalsEventsSummaryExists(ctx context.Context, exec boil.ContextExecutor, dateID string, integrationID string, powerTrainType string, deviceDefinitionID string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"device_data_api\".\"report_vehicle_signals_events_summary\" where \"date_id\"=$1 AND \"integration_id\"=$2 AND \"power_train_type\"=$3 AND \"device_definition_id\"=$4 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, dateID, integrationID, powerTrainType, deviceDefinitionID)
	}
	row := exec.QueryRowContext(ctx, sql, dateID, integrationID, powerTrainType, deviceDefinitionID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if report_vehicle_signals_events_summary exists")
	}

	return exists, nil
}
