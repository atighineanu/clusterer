package utils

type MU struct {
	Prefix         string
	Incident       string
	ReleaseRequest string
}
type Command struct {
	Deploy    string
	StackName string
	Pool      struct {
		Name string
		Path string
	}
	Network struct {
		Name string
		IP   string
	}
	Node               map[string]string
	SeedVol_Ubuntu     string
	SeedVM_Ubuntu      string
	SeedVol_Leap       string
	SeedVol_SLES15_2   string
	SeedVM_Leap        string
	SeedVM_SLES15_2    string
	Workers            Nodes
	Masters            Nodes
	SumaOrgCredentials SUMAOrgCreds
	SSHKeyLocation     string `json:SSHKeyLocation`
}

type Nodes struct {
	Count  int
	Distro string
}

type SUMAOrgCreds struct {
	Username string
	Password string
}
