package jamf

type BaseType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type User struct {
	BaseType
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	EmailAddress string `json:"email_address"`
	Username     string `json:"username"`
	Sites        []struct {
		Site BaseType `json:"site"`
	} `json:"sites"`
}

type BaseAccount struct {
	Users  []User  `json:"users"`
	Groups []Group `json:"groups"`
}

type UserAccount struct {
	BaseType
	FullName     string   `json:"full_name"`
	Email        string   `json:"email"`
	EmailAddress string   `json:"email_address"`
	Enabled      string   `json:"enabled"`
	AccessLevel  string   `json:"access_level"`
	PrivilegeSet string   `json:"privilege_set"`
	Site         BaseType `json:"site"`
}

type Site struct {
	BaseType
}

type UserGroup struct {
	BaseType
	IsSmart bool   `json:"is_smart"`
	Site    Site   `json:"site"`
	Users   []User `json:"users"`
}

type TokenDetails struct {
	Account Account `json:"account"`
	Sites   []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"sites"`
	AuthenticationType string `json:"authenticationType"`
}

type Account struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	RealName string `json:"realName"`
	Email    string `json:"email"`
}

type Group struct {
	BaseType
	AccessLevel  string   `json:"access_level"`
	PrivilegeSet string   `json:"privilege_set"`
	Site         BaseType `json:"site"`
	Members      []struct {
		User BaseType `json:"user"`
	} `json:"members"`
}
