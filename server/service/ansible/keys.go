package ansible

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
)

type KeyService struct {
}

var KeyServiceApp = new(KeyService)

func (keyService *KeyService) GetAccessKey(projectID int, accessKeyID int) (key ansible.AccessKey, err error) {
	err = global.GVA_DB.Where("project_id=? and id =?", projectID, accessKeyID).First(&key).Error
	return
}

//func (m *AccessKey) GetAccessKeyRefs(projectID int, keyID int) (refs ObjectReferrers, error) {
//	t,s,i:=&Template{},&Schedule{},&Inventory{}
//	refs.Templates, err = d.getObjectRefsFrom(projectID, objectProps, objectID, db.TemplateProps)
//	if err != nil {
//		return
//	}
//
//	refs.Inventories, err = d.getObjectRefsFrom(projectID, objectProps, objectID, db.InventoryProps)
//	if err != nil {
//		return
//	}
//
//	templates, err := d.getObjectRefsFrom(projectID, objectProps, objectID, db.ScheduleProps)
//	if err != nil {
//		return
//	}
//
//	for _, st := range templates {
//		exists := false
//		for _, tpl := range refs.Templates {
//			if tpl.ID == st.ID {
//				exists = true
//				break
//			}
//		}
//		if exists {
//			continue
//		}
//		refs.Templates = append(refs.Templates, st)
//	}
//
//	return
//}

func (keyService *KeyService) GetAccessKeys(projectID int, sortInverted bool, sortBy string) ([]ansible.AccessKey, error) {
	var keys []ansible.AccessKey
	db := global.GVA_DB.Model(&ansible.AccessKey{})
	order := ""
	if sortInverted {
		order = "desc"
	}
	db = db.Where("project_id=?", projectID).Order(sortBy + " " + order)
	err := db.Find(&keys).Error
	return keys, err
}

func (keyService *KeyService) UpdateAccessKey(key ansible.AccessKey) error {
	var oldKey ansible.AccessKey
	upDateMap := make(map[string]interface{})
	upDateMap["name"] = key.Name
	upDateMap["secret"] = key.Secret
	upDateMap["type"] = key.Type

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ? and project_id = ?", key.ID, key.ProjectID).Find(&oldKey)
		txErr := db.Updates(upDateMap).Error
		if txErr != nil {
			global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}
		return nil
	})
	return err
}

func (keyService *KeyService) CreateAccessKey(key *ansible.AccessKey) (newKey *ansible.AccessKey, err error) {
	err = keyService.SerializeSecret(key)
	if err != nil {
		return
	}
	err = global.GVA_DB.Create(&key).Error
	return key, err
}

func (keyService *KeyService) DeleteAccessKey(projectID int, accessKeyID int) error {
	err := global.GVA_DB.Where("id = ? and project_id = ?", accessKeyID, projectID).First(&ansible.AccessKey{}).Error
	if err != nil {
		return err
	}
	var key ansible.AccessKey
	return global.GVA_DB.Where("id = ? and project_id = ?", accessKeyID, projectID).First(&key).Delete(&key).Error
}

func (keyService *KeyService) Install(key *ansible.AccessKey, usage ansible.AccessKeyRole) error {
	rnd, err := rand.Int(rand.Reader, big.NewInt(1000000000))
	if err != nil {
		return err
	}

	key.InstallationKey = rnd.Int64()

	if key.Type == ansible.AccessKeyNone {
		return nil
	}

	path := keyService.GetPath(key)

	err = keyService.DeserializeSecret(key)

	if err != nil {
		return err
	}

	switch usage {
	case ansible.AccessKeyRoleGit:
		switch key.Type {
		case ansible.AccessKeySSH:
			if key.SshKey.Passphrase != "" {
				return fmt.Errorf("ssh key with passphrase not supported")
			}
			return ioutil.WriteFile(path, []byte(key.SshKey.PrivateKey+"\n"), 0600)
		}
	case ansible.AccessKeyRoleAnsiblePasswordVault:
		switch key.Type {
		case ansible.AccessKeyLoginPassword:
			return ioutil.WriteFile(path, []byte(key.LoginPassword.Password), 0600)
		}
	case ansible.AccessKeyRoleAnsibleBecomeUser:
		switch key.Type {
		case ansible.AccessKeyLoginPassword:
			content := make(map[string]string)
			content["ansible_become_user"] = key.LoginPassword.Login
			content["ansible_become_password"] = key.LoginPassword.Password
			var bytes []byte
			bytes, err = json.Marshal(content)
			if err != nil {
				return err
			}
			return ioutil.WriteFile(path, bytes, 0600)
		default:
			return fmt.Errorf("access key type not supported for ansible user")
		}
	case ansible.AccessKeyRoleAnsibleUser:
		switch key.Type {
		case ansible.AccessKeySSH:
			if key.SshKey.Passphrase != "" {
				return fmt.Errorf("ssh key with passphrase not supported")
			}
			return ioutil.WriteFile(path, []byte(key.SshKey.PrivateKey+"\n"), 0600)
		case ansible.AccessKeyLoginPassword:
			content := make(map[string]string)
			content["ansible_user"] = key.LoginPassword.Login
			content["ansible_password"] = key.LoginPassword.Password
			var bytes []byte
			bytes, err = json.Marshal(content)
			if err != nil {
				return err
			}
			return ioutil.WriteFile(path, bytes, 0600)

		default:
			return fmt.Errorf("access key type not supported for ansible user")
		}
	}

	return nil
}

func (keyService *KeyService) Destroy(key *ansible.AccessKey) error {
	path := keyService.GetPath(key)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil
	}
	return os.Remove(path)
}

// GetPath returns the location of the access key once written to disk
func (keyService *KeyService) GetPath(key *ansible.AccessKey) string {
	return global.GVA_CONFIG.Ansible.TmpPath + "/access_key_" + strconv.FormatInt(key.InstallationKey, 10)
}

func (keyService *KeyService) Validate(key *ansible.AccessKey, validateSecretFields bool) error {
	if key.Name == "" {
		return fmt.Errorf("name can not be empty")
	}

	if !validateSecretFields {
		return nil
	}

	switch key.Type {
	case ansible.AccessKeySSH:
		if key.SshKey.PrivateKey == "" {
			return fmt.Errorf("private key can not be empty")
		}
	case ansible.AccessKeyLoginPassword:
		if key.LoginPassword.Password == "" {
			return fmt.Errorf("password can not be empty")
		}
	}

	return nil
}

func (keyService *KeyService) SerializeSecret(key *ansible.AccessKey) error {
	var plaintext []byte
	var err error

	switch key.Type {
	case ansible.AccessKeySSH:
		plaintext, err = json.Marshal(key.SshKey)
		if err != nil {
			return err
		}
	case ansible.AccessKeyLoginPassword:
		plaintext, err = json.Marshal(key.LoginPassword)
		if err != nil {
			return err
		}
	case ansible.AccessKeyPAT:
		plaintext = []byte(key.PAT)
	case ansible.AccessKeyNone:
		key.Secret = nil
		return nil
	default:
		return fmt.Errorf("invalid access token type")
	}

	encryptionString := global.GVA_CONFIG.Ansible.GetAccessKeyEncryption()

	if encryptionString == "" {
		secret := base64.StdEncoding.EncodeToString(plaintext)
		key.Secret = &secret
		return nil
	}

	encryption, err := base64.StdEncoding.DecodeString(encryptionString)

	if err != nil {
		return err
	}

	c, err := aes.NewCipher(encryption)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	secret := base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, plaintext, nil))
	key.Secret = &secret

	return nil
}

func (keyService *KeyService) unmarshalAppropriateField(key *ansible.AccessKey, secret []byte) (err error) {
	switch key.Type {
	case ansible.AccessKeySSH:
		sshKey := ansible.SshKey{}
		err = json.Unmarshal(secret, &sshKey)
		if err == nil {
			key.SshKey = sshKey
		}
	case ansible.AccessKeyLoginPassword:
		loginPass := ansible.LoginPassword{}
		err = json.Unmarshal(secret, &loginPass)
		if err == nil {
			key.LoginPassword = loginPass
		}
	case ansible.AccessKeyPAT:
		key.PAT = string(secret)
	}
	return
}

//func (key *AccessKey) ClearSecret() {
//	key.LoginPassword = LoginPassword{}
//	key.SshKey = SshKey{}
//	key.PAT = ""
//}

func (keyService *KeyService) DeserializeSecret(key *ansible.AccessKey) error {
	if key.Secret == nil || *key.Secret == "" {
		return nil
	}

	ciphertext := []byte(*key.Secret)

	if ciphertext[len(*key.Secret)-1] == '\n' { // not encrypted private key, used for back compatibility
		if key.Type != ansible.AccessKeySSH {
			return fmt.Errorf("invalid access key type")
		}
		key.SshKey = ansible.SshKey{
			PrivateKey: *key.Secret,
		}
		return nil
	}

	ciphertext, err := base64.StdEncoding.DecodeString(*key.Secret)
	if err != nil {
		return err
	}

	encryptionString := global.GVA_CONFIG.Ansible.GetAccessKeyEncryption()

	if encryptionString == "" {
		err = keyService.unmarshalAppropriateField(key, ciphertext)
		if _, ok := err.(*json.SyntaxError); ok {
			err = fmt.Errorf("cannot decrypt access key, perhaps encryption key was changed")
		}
		return err
	}

	encryption, err := base64.StdEncoding.DecodeString(encryptionString)
	if err != nil {
		return err
	}

	c, err := aes.NewCipher(encryption)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	ciphertext, err = gcm.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		if err.Error() == "cipher: message authentication failed" {
			err = fmt.Errorf("cannot decrypt access key, perhaps encryption key was changed")
		}
		return err
	}

	return keyService.unmarshalAppropriateField(key, ciphertext)
}
