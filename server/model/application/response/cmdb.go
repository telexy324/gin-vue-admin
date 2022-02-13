package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
)

type ApplicationServerResponse struct {
	Servers []application.ApplicationServer `json:"servers"`
}

type SystemRelationsResponse struct {
	Path RelationPath `json:"path"`
}

type RelationPath struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

type Node struct {
	Id    int    `json:"id"`
	Type  int    `json:"type"` //0 outer 1 inner
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type Link struct {
	VectorType     int      `json:"vector_type"` //0 outer 1 inner
	VectorStrValue string   `json:"vector_str_value"`
	Property       Property `json:"property"`
	StartNodeId    int      `json:"start_node_id"`
	EndNodeId      int      `json:"end_node_id"`
}

type Property struct {
	Relation         string `json:"relation"`
	Url              string `json:"url"`
	ServerUpdateDate string `json:"server_update_date"`
}

//message GetEnterpriseRelationChartRsp {
//x.common.def.RelationPath path = 1;
//bool convert_h5 = 999;
//}
//
//message RelationPath {
//repeated Node nodes = 1;
//repeated Link links = 2;
//}
//
//message Node {
//// company_id、people_id
//int64 id = 1;
//x.common.graph.consts.GraphNodeType type = 2;
//// company_name 、 people_name
//string name = 3;
//// value: number of relational nodes
//int64 value = 4;
//// id to_string
//string s_id = 5;
//}
//
//message Link {
////GraphVectorType
//x.common.graph.consts.GraphVectorType vector_type = 1;
////type string value
//string vector_str_value = 2;
//// property value
//Property property = 3;
//// start node id
//int64 start_node_id = 4;
//// end node id
//int64 end_node_id = 5;
//// string(start node id)
//string s_start_node_id = 6;
//// string(end node id)
//string s_end_node_id = 7;
//}
//
//message Property {
//// 职位
//string position = 1;
//// 持股比例
//double shareholding_ratio = 2;
//// 持股类型名称（investment、shareholder）
//string shareholding_name = 3;
//// 股权更新时间
//string shareholding_update_date = 4;
//}
//
//enum GraphNodeType {
//gn_unknown = 0;
//gn_company = 1;
//gn_people = 2;
//}
//
//enum GraphVectorType {
//gv_unknown = 0;
//gv_shareholder = 1;
//gv_employee = 2;
//gv_subsidiary = 3;
//}
