# Google Calendar Hot Tub Scheduler

**Date:** 2026-01-17
**Status:** Approved

## Overview

Add scheduling functionality to the pool controller that triggers the hot tub based on Google Calendar events. The scheduler polls a public iCal feed every 5 minutes and turns on the spa circuit when a timed event is active, respecting pool pump/cleaner conflicts.

## Requirements

- Turn on hot tub (circuit 500) when a calendar event is active
- Ignore all-day events; only timed events trigger the hot tub
- Do not turn on hot tub if pool pump (505) or cleaner (501) is running
- If pool turns off during an event, start the hot tub at that point
- Turn off hot tub when the calendar event ends (once, only for events the scheduler started)
- Do not interfere with manual hot tub usage
- Poll calendar every 5 minutes
- Configuration via environment variable

## Calendar Integration

**Source:** Public iCal URL (no authentication required)

**URL:**
```
https://calendar.google.com/calendar/ical/8af9ec772b71063fd95d6470cbe0e93160498e246c21fc90a298f19777e1c208%40group.calendar.google.com/public/basic.ics
```

**Parsing:** Use `github.com/arran4/golang-ical` library to handle iCal format quirks (timezones, recurrence, multi-line values).

## State Machine

For each calendar event, track state by event UID:

```
States: pending | started | ended
```

### State Transitions

| Current State | Condition | Action | New State |
|---------------|-----------|--------|-----------|
| (none) | Event active, pool running | Log waiting | pending |
| (none) | Event active, pool off | Turn on spa | started |
| pending | Event active, pool running | No action | pending |
| pending | Event active, pool off | Turn on spa | started |
| pending | Event ended | No action | ended |
| started | Event active | No action | started |
| started | Event ended | Turn off spa | ended |
| ended | Any | Cleanup after event passes | (removed) |

### Key Invariant

Only turn OFF the hot tub if we previously turned it ON for that specific event (state was `started`). This prevents interfering with manual usage.

## Architecture

### Package Structure

```
internal/scheduler/
  scheduler.go    # Main Scheduler struct and run loop
  calendar.go     # iCal fetching and parsing
  state.go        # Event state tracking (in-memory)
  scheduler_test.go
  calendar_test.go
```

### Scheduler Struct

```go
type Scheduler struct {
    calendarURL   string
    pollInterval  time.Duration
    bridge        *pool.Bridge
    eventStates   map[string]State
    mu            sync.Mutex
    log           *log.Logger
}

type CalendarEvent struct {
    UID       string
    Start     time.Time
    End       time.Time
    AllDay    bool
}

type State string
const (
    StatePending State = "pending"
    StateStarted State = "started"
    StateEnded   State = "ended"
)
```

### Main Loop

```go
func (s *Scheduler) Run(ctx context.Context) {
    ticker := time.NewTicker(s.pollInterval)
    defer ticker.Stop()

    s.poll()  // Run immediately on start

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            s.poll()
        }
    }
}
```

### Conflict Checking

```go
func (s *Scheduler) isPoolRunning() (bool, error) {
    if err := s.bridge.Update(); err != nil {
        return false, fmt.Errorf("failed to update bridge status: %w", err)
    }

    poolOn := s.bridge.GetCircuitState(gateway.CircuitPool) > 0
    cleanerOn := s.bridge.GetCircuitState(gateway.CircuitCleaner) > 0

    return poolOn || cleanerOn, nil
}
```

### Event Processing

```go
func (s *Scheduler) processEvent(event CalendarEvent, now time.Time) {
    s.mu.Lock()
    defer s.mu.Unlock()

    state := s.eventStates[event.UID]

    // Event is currently active
    if now.After(event.Start) && now.Before(event.End) {
        switch state {
        case "", StatePending:
            poolRunning, err := s.isPoolRunning()
            if err != nil {
                s.log.Printf("error checking pool status: %v", err)
                return
            }
            if poolRunning {
                s.log.Printf("event %s waiting, pool is running", event.UID)
                s.eventStates[event.UID] = StatePending
                return
            }
            if err := s.bridge.SetCircuit(gateway.CircuitSpa, 1); err != nil {
                s.log.Printf("failed to turn on spa: %v", err)
                return
            }
            s.log.Printf("started spa for event %s", event.UID)
            s.eventStates[event.UID] = StateStarted

        case StateStarted:
            // Already started, do nothing
        }
        return
    }

    // Event has ended
    if now.After(event.End) && state == StateStarted {
        if err := s.bridge.SetCircuit(gateway.CircuitSpa, 0); err != nil {
            s.log.Printf("failed to turn off spa: %v", err)
            return
        }
        s.log.Printf("stopped spa for event %s", event.UID)
        s.eventStates[event.UID] = StateEnded
    }
}
```

## Integration with main.go

```go
if calendarURL := os.Getenv(scheduler.EnvCalendarURL); calendarURL != "" {
    log.Printf("Starting calendar scheduler with 5-minute poll interval")
    sched := scheduler.New(calendarURL, bridge, 5*time.Minute)
    go sched.Run(ctx)
}
```

The scheduler is optional - if `CALENDAR_URL` isn't set, it doesn't start.

## Configuration

### Environment Variable

```
CALENDAR_URL=https://calendar.google.com/calendar/ical/...
```

### Makefile.local

```makefile
CALENDAR_URL = https://calendar.google.com/calendar/ical/8af9ec772b71063fd95d6470cbe0e93160498e246c21fc90a298f19777e1c208%40group.calendar.google.com/public/basic.ics
```

### Systemd Service

Add to environment:
```
Environment=CALENDAR_URL=https://calendar.google.com/calendar/ical/...
```

## Logging

Use a prefix logger for scheduler output:

```go
log: log.New(os.Stdout, "[scheduler] ", log.LstdFlags)
```

Example output:
```
[scheduler] 2024/01/15 10:00:00 polling calendar, found 2 active events
[scheduler] 2024/01/15 10:00:00 event abc123 waiting, pool is running
2024/01/15 10:00:01 GET /pool 200 12ms
[scheduler] 2024/01/15 10:05:00 started spa for event abc123
```

Filter with: `make logs | grep "\[scheduler\]"`

## Testing

### Unit Tests

```go
// calendar_test.go
func TestParseICal_TimedEvent(t *testing.T) { ... }
func TestParseICal_AllDayEvent_Ignored(t *testing.T) { ... }
func TestParseICal_RecurringEvent(t *testing.T) { ... }

// scheduler_test.go
func TestProcessEvent_PoolRunning_StaysPending(t *testing.T) { ... }
func TestProcessEvent_PoolOff_StartsSpa(t *testing.T) { ... }
func TestProcessEvent_AlreadyStarted_NoAction(t *testing.T) { ... }
func TestProcessEvent_EventEnds_StopsSpa(t *testing.T) { ... }
func TestProcessEvent_NeverStarted_NoStopOnEnd(t *testing.T) { ... }
```

### Mock Interface

```go
type BridgeInterface interface {
    Update() error
    GetCircuitState(circuitID int) int
    SetCircuit(circuitID, state int) error
}
```

### Integration Test

Fetches the real calendar, finds times with/without events, simulates scheduler behavior:

```go
func TestScheduler_RealCalendar(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    calendarURL := os.Getenv("CALENDAR_URL")
    if calendarURL == "" {
        t.Skip("CALENDAR_URL not set")
    }

    events, err := fetchEvents(calendarURL)
    require.NoError(t, err)

    // Find times during and outside events, simulate state transitions
}
```

## Debug CLI Tool

New binary: `cmd/pool-calendar/main.go`

Shows current calendar status and what the scheduler would do:

```go
func main() {
    calendarURL := os.Getenv("CALENDAR_URL")
    events, _ := scheduler.FetchEvents(calendarURL)
    now := time.Now()

    fmt.Printf("Current time: %s\n\n", now.Format(time.RFC1123))

    active := filterActive(events, now)
    if len(active) == 0 {
        fmt.Println("No active events. Hot tub would NOT be triggered.")
    } else {
        fmt.Println("Active events:")
        for _, e := range active {
            fmt.Printf("  - %s to %s\n", e.Start.Format("3:04 PM"), e.End.Format("3:04 PM"))
        }
        fmt.Println("\nHot tub WOULD be triggered (if pool is off).")
    }

    fmt.Println("\nUpcoming events (next 24h):")
    // ... list upcoming events
}
```

### Makefile Additions

```makefile
build-calendar:
	go build -o bin/pool-calendar ./cmd/pool-calendar

calendar-status:
	CALENDAR_URL=$(CALENDAR_URL) go run ./cmd/pool-calendar
```

## Files to Create

| File | Purpose |
|------|---------|
| `internal/scheduler/scheduler.go` | Main loop, state machine |
| `internal/scheduler/calendar.go` | iCal fetch/parse |
| `internal/scheduler/state.go` | Event state types |
| `internal/scheduler/scheduler_test.go` | State machine unit tests |
| `internal/scheduler/calendar_test.go` | Parsing unit tests |
| `internal/scheduler/integration_test.go` | Real calendar integration test |
| `cmd/pool-calendar/main.go` | Debug CLI tool |

## Files to Modify

| File | Change |
|------|--------|
| `cmd/pool-controller/main.go` | Start scheduler if `CALENDAR_URL` is set |
| `Makefile` | Add `build-calendar` and `calendar-status` targets |
| `go.mod` | Add `github.com/arran4/golang-ical` dependency |

## Dependencies

```
github.com/arran4/golang-ical
```
