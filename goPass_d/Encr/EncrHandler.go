package encr

type Encpr interface {
	Encrypt([]byte) ([]byte, []byte, error)
	//			Encrypted Data, Nonce, error
	Decrypt([]byte, []byte) ([]byte, error)
	GenNonce() ([]byte, error)
	Init([]byte, []byte) error
}

type EncFactory interface {
	Create([]byte, []byte) Encpr
}

type AesGCMFactory struct{}

func (a AesGCMFactory) Create(key []byte, salt []byte) Encpr {
	var aesGCM = new(GCM_Encr)
	aesGCM.Init(key, salt)
	return aesGCM
}
