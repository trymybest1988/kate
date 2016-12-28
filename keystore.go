package kate

type KeyStore interface {
	GetKey(string) string
	AddKey(string, string)
}

type MemKeyStore map[string]string

func NewMemKeyStore() KeyStore {
	return make(MemKeyStore)
}

func (s MemKeyStore) GetKey(appId string) (appKey string) {
	appKey = s[appId]
	return
}

func (s MemKeyStore) AddKey(appId, appKey string) {
	s[appId] = appKey
}
