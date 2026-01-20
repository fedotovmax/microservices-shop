package keys

type DatastarRequestOptions struct {
	Headers map[string]string `json:"headers"`
}

func NewDatastarRequestOptions(h map[string]string) DatastarRequestOptions {
	return DatastarRequestOptions{
		Headers: h,
	}
}
