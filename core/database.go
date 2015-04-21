package core

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"
	"github.com/lann/squirrel"
	"github.com/mgutz/ansi"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/utils"
)

const (
	// schema migrations database table
	MIGRATIONS_DB_TABLE = "schema_migrations"
)

var (
	yellow = ansi.ColorFunc("yellow+h:black")
	green  = ansi.ColorFunc("green+h:black")
	red    = ansi.ColorFunc("red+h:black")
)

/* @TODO: add possibility to post process. One idea is a list of functions that will be
called after migration. Very useful for e.g. calling createsuperuser
*/
type Migrationer interface {
	// returns identifier of migration (must be unique for plugin)
	ID() string

	// returns list of sql queries
	Queries() ([]string, error)

	// something
	Dependencies() []string

	// sets dependencies
	SetDependencies([]string)

	// runs PostMigrate
	PostMigrate(context *context.Context) error
}

// Simple migration
type simpleMigration struct {
	id             string
	queries        []string
	dependencies   []string
	postmigrations []func(context *context.Context) error
}

func (s *simpleMigration) ID() string                    { return s.id }
func (s *simpleMigration) Queries() ([]string, error)    { return s.queries, nil }
func (s *simpleMigration) Dependencies() []string        { return s.dependencies }
func (s *simpleMigration) SetDependencies(deps []string) { s.dependencies = deps }
func (s *simpleMigration) PostMigrate(context *context.Context) error {
	for _, pmf := range s.postmigrations {
		if err := pmf(context); err != nil {
			return err
		}
	}
	return nil
}

// Constructor for simple migration
func NewMigration(id string, queries, dependencies []string, postmigrations ...func(context *context.Context) error) Migrationer {
	return &simpleMigration{id: id, queries: queries,
		dependencies: dependencies, postmigrations: postmigrations}
}

/* Schema editor
is responsible for database migrations
*/
type SchemaEditor struct {
	pr      *PluginRegistry
	context *context.Context
}

// Constructor function for schema editor
func NewSchemaEditor(context *context.Context, pr *PluginRegistry) *SchemaEditor {
	return &SchemaEditor{pr: pr, context: context}
}

/* Performs migrations
First step is to check whether schema migrations has been installed in db
Then we create schema dependencies resolver to resolve dependencies and
return migrations in correct order.
Then we iterate over this array and check if migration has already been applied
If not we perform migration and save this migration to schem migration table
*/
func (s *SchemaEditor) Migrate() (err error) {
	glog.V(2).Infof("patrol: running migrate.")

	var count int
	if count, err = s.PendingMigrations(); err != nil {
		return
	}
	if count == 0 {
		return fmt.Errorf("no pending migrations")
	}

	if err = s.prepareSchemaMigrations(); err != nil {
		return
	}

	resolver := NewSchemaDependenciesResolver(s.pr)
	items, err := resolver.Resolve()
	if err != nil {
		return err
	}

	for _, item := range items {
		isApplied, errApplied := s.IsAppliedMigration(item.migration, item.plugin)
		if errApplied != nil {
			return errApplied
		}

		fmt.Printf("migration %s:%s ", yellow(item.plugin.ID()), yellow(item.migration.ID()))

		if isApplied {
			fmt.Println("already applied.")
			continue
		}

		// process migration
		if err := s.ProcessMigration(item.migration, item.plugin); err != nil {
			fmt.Printf(red("... Error %s\n"), err.Error())
			return err
		}
		fmt.Println(green("... OK"))

		fmt.Printf("    post migrate %s:%s ", item.plugin.ID(), item.migration.ID())
		if err := item.migration.PostMigrate(s.context); err != nil {
			fmt.Printf(red("...failed. %s\n"), err.Error())
		}
		fmt.Println(green("... OK"))
	}

	return nil
}

/* Process migration runs migration in transaction and inserts into schema table
information about run migration
*/
func (s *SchemaEditor) ProcessMigration(migration Migrationer, plugin Pluginer) (err error) {

	var tx *sqlx.Tx
	if tx, err = s.context.DB.Beginx(); err != nil {
		return
	}

	var queries []string
	if queries, err = migration.Queries(); err != nil {
		defer tx.Rollback()
		return
	}

	// process queries
	for _, query := range queries {
		glog.V(3).Infof("migration %s:%s running query: %s", plugin.ID(), migration.ID(), query)
		if _, err = tx.Exec(query); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %s:%s failed with: %s query: %s", plugin.ID(), migration.ID(), err, query)
		}
	}

	arrayValues := []interface{}{}
	for _, q := range queries {
		arrayValues = append(arrayValues, q)
	}

	qb := utils.QueryBuilder().
		Insert(MIGRATIONS_DB_TABLE).
		Columns("id", "plugin_id", "queries").
		Values(migration.ID(), plugin.ID(), squirrel.Expr("ARRAY["+utils.DBArrayPlaceholder(len(queries))+"]", arrayValues...))

	var (
		args  []interface{}
		query string
	)

	if query, args, err = qb.ToSql(); err != nil {
		return
	}

	if _, err = tx.Exec(query, args...); err != nil {
		tx.Rollback()
		return
	}

	return tx.Commit()
}

/*Checks if schema migrations table exists, if not it will be created
 */
func (s *SchemaEditor) prepareSchemaMigrations() (err error) {
	qb := utils.QueryBuilder().
		Select("COUNT(*)").
		From("information_schema.tables").
		Where("table_name = ?", MIGRATIONS_DB_TABLE)

	var (
		args  []interface{}
		count int
		query string
	)

	if query, args, err = qb.ToSql(); err != nil {
		return
	}

	if err = s.context.DB.Get(&count, query, args...); err != nil {
		return
	}

	// schema migrations table already exists
	if count > 0 {
		return
	}

	glog.V(2).Info("Create schema migrations table.")
	schema := fmt.Sprintf(`CREATE TABLE %s (
		        id character varying(200) NOT NULL,
		        plugin_id character varying(200) NOT NULL,
		        created timestamp NULL default now(),
		        queries text[],
		        PRIMARY KEY(id));`, MIGRATIONS_DB_TABLE)
	_, err = s.context.DB.Exec(schema)
	return err
}

/* Useful as check if system is migrated, returns whether given migration is
   already migrated
*/
func (s *SchemaEditor) IsAppliedMigration(migration Migrationer, plugin Pluginer) (result bool, err error) {
	qb := utils.QueryBuilder().
		Select("COUNT(*)").
		From(MIGRATIONS_DB_TABLE).
		Where("id=? AND plugin_id=?", migration.ID(), plugin.ID())

	// prepare variables
	var args []interface{}
	var count int
	var query string

	query, args, err = qb.ToSql()
	if err = s.context.DB.QueryRowx(query, args...).Scan(&count); err != nil {
		result = false
	} else {
		result = count > 0
	}
	return
}

// returns count of pending migrations
func (s *SchemaEditor) PendingMigrations() (count int, err error) {
	if err = s.prepareSchemaMigrations(); err != nil {
		return
	}
	err = s.pr.Do(func(plugin Pluginer) error {
		for _, migration := range plugin.Migrations() {
			if isApplied, isError := s.IsAppliedMigration(migration, plugin); isError != nil {
				return isError
			} else if isApplied == false {
				count++
			}
		}
		return nil
	})
	return
}

/* SchemaDependenciesResolver
resolves dependencies fo migrations

@TODO: implement dependencies resolver. Current implementatation just returns
migrations in order since the CORE is developed so our migrations are not
really depending on each other.
*/
type SchemaDependenciesResolver struct {
	pr          *PluginRegistry
	pluginitems []*SchemaDependenciesResolverItem

	Result []*SchemaDependenciesResolverItem
}

// Constructor function for SchemaDependenciesResolver
func NewSchemaDependenciesResolver(pr *PluginRegistry) *SchemaDependenciesResolver {
	return &SchemaDependenciesResolver{
		pr: pr,
	}
}

func (s *SchemaDependenciesResolver) Processed(id string) (result *SchemaDependenciesResolverItem, err error) {
	err = fmt.Errorf("not found")
	for _, item := range s.Result {
		if id == item.ID() {
			result, err = item, nil
			break
		}
	}
	return
}

func (s *SchemaDependenciesResolver) IsProcessed(id string) (result bool) {
	for _, item := range s.Result {
		if id == item.ID() {
			result = true
			break
		}
	}
	return
}

func (s *SchemaDependenciesResolver) IsValidDependency(id string) (result bool) {
	s.pr.Do(func(plugin Pluginer) error {
		for _, migration := range plugin.Migrations() {
			identifier := plugin.ID() + ":" + migration.ID()
			if identifier == id {
				result = true
			}
		}
		return nil
	})
	return
}

func (s *SchemaDependenciesResolver) fillData() (err error) {

	s.Result = make([]*SchemaDependenciesResolverItem, 0)

	s.pr.Do(func(plugin Pluginer) error {
		var prev *SchemaDependenciesResolverItem
		for i, migration := range plugin.Migrations() {
			item := &SchemaDependenciesResolverItem{
				plugin:    plugin,
				migration: migration,
			}

			if prev != nil {
				prev.next = item
			}

			if i == 0 {
				s.pluginitems = append(s.pluginitems, item)
			}

			prev = item
		}
		return nil
	})

	return
}

func (sdr *SchemaDependenciesResolver) Resolve() ([]*SchemaDependenciesResolverItem, error) {

	// fill initial data
	sdr.fillData()

	var changed int
	var maxTimesNonChanged = 3

	timesNonChanged := 0

	for {

		changed = 0

		for i, pitem := range sdr.pluginitems {

			if pitem == nil {
				continue
			}

			actual := pitem
			for {

				remaining := []*SchemaDependenciesResolverItem{}
				for _, dependency := range actual.migration.Dependencies() {

					// test if dependency is valid identifier
					if !sdr.IsValidDependency(dependency) {
						return nil, fmt.Errorf("migration `%s:%s` has invalid dependency on `%s`.", actual.plugin.ID(), actual.migration.ID(), dependency)
					}

					if processed, err := sdr.Processed(dependency); err != nil {
						remaining = append(remaining, processed)
					}
				}

				// no more remaining
				if len(remaining) == 0 {
					sdr.Result = append(sdr.Result, actual)
					changed++
					actual = actual.next

				} else {
					actual = actual
					break
				}

				// to next run
				if actual == nil {
					break
				}
			}

			sdr.pluginitems[i] = actual
		}

		if changed == 0 {
			timesNonChanged++

			if timesNonChanged == maxTimesNonChanged {
				deadlocking := []string{}
				for _, pitem := range sdr.pluginitems {
					if pitem == nil {
						continue
					}
					deadlocking = append(deadlocking, pitem.ID())
				}

				if len(deadlocking) > 0 {
					return nil, fmt.Errorf("This migrations are deadlocking %#v", deadlocking)
				}
				break
			}
		} else {
			timesNonChanged = 0
		}
	}

	return sdr.Result, nil
}

// Item returned by SchemaDependenciesResolver.Resolve
type SchemaDependenciesResolverItem struct {
	migration Migrationer
	plugin    Pluginer
	next      *SchemaDependenciesResolverItem
}

func (s *SchemaDependenciesResolverItem) ID() string {
	return fmt.Sprintf("%s:%s", s.plugin.ID(), s.migration.ID())
}
