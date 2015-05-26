package apitest

import (
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/utils"
)

/*
	Create new EventGroup for given project
*/
func CreateEventGroup(ctx *context.Context, project *models.Project) (eventgroup *models.EventGroup, err error) {
	eventgroup = models.NewEventGroup(func(eg *models.EventGroup) {
		eg.ProjectID = project.ID.ToForeignKey()
		eg.Logger = "root"
		eg.Level = 10
		eg.Message = "message"
		eg.Culprit = ""
		eg.Checksum = utils.RandomString(32)
		eg.Platform = "any"
		eg.Status = models.EVENT_GROUP_STATUS_UNRESOLVED
		eg.TimesSeen = 1
	})
	err = eventgroup.Insert(ctx)
	return
}

/*
Creates random events for given EventGroup
*/
func CreateEvents(ctx *context.Context, eventgroup *models.EventGroup, count int) (result []*models.Event, err error) {

	// prepare space
	result = make([]*models.Event, count)

	for i := 0; i < count; i++ {
		event := models.NewEvent(func(e *models.Event) {
			e.EventGroupID = eventgroup.ID.ToForeignKey()
			e.EventID = utils.RandomString(32)
			e.ProjectID = eventgroup.ProjectID
			e.Message = "message " + utils.RandomString(10)
			e.Platform = eventgroup.Platform
			e.Datetime = utils.NowTruncated()
			e.TimeSpent = 0
		})
		if err = event.Insert(ctx); err != nil {
			return
		}
		result[i] = event
	}

	return
}
