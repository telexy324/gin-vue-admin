package source

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/gookit/color"
	"gorm.io/gorm"
)

var AuthoritiesMenus = new(authoritiesMenus)

type authoritiesMenus struct{}

type AuthorityMenus struct {
	AuthorityId string `gorm:"column:sys_authority_authority_id"`
	BaseMenuId  uint   `gorm:"column:sys_base_menu_id"`
}

var authorityMenus = []AuthorityMenus{
	{"888", 1},
	{"888", 2},
	{"888", 3},
	{"888", 4},
	{"888", 5},
	{"888", 6},
	{"888", 7},
	{"888", 8},
	{"888", 9},
	{"888", 10},
	{"888", 11},
	{"888", 12},
	{"888", 13},
	{"888", 14},
	{"888", 15},
	{"888", 16},
	{"888", 17},
	{"888", 18},
	{"888", 19},
	{"888", 20},
	{"888", 22},
	{"888", 23},
	{"888", 24},
	{"888", 25},
	{"888", 26},
	{"888", 27},
	{"888", 34},
	{"888", 35},
	{"888", 36},
	{"888", 37},
	{"888", 38},
	{"888", 39},
	{"888", 40},
	{"888", 42},
	{"888", 43},
	{"888", 44},
	{"888", 45},
	{"888", 46},
	{"8881", 1},
	{"8881", 2},
	{"8881", 8},
	{"9528", 1},
	{"9528", 2},
	{"9528", 3},
	{"9528", 4},
	{"9528", 5},
	{"9528", 6},
	{"9528", 7},
	{"9528", 8},
	{"9528", 9},
	{"9528", 10},
	{"9528", 11},
	{"9528", 12},
	{"9528", 14},
	{"9528", 15},
	{"9528", 16},
	{"9528", 17},
	{"9527", 23},
	{"9527", 26},
	{"9527", 27},
	{"9527", 34},
	{"9527", 35},
	{"9527", 36},
	{"9527", 37},
	{"9527", 38},
	{"9527", 39},
	{"9527", 40},
	{"9527", 42},
	{"9527", 43},
	{"9527", 44},
	{"9527", 45},
	{"9527", 46},
	{"9529", 23},
	{"9529", 26},
	{"9529", 27},
	{"9529", 34},
	{"9529", 35},
	{"9529", 36},
	{"9529", 37},
	{"9529", 38},
	{"9529", 39},
	{"9529", 40},
	{"9529", 42},
	{"9529", 43},
	{"9529", 44},
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@description: sys_authority_menus 表数据初始化
func (a *authoritiesMenus) Init() error {
	return global.GVA_DB.Table("sys_authority_menus").Transaction(func(tx *gorm.DB) error {
		if tx.Where("sys_authority_authority_id IN ('888', '8881', '9528', '9527', '9529')").Find(&[]AuthorityMenus{}).RowsAffected == 84 {
			color.Danger.Println("\n[Mysql] --> sys_authority_menus 表的初始数据已存在!")
			return nil
		}
		if err := tx.Create(&authorityMenus).Error; err != nil { // 遇到错误时回滚事务
			return err
		}
		color.Info.Println("\n[Mysql] --> sys_authority_menus 表初始数据成功!")
		return nil
	})
}
