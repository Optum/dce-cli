package configs

type masterAccount struct {
	Profile *string
}

type admin struct {
	MasterAccount *masterAccount
}

type auth struct {
	LoginUrl *string
}

// Config contains config
type Config struct {
	Admin *admin
	Auth  *auth
}
