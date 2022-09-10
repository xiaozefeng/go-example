package bee

type Middleware func(next HandleFunc) HandleFunc
