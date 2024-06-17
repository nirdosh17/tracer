package tracer

const (
	DEFAULT_HOPS       = 64
	DEFAULT_TIMEOUT_MS = 100
)

type TracerConfig struct {
	MaxHops   int
	TimeoutMs int
}

func NewConfig() *TracerConfig {
	return &TracerConfig{
		MaxHops:   DEFAULT_HOPS,
		TimeoutMs: DEFAULT_TIMEOUT_MS,
	}
}

func (t *TracerConfig) Hops() int {
	if t.MaxHops == 0 {
		t.MaxHops = DEFAULT_HOPS
	}
	return t.MaxHops
}

func (t *TracerConfig) Timeout() int {
	if t.TimeoutMs == 0 {
		t.TimeoutMs = DEFAULT_TIMEOUT_MS
	}
	return t.TimeoutMs
}

func (t *TracerConfig) WithHops(h int) *TracerConfig {
	t.MaxHops = h
	return t
}

func (t *TracerConfig) WithTimeout(to int) *TracerConfig {
	t.TimeoutMs = to
	return t
}
