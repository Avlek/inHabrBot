package impl

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) Run() {
	crawler := NewCrawler()
	errCh := make(chan error)

	go crawler.Run(errCh)

	select {
	case <-errCh:
		return
	}
}
