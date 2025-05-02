package client

// DummyLogger is a no-op logger implementation for the temporal client.
type DummyLogger struct{}

// Debug is a no-op method for logging debug messages.
func (log *DummyLogger) Debug(_ string, _ ...interface{}) {
}

// Info is a no-op method for logging info messages.
func (log *DummyLogger) Info(_ string, _ ...interface{}) {
}

// Warn is a no-op method for logging warning messages.
func (log *DummyLogger) Warn(_ string, _ ...interface{}) {
}

// Error is a no-op method for logging error messages.
func (log *DummyLogger) Error(_ string, _ ...interface{}) {
}
