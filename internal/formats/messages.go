package formats

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/lachlanharrisdev/praetor/internal/events"
)

type Level int

const (
	LevelInfo Level = iota
	LevelSuccess
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "info"
	case LevelSuccess:
		return "success"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	default:
		return "unknown"
	}
}

func (l Level) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.String())
}

func (l *Level) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		switch strings.ToLower(s) {
		case "info":
			*l = LevelInfo
		case "success":
			*l = LevelSuccess
		case "warn", "warning":
			*l = LevelWarn
		case "error", "err":
			*l = LevelError
		default:
			return fmt.Errorf("invalid level %q", s)
		}
		return nil
	}

	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		*l = Level(i)
		return nil
	}

	return fmt.Errorf("invalid level %s", string(data))
}

var (
	rendererMu        sync.RWMutex
	messageRenderers  = map[Format]messageRenderer{}
	defaultEmitter    *Emitter
	defaultEmitterMux sync.Once
	defaultEmitterMu  sync.RWMutex
)

type Message struct {
	Timestamp string         `json:"timestamp,omitempty"`
	Level     Level          `json:"level"`
	Text      string         `json:"text,omitempty"`
	Fields    map[string]any `json:"fields,omitempty"`
	Event     *events.Event  `json:"event,omitempty"`
}

type Options struct {
	Format       Format
	Writer       io.Writer
	UseTimestamp bool
}

type messageRenderer func([]Message, Options) (string, error)

// RegisterMessageRenderer registers a renderer for a given format.
// format-specific files must call this in init().
func RegisterMessageRenderer(format Format, renderer messageRenderer) {
	rendererMu.Lock()
	defer rendererMu.Unlock()
	messageRenderers[format] = renderer
}

// RenderMessages renders a batch of messages using the renderer registered for the format.
func RenderMessages(format Format, messages []Message, opts Options) (string, error) {
	rendererMu.RLock()
	renderer, ok := messageRenderers[format]
	rendererMu.RUnlock()
	if !ok {
		return "", fmt.Errorf("renderer not registered for format %s", format.String())
	}
	return renderer(messages, opts)
}

// Emitter serializes access to a renderer and writes rendered output to the configured writer.
type Emitter struct {
	mu   sync.Mutex
	opts Options
	w    io.Writer
}

// NewEmitter constructs an emitter with the provided options
func NewEmitter(opts Options) *Emitter {
	w := opts.Writer
	if w == nil {
		w = os.Stdout
	}
	return &Emitter{opts: opts, w: w}
}

// Emit renders and writes a single message
func (e *Emitter) Emit(m Message) {
	e.mu.Lock()
	defer e.mu.Unlock()

	out, err := RenderMessages(e.opts.Format, []Message{m}, e.opts)
	if err != nil {
		fmt.Fprintf(e.w, "render error: %v\n", err)
		return
	}
	fmt.Fprint(e.w, out)
}
func (e *Emitter) Emitf(level Level, format string, args ...any) {
	e.Emit(Message{Level: level, Text: fmt.Sprintf(format, args...)})
}

// Event emits an event wrapped as a message.
func (e *Emitter) Event(ev *events.Event) {
	e.Emit(Message{Level: LevelInfo, Event: ev})
}

// Default returns the process-wide default emitter, creating it if needed.
func Default() *Emitter {
	defaultEmitterMux.Do(func() {
		defaultEmitterMu.Lock()
		defer defaultEmitterMu.Unlock()
		defaultEmitter = NewEmitter(Options{
			Format:       FormatTerminal,
			Writer:       os.Stdout,
			UseTimestamp: false,
		})
	})

	defaultEmitterMu.RLock()
	defer defaultEmitterMu.RUnlock()
	return defaultEmitter
}

// SetDefault overrides the process-wide default emitter.
func SetDefault(em *Emitter) {
	defaultEmitterMux.Do(func() {})
	defaultEmitterMu.Lock()
	defaultEmitter = em
	defaultEmitterMu.Unlock()
}

func Info(msg string)    { Default().Emit(Message{Level: LevelInfo, Text: msg}) }
func Success(msg string) { Default().Emit(Message{Level: LevelSuccess, Text: msg}) }
func Warn(msg string)    { Default().Emit(Message{Level: LevelWarn, Text: msg}) }
func Error(msg string)   { Default().Emit(Message{Level: LevelError, Text: msg}) }

func Infof(format string, args ...any)    { Default().Emitf(LevelInfo, format, args...) }
func Successf(format string, args ...any) { Default().Emitf(LevelSuccess, format, args...) }
func Warnf(format string, args ...any)    { Default().Emitf(LevelWarn, format, args...) }
func Errorf(format string, args ...any)   { Default().Emitf(LevelError, format, args...) }

func Emit(m Message)             { Default().Emit(m) }
func EmitEvent(ev *events.Event) { Default().Event(ev) }
