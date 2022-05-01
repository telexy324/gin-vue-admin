package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
)

type TemplatesResponse struct {
	Templates []ansible.Template `json:"templates"`
}

type TemplateResponse struct {
	Template ansible.Template `json:"template"`
}
