package models

import (
	"encoding/json"
	"net/http"

	"github.com/lann/squirrel"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/utils"
)

/*
Base Manager implementation.
Every manager that works with database should compose from this object
*/
type Manager struct {
}

/*
Remarshal method takes source object, marshalls it to json and then unmarshalls
to target
*/
func (m *Manager) Remarshal(source, target interface{}) (err error) {
	var body []byte
	if body, err = json.Marshal(source); err != nil {
		return
	}

	return json.Unmarshal(body, target)
}

// returns paging set with settings
func (m *Manager) NewPaging(params ...*utils.PagingParams) *utils.Paging {
	pp := m.NewPagingParamsFromList(params)
	return utils.NewPaging(
		settings.PAGING_MIN_LIMIT,
		settings.PAGING_MAX_LIMIT,
		settings.PAGING_DEFAULT_LIMIT,
		pp,
	)
}

// returns paging set with settings
func (m *Manager) NewPagingFromRequest(r *http.Request, params ...*utils.PagingParams) *utils.Paging {
	paging := m.NewPaging(params...)
	paging.ReadRequest(r)
	return paging
}

func (m *Manager) NewPagingParams() *utils.PagingParams {
	return utils.NewPagingParams(
		settings.PAGING_DEFAULT_LIMIT_PARAM_NAME,
		settings.PAGING_DEFAULT_PAGE_PARAM_NAME,
	)
}

func (m *Manager) NewPagingParamsFromList(params []*utils.PagingParams) *utils.PagingParams {
	if len(params) > 0 {
		return params[0]
	}
	return m.NewPagingParams()
}

func (m *Manager) NewOrdering(allowed ...string) *utils.Ordering {
	return utils.NewOrdering(settings.ORDERING_DEFAULT_PARAM_NAME, allowed...)
}

/* Various Query filter funcs
 */
func (m *Manager) QueryFilterID(id interface{}) utils.QueryFunc {
	return utils.QueryFilterID(id)
}

func (m *Manager) QueryFilterWhere(pred interface{}, args ...interface{}) utils.QueryFunc {
	return utils.QueryFilterWhere(pred, args...)
}

/* Apply paging to SelectBuilder
 */
func (m *Manager) QueryFilterPaging(paging *utils.Paging) utils.QueryFunc {
	return func(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
		return paging.UpdateBuilder(builder)
	}
}
