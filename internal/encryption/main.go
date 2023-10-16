package encryption

import (
	"github.com/arpanrec/secureserver/internal/serverconfig"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/ProtonMail/gopenpgp/v2/helper"
)

var encryptionConfig serverconfig.EncryptionConfig

var mu = &sync.Mutex{}
var mo = &sync.Once{}

func setGPGInfo() serverconfig.EncryptionConfig {
	mu.Lock()
	mo.Do(func() {
		encryptionConfig = serverconfig.GetConfig().Encryption
		gpgPrivateKey, err := os.ReadFile(encryptionConfig.GPGPrivateKeyFile)
		if err != nil {
			log.Fatalln("Error reading private key: ", err)
		}
		encryptionConfig.GPGPrivateKey = string(gpgPrivateKey)

		gpgPublicKey, err1 := os.ReadFile(encryptionConfig.GPGPublicKeyFile)
		if err1 != nil {
			log.Fatalln("Error reading public key: ", err1)
		}
		encryptionConfig.GPGPublicKey = string(gpgPublicKey)

		gpgPassphrase, err2 := os.ReadFile(encryptionConfig.GPGPassphraseFile)
		if err2 != nil {
			log.Fatalln("Error reading passphrase: ", err2)
		}
		gpgPassphraseSanitized := strings.Split(string(gpgPassphrase), "\n")[0]
		log.Printf("Passphrase: %s", gpgPassphraseSanitized)
		encryptionConfig.GPGPrivateKeyPassphrase = []byte(gpgPassphraseSanitized)

		if encryptionConfig.GPGDeleteKeys {
			log.Println("Deleting keys")
			err3 := os.Remove(encryptionConfig.GPGPrivateKeyFile)
			if err3 != nil {
				log.Fatalln("Error deleting private key: ", err3)
			}
			err4 := os.Remove(encryptionConfig.GPGPublicKeyFile)
			if err4 != nil {
				log.Fatalln("Error deleting public key: ", err4)
			}
			err5 := os.Remove(encryptionConfig.GPGPassphraseFile)
			if err5 != nil {
				log.Fatalln("Error deleting passphrase: ", err5)
			}
		}
	})

	mu.Unlock()
	return encryptionConfig
}

func EncryptMessage(message *string) error {
	setGPGInfo()
	armor, err := helper.EncryptMessageArmored(encryptionConfig.GPGPublicKey, *message)
	if err != nil {
		log.Println("Error encrypting message: ", err)
	}
	*message = armor
	return err
}

func DecryptMessage(armor *string) error {
	setGPGInfo()
	decrypted, err := helper.DecryptMessageArmored(encryptionConfig.GPGPrivateKey, encryptionConfig.GPGPrivateKeyPassphrase, *armor)
	if err != nil {
		log.Println("Error decrypting message: ", err)
	}
	*armor = decrypted
	return err
}
