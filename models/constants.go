/*
Constants and enums

for models
*/
package models

const (
	MAX_TAG_KEY_LENGTH   = 32
	MAX_TAG_VALUE_LENGTH = 200
	MAX_CULPRIT_LENGTH   = 200
)

/*
MemberType
is used in TeamMemeber and will be used in the future in Organisation
*/
type MemberType int

func (m MemberType) String() (result string) { return MEMBER_TYPE_MAPPING[m] }
func (m MemberType) IsValid(choices ...MemberType) (result bool) {
	if len(choices) == 0 {
		choices = MEMBER_TYPE_LIST
	}
	for _, v := range choices {
		if v == m {
			return true
		}
	}
	return
}

const (
	MEMBER_TYPE_ADMIN = MemberType(iota + 1)
	MEMBER_TYPE_MEMBER
)

var (
	MEMBER_TYPE_LIST    = []MemberType{MEMBER_TYPE_ADMIN, MEMBER_TYPE_MEMBER}
	MEMBER_TYPE_MAPPING = map[MemberType]string{
		MEMBER_TYPE_MEMBER: "member",
		MEMBER_TYPE_ADMIN:  "admin",
	}
)

/*
team status field
*/
type TeamStatus int

func (t TeamStatus) String() string { return TEAM_STATUS_MAPPING[t] }
func (t TeamStatus) IsValid(choices ...TeamStatus) (result bool) {
	if len(choices) == 0 {
		choices = TEAM_STATUS_LIST
	}
	for _, i := range choices {
		if i == t {
			return true
		}
	}
	return
}

const (
	TEAM_STATUS_VISIBLE = TeamStatus(iota + 1)
	TEAM_STATUS_PENDING_DELETION
	TEAM_STATUS_DELETION_IN_PROGRESS
)

var (
	TEAM_STATUS_LIST = []TeamStatus{
		TEAM_STATUS_VISIBLE,
		TEAM_STATUS_PENDING_DELETION,
		TEAM_STATUS_DELETION_IN_PROGRESS,
	}
	TEAM_STATUS_MAPPING = map[TeamStatus]string{
		TEAM_STATUS_VISIBLE:              "visible",
		TEAM_STATUS_PENDING_DELETION:     "pending_deletion",
		TEAM_STATUS_DELETION_IN_PROGRESS: "deletion in progress",
	}
)

/*
EventGroup model ana manager

*/
// status for eventgroup
type EventGroupStatus int

func (e EventGroupStatus) String() string { return EVENT_GROUP_STATUS_MAPPING[e] }
func (e EventGroupStatus) IsValid(choices ...EventGroupStatus) bool {
	if len(choices) == 0 {
		choices = EVENT_GROUP_STATUS_LIST
	}
	for _, v := range choices {
		if e == v {
			return true
		}
	}
	return false
}

const (
	EVENT_GROUP_STATUS_UNRESOLVED = EventGroupStatus(iota + 1)
	EVENT_GROUP_STATUS_RESOLVED
	EVENT_GROUP_STATUS_MUTED
)

var (
	// List of all statuses
	EVENT_GROUP_STATUS_LIST = []EventGroupStatus{
		EVENT_GROUP_STATUS_UNRESOLVED,
		EVENT_GROUP_STATUS_RESOLVED,
		EVENT_GROUP_STATUS_MUTED,
	}
	EVENT_GROUP_STATUS_MAPPING = map[EventGroupStatus]string{
		EVENT_GROUP_STATUS_UNRESOLVED: "unresolved",
		EVENT_GROUP_STATUS_RESOLVED:   "resolved",
		EVENT_GROUP_STATUS_MUTED:      "muted",
	}
)
