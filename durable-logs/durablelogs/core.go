package durablelogs

type DlServer struct {
	directory string
}

func NewDLServer(directory string) *DlServer {
	return &DlServer{
		directory: directory,
	}
}
