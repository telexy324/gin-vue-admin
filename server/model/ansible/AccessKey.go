package ansible

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type AccessKeyType string

const (
	AccessKeySSH           AccessKeyType = "ssh"
	AccessKeyNone          AccessKeyType = "none"
	AccessKeyLoginPassword AccessKeyType = "login_password"
	AccessKeyPAT           AccessKeyType = "pat"
)

// AccessKey represents a key used to access a machine with ansible from semaphore
type AccessKey struct {
	global.GVA_MODEL
	Name string `gorm:"column:name" json:"name"`
	// 'ssh/login_password/none'
	Type AccessKeyType `gorm:"column:type" json:"type"`

	ProjectID *int `gorm:"column:project_id" json:"projectId"`

	// Secret used internally, do not assign this field.
	// You should use methods SerializeSecret to fill this field.
	Secret *string `gorm:"column:secret" json:"-"`

	LoginPassword  LoginPassword `gorm:"-" json:"loginPassword"`
	SshKey         SshKey        `gorm:"-" json:"ssh"`
	PAT            string        `gorm:"-" json:"pat"`
	OverrideSecret bool          `gorm:"-" json:"overrideSecret"`

	InstallationKey int64 `gorm:"-" json:"-"`
}

func (m *AccessKey) TableName() string {
	return "ansible_access_keys"
}

type LoginPassword struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SshKey struct {
	Login      string `json:"login"`
	Passphrase string `json:"passphrase"`
	PrivateKey string `json:"privateKey"`
}

type AccessKeyRole int

type ObjectReferrers struct {
	Templates   []Template  `json:"templates"`
	Inventories []Inventory `json:"inventories"`
}

const (
	AccessKeyRoleAnsibleUser = iota
	AccessKeyRoleAnsibleBecomeUser
	AccessKeyRoleAnsiblePasswordVault
	AccessKeyRoleGit
)
