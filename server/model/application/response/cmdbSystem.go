package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
)

type ApplicationSystemsResponse struct {
	Systems []application.ApplicationSystem `json:"systems"`
}

type ApplicationSystemResponse struct {
	System application.ApplicationSystem `json:"system"`
	Admins []application.Admin           `json:"admins"`
}

type SystemRelationsResponse struct {
	Path RelationPath `json:"path"`
}

//type RelationPath struct {
//	Nodes []Node `json:"nodes"`
//	Links []Link `json:"links"`
//}
//
//type Node struct {
//	Id    int    `json:"id"`
//	Type  int    `json:"type"` //0 outer 1 inner
//	Name  string `json:"name"`
//	Value int    `json:"value"`
//}
//
//type Link struct {
//	VectorType     int      `json:"vector_type"` //0 outer 1 inner
//	VectorStrValue string   `json:"vector_str_value"`
//	Property       Property `json:"property"`
//	StartNodeId    int      `json:"start_node_id"`
//	EndNodeId      int      `json:"end_node_id"`
//}
//
//type Property struct {
//	Relation         string `json:"relation"`
//	Url              string `json:"url"`
//	ServerUpdateDate string `json:"server_update_date"`
//}
