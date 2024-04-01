package label

type Message struct {
	ID          string
	Labels      map[string]string
	Annotations map[string]string
}

// TODO: aggregation
