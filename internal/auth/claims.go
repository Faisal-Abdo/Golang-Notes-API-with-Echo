package auth

type Claims struct {
	PreferredUsername string `json:"preferred_username"`

	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}
