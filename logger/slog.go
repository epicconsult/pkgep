package logger

type Slog struct{}

func (l *Slog) Info() {}

func (l *Slog) Error() {}

func (l *Slog) Warn() {}

func (l *Slog) Trace() {}
