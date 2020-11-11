package data

type Command struct {
	Deploy    string
	StackName string
	Pool      struct {
		Name string
		Path string
	}
	Node             map[string]string
	SeedVol_Leap     string
	SeedVol_SLES15_2 string
	SeedVM_Leap      string
	SeedVM_SLES15_2  string
	Workers          Nodes
	Masters          Nodes
}

type Nodes struct {
	Count  int
	Distro string
}
