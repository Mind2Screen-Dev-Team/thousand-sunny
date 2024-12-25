package xpanic

import (
	"io"
	"time"

	"github.com/DataDog/gostackparse"
)

type Stack struct {
	// ID is the goroutine id (aka `goid`).
	ID int `json:"id"`
	// State is the `atomicstatus` of the goroutine, or if "waiting" the
	// `waitreason`.
	State string `json:"state"`
	// Wait is the approximate duration a goroutine has been waiting or in a
	// syscall as determined by the first gc after the wait started. Aka
	// `waitsince`.
	Wait time.Duration `json:"wait"`
	// LockedToThread is true if the goroutine is locked by a thread, aka
	// `lockedm`.
	LockedToThread bool `json:"lockedToThread"`
	// Stack is the stack trace of the goroutine.
	Stack []*Frame `json:"stack"`
	// FramesElided is true if the stack trace contains a message indicating that
	// additional frames were elided. This happens when the stack depth exceeds
	// 100.
	FramesElided bool `json:"framesElided"`
	// CreatedBy is the frame that created this goroutine, nil for main().
	CreatedBy *Frame `json:"createdBy"`
	// Ancestors are the Goroutines that created this goroutine.
	// See GODEBUG=tracebackancestors=n in https://pkg.go.dev/runtime.
	Ancestor *Stack `json:"ancestor,omitempty"`
}

// Frame is a single call frame on the stack.
type Frame struct {
	// Func is the name of the function, including package name, e.g. "main.main"
	// or "net/http.(*Server).Serve".
	Func string `json:"func"`
	// File is the absolute path of source file e.g.
	// "/go/src/example.org/example/main.go".
	File string `json:"file"`
	// Line is the line number of inside of the source file that was active when
	// the sample was taken.
	Line int `json:"line"`
}

func ParseStack(r io.Reader) []Stack {
	if r == nil {
		return nil
	}
	_stacks, _ := gostackparse.Parse(r)
	return ConvertGoroutines(_stacks)
}

// ConvertGoroutines converts a slice of *gostackparse.Goroutine to a slice of Stack.
func ConvertGoroutines(goroutines []*gostackparse.Goroutine) []Stack {
	stacks := make([]Stack, len(goroutines))
	for i, g := range goroutines {
		stacks[i] = convertGoroutine(g)
	}
	return stacks
}

// convertGoroutine converts a single *gostackparse.Goroutine to a Stack.
func convertGoroutine(g *gostackparse.Goroutine) Stack {
	return Stack{
		ID:             int(g.ID),
		State:          g.State,
		Wait:           g.Wait,
		LockedToThread: g.LockedToThread,
		Stack:          convertFrames(g.Stack),
		FramesElided:   g.FramesElided,
		CreatedBy:      convertFrame(g.CreatedBy),
		Ancestor:       convertAncestor(g.Ancestor),
	}
}

// convertFrames converts a slice of gostackparse.Frame to a slice of Frame.
func convertFrames(frames []*gostackparse.Frame) []*Frame {
	converted := make([]*Frame, len(frames))
	for i, f := range frames {
		converted[i] = &Frame{
			Func: f.Func,
			File: f.File,
			Line: f.Line,
		}
	}
	return converted
}

// convertFrame converts a single gostackparse.Frame to a Frame pointer.
func convertFrame(frame *gostackparse.Frame) *Frame {
	if frame == nil {
		return nil
	}
	return &Frame{
		Func: frame.Func,
		File: frame.File,
		Line: frame.Line,
	}
}

// convertAncestor recursively converts a *gostackparse.Goroutine ancestor to a *Stack.
func convertAncestor(ancestor *gostackparse.Goroutine) *Stack {
	if ancestor == nil {
		return nil
	}
	stack := convertGoroutine(ancestor)
	return &stack
}
