package entity

type Organization struct {
	ID       *string        `json:"id"`
	ParentID *string        `json:"parent_id"`
	GroupID  *string        `json:"group_id"`
	Group    *Group         `json:"group"`
	Children []Organization `json:"children"`
}
