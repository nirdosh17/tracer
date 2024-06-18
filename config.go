package tracer

const (
	DEFAULT_HOPS            = 64
	DEFAULT_TIMEOUT_SECONDS = 5
	DEFAULT_MAX_RETRIES     = 2
)

type TracerConfig struct {
	MaxHops        int
	TimeoutSeconds int
	MaxRetries     int
}

func NewConfig() *TracerConfig {
	return &TracerConfig{
		MaxHops:        DEFAULT_HOPS,
		TimeoutSeconds: DEFAULT_TIMEOUT_SECONDS,
		MaxRetries:     DEFAULT_MAX_RETRIES,
	}
}

func (t *TracerConfig) Hops() int {
	if t.MaxHops == 0 {
		t.MaxHops = DEFAULT_HOPS
	}
	return t.MaxHops
}

func (t *TracerConfig) Timeout() int {
	if t.TimeoutSeconds == 0 {
		t.TimeoutSeconds = DEFAULT_TIMEOUT_SECONDS
	}
	return t.TimeoutSeconds
}

func (t *TracerConfig) Retries() int {
	if t.MaxRetries == 0 {
		t.MaxRetries = DEFAULT_MAX_RETRIES
	}
	return t.MaxRetries
}

func (t *TracerConfig) WithHops(h int) *TracerConfig {
	t.MaxHops = h
	return t
}

func (t *TracerConfig) WithTimeout(to int) *TracerConfig {
	t.TimeoutSeconds = to
	return t
}

func (t *TracerConfig) WithMaxRetries(n int) *TracerConfig {
	t.MaxRetries = n
	return t
}
