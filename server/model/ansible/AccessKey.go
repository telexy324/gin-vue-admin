package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"

	"github.com/ansible-semaphore/semaphore/util"
)

type AccessKeyType string

const (
	AccessKeySSH           AccessKeyType = "ssh"
	AccessKeyNone          AccessKeyType = "none"
	AccessKeyLoginPassword AccessKeyType = "login_password"
	AccessKeyPAT           AccessKeyType = "pat"
)

// AccessKey represents a key used to access a machine with ansible from semaphore
type AccessKey struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name" binding:"required"`
	// 'ssh/login_password/none'
	Type AccessKeyType `db:"type" json:"type" binding:"required"`

	ProjectID *int `db:"project_id" json:"project_id"`

	// Secret used internally, do not assign this field.
	// You should use methods SerializeSecret to fill this field.
	Secret *string `db:"secret" json:"-"`

	LoginPassword  LoginPassword `db:"-" json:"login_password"`
	SshKey         SshKey        `db:"-" json:"ssh"`
	PAT            string        `db:"-" json:"pat"`
	OverrideSecret bool          `db:"-" json:"override_secret"`

	InstallationKey int64 `db:"-" json:"-"`
}

type LoginPassword struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SshKey struct {
	Login      string `json:"login"`
	Passphrase string `json:"passphrase"`
	PrivateKey string `json:"private_key"`
}

type AccessKeyRole int

const (
	AccessKeyRoleAnsibleUser = iota
	AccessKeyRoleAnsibleBecomeUser
	AccessKeyRoleAnsiblePasswordVault
	AccessKeyRoleGit
)

func (key *AccessKey) Install(usage AccessKeyRole) error {
	rnd, err := rand.Int(rand.Reader, big.NewInt(1000000000))
	if err != nil {
		return err
	}

	key.InstallationKey = rnd.Int64()

	if key.Type == AccessKeyNone {
		return nil
	}

	path := key.GetPath()

	err = key.DeserializeSecret()

	if err != nil {
		return err
	}

	switch usage {
	case AccessKeyRoleGit:
		switch key.Type {
		case AccessKeySSH:
			if key.SshKey.Passphrase != "" {
				return fmt.Errorf("ssh key with passphrase not supported")
			}
			return ioutil.WriteFile(path, []byte(key.SshKey.PrivateKey+"\n"), 0600)
		}
	case AccessKeyRoleAnsiblePasswordVault:
		switch key.Type {
		case AccessKeyLoginPassword:
			return ioutil.WriteFile(path, []byte(key.LoginPassword.Password), 0600)
		}
	case AccessKeyRoleAnsibleBecomeUser:
		switch key.Type {
		case AccessKeyLoginPassword:
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
	case AccessKeyRoleAnsibleUser:
		switch key.Type {
		case AccessKeySSH:
			if key.SshKey.Passphrase != "" {
				return fmt.Errorf("ssh key with passphrase not supported")
			}
			return ioutil.WriteFile(path, []byte(key.SshKey.PrivateKey+"\n"), 0600)
		case AccessKeyLoginPassword:
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

func (key AccessKey) Destroy() error {
	path := key.GetPath()
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil
	}
	return os.Remove(path)
}

// GetPath returns the location of the access key once written to disk
func (key AccessKey) GetPath() string {
	return util.Config.TmpPath + "/access_key_" + strconv.FormatInt(key.InstallationKey, 10)
}

func (key AccessKey) Validate(validateSecretFields bool) error {
	if key.Name == "" {
		return fmt.Errorf("name can not be empty")
	}

	if !validateSecretFields {
		return nil
	}

	switch key.Type {
	case AccessKeySSH:
		if key.SshKey.PrivateKey == "" {
			return fmt.Errorf("private key can not be empty")
		}
	case AccessKeyLoginPassword:
		if key.LoginPassword.Password == "" {
			return fmt.Errorf("password can not be empty")
		}
	}

	return nil
}

func (key *AccessKey) SerializeSecret() error {
	var plaintext []byte
	var err error

	switch key.Type {
	case AccessKeySSH:
		plaintext, err = json.Marshal(key.SshKey)
		if err != nil {
			return err
		}
	case AccessKeyLoginPassword:
		plaintext, err = json.Marshal(key.LoginPassword)
		if err != nil {
			return err
		}
	case AccessKeyPAT:
		plaintext = []byte(key.PAT)
	case AccessKeyNone:
		key.Secret = nil
		return nil
	default:
		return fmt.Errorf("invalid access token type")
	}

	encryptionString := util.Config.GetAccessKeyEncryption()

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

func (key *AccessKey) unmarshalAppropriateField(secret []byte) (err error) {
	switch key.Type {
	case AccessKeySSH:
		sshKey := SshKey{}
		err = json.Unmarshal(secret, &sshKey)
		if err == nil {
			key.SshKey = sshKey
		}
	case AccessKeyLoginPassword:
		loginPass := LoginPassword{}
		err = json.Unmarshal(secret, &loginPass)
		if err == nil {
			key.LoginPassword = loginPass
		}
	case AccessKeyPAT:
		key.PAT = string(secret)
	}
	return
}

//func (key *AccessKey) ClearSecret() {
//	key.LoginPassword = LoginPassword{}
//	key.SshKey = SshKey{}
//	key.PAT = ""
//}

func (key *AccessKey) DeserializeSecret() error {
	if key.Secret == nil || *key.Secret == "" {
		return nil
	}

	ciphertext := []byte(*key.Secret)

	if ciphertext[len(*key.Secret)-1] == '\n' { // not encrypted private key, used for back compatibility
		if key.Type != AccessKeySSH {
			return fmt.Errorf("invalid access key type")
		}
		key.SshKey = SshKey{
			PrivateKey: *key.Secret,
		}
		return nil
	}

	ciphertext, err := base64.StdEncoding.DecodeString(*key.Secret)
	if err != nil {
		return err
	}

	encryptionString := util.Config.GetAccessKeyEncryption()

	if encryptionString == "" {
		err = key.unmarshalAppropriateField(ciphertext)
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

	return key.unmarshalAppropriateField(ciphertext)
}
