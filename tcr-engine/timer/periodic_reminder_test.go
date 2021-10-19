/*
Copyright (c) 2021 Murex

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package timer

import (
	"flag"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var testTimeout time.Duration
var testTickPeriod time.Duration

func TestMain(m *testing.M) {
	if !flag.Parsed() {
		flag.Parse()
	}
	// Most tests in this file are designed so that reminders fire twice, and then go in time out.
	// We differentiate CI and local machine to optimize test speed execution when run on local machine
	// while not failing when run on CI (which runs slower)
	if testing.Short() {
		fmt.Println("Running tests with fast timers")
		testTimeout = 100 * time.Millisecond
		testTickPeriod = 40 * time.Millisecond
	} else {
		fmt.Println("Running tests with long timers")
		testTimeout = 1000 * time.Millisecond
		testTickPeriod = 400 * time.Millisecond
	}
	// Run tests
	exitVal := m.Run()
	os.Exit(exitVal)
}

// Timeout

func Test_default_timeout_is_5_min(t *testing.T) {
	r := NewPeriodicReminder(0, testTickPeriod, func(ctx ReminderContext) {})
	assert.Equal(t, 5*time.Minute, r.timeout)
}

func Test_init_with_non_default_timeout(t *testing.T) {
	r := NewPeriodicReminder(testTimeout, testTickPeriod, func(ctx ReminderContext) {})
	assert.Equal(t, testTimeout, r.timeout)
}

func Test_ticking_stops_after_timeout(t *testing.T) {
	r := NewPeriodicReminder(testTimeout, testTickPeriod, func(ctx ReminderContext) {})
	r.Start()
	time.Sleep(testTimeout * 2)
	assert.Equal(t, 2, r.tickCounter)
	assert.Equal(t, StoppedAfterTimeOut, r.state)
}

// Tick Period

func Test_default_tick_period_is_1_min(t *testing.T) {
	r := NewPeriodicReminder(testTimeout, 0, func(ctx ReminderContext) {})
	assert.Equal(t, 1*time.Minute, r.tickPeriod)
}

func Test_init_with_non_default_tick_period(t *testing.T) {
	r := NewPeriodicReminder(testTimeout, testTickPeriod, func(ctx ReminderContext) {})
	assert.Equal(t, testTickPeriod, r.tickPeriod)
}

// Starting PeriodicReminder

func Test_start_reminder(t *testing.T) {
	r := NewPeriodicReminder(testTimeout, testTickPeriod, func(ctx ReminderContext) {})
	time.Sleep(testTimeout)
	assert.Equal(t, 0, r.tickCounter)
	r.Start()
	time.Sleep(testTimeout)
	assert.Equal(t, 2, r.tickCounter)
}

func Test_start_reminder_triggers_start_event(t *testing.T) {
	var eventFired = false
	r := NewPeriodicReminder(testTimeout, testTickPeriod, func(ctx ReminderContext) {
		if ctx.eventType == StartEvent {
			eventFired = true
		}
	})
	r.Start()
	r.Stop()

	assert.True(t, eventFired)
}

// Stopping PeriodicReminder

func Test_stop_reminder_before_1st_tick(t *testing.T) {
	r := NewPeriodicReminder(testTimeout, testTickPeriod, func(ctx ReminderContext) {})
	r.Start()
	time.Sleep(testTickPeriod / 2)
	r.Stop()
	time.Sleep(testTimeout)

	assert.Equal(t, 0, r.tickCounter)
	assert.Equal(t, StoppedAfterInterruption, r.state)
}

func Test_stop_reminder_between_1st_and_2nd_tick(t *testing.T) {
	r := NewPeriodicReminder(testTimeout, testTickPeriod, func(ctx ReminderContext) {})
	r.Start()
	time.Sleep(testTickPeriod + testTickPeriod/2)
	r.Stop()
	time.Sleep(testTimeout)

	assert.Equal(t, 1, r.tickCounter)
	assert.Equal(t, StoppedAfterInterruption, r.state)
}

func Test_stop_reminder_after_timeout(t *testing.T) {
	r := NewPeriodicReminder(testTimeout, testTickPeriod, func(ctx ReminderContext) {})
	r.Start()
	time.Sleep(testTimeout * 2)
	r.Stop()

	assert.Equal(t, 2, r.tickCounter)
	assert.Equal(t, StoppedAfterTimeOut, r.state)
}

// PeriodicReminder tick counter

func Test_can_track_number_of_ticks_fired(t *testing.T) {
	r := NewPeriodicReminder(testTimeout, testTickPeriod, func(ctx ReminderContext) {})
	r.Start()
	assert.Equal(t, 0, r.tickCounter)
	time.Sleep(testTickPeriod / 2)
	assert.Equal(t, 0, r.tickCounter)
	time.Sleep(testTickPeriod)
	assert.Equal(t, 1, r.tickCounter)
	time.Sleep(testTickPeriod)
	assert.Equal(t, 2, r.tickCounter)
	time.Sleep(testTickPeriod)
	assert.Equal(t, 2, r.tickCounter)
	assert.Equal(t, StoppedAfterTimeOut, r.state)
}

// PeriodicReminder callback function

func Test_callback_function_can_know_current_tick_index(t *testing.T) {
	var index int
	r := NewPeriodicReminder(testTimeout, testTickPeriod, func(ctx ReminderContext) {
		if ctx.eventType == PeriodicEvent {
			index = ctx.index
		}
	})
	r.Start()
	time.Sleep(testTickPeriod + testTickPeriod/2)
	assert.Equal(t, 0, index)
	time.Sleep(testTickPeriod)
	assert.Equal(t, 1, index)
}

func Test_callback_function_can_know_timestamp(t *testing.T) {
	var tsPeriodic [2]time.Time
	var tsStart time.Time
	var tsTimeout time.Time
	r := NewPeriodicReminder(testTimeout, testTickPeriod, func(ctx ReminderContext) {
		switch ctx.eventType {
		case StartEvent:
			tsStart = ctx.timestamp
		case PeriodicEvent:
			tsPeriodic[ctx.index] = ctx.timestamp
		case TimeoutEvent:
			tsTimeout = ctx.timestamp
		}
	})
	r.Start()
	time.Sleep(testTimeout * 2)
	tsEnd := time.Now()

	assert.True(t, tsStart.Before(tsPeriodic[0]))
	assert.True(t, tsPeriodic[0].Before(tsPeriodic[1]))
	assert.True(t, tsPeriodic[1].Before(tsTimeout))
	assert.True(t, tsTimeout.Before(tsEnd))
}

func Test_callback_function_can_know_elapsed_time_since_start(t *testing.T) {
	r := NewPeriodicReminder(testTimeout, testTickPeriod, func(ctx ReminderContext) {
		var expected time.Duration
		switch ctx.eventType {
		case StartEvent:
			expected = 0
		case PeriodicEvent:
			expected = testTickPeriod * time.Duration(ctx.index+1)
		case TimeoutEvent:
			expected = testTimeout
		}
		assert.Equal(t, expected, ctx.elapsed)
	})
	r.Start()
	time.Sleep(testTimeout)
}

func Test_callback_function_can_know_remaining_time_until_end(t *testing.T) {
	r := NewPeriodicReminder(testTimeout, testTickPeriod, func(ctx ReminderContext) {
		var expected time.Duration
		switch ctx.eventType {
		case StartEvent:
			expected = testTimeout
		case PeriodicEvent:
			expected = testTimeout - testTickPeriod*time.Duration(ctx.index+1)
		case TimeoutEvent:
			expected = 0
		}
		assert.Equal(t, expected, ctx.remaining)
	})
	r.Start()
	time.Sleep(testTimeout)
}

func Test_callback_function_can_know_max_index_value(t *testing.T) {
	var expected = int(testTimeout/testTickPeriod) - 1
	r := NewPeriodicReminder(testTimeout, testTickPeriod, func(ctx ReminderContext) {
		assert.Equal(t, expected, ctx.indexMax)
	})
	r.Start()
	time.Sleep(testTimeout)
}

// Time elapsed since timer started

func Test_retrieving_time_elapsed_since_timer_started(t *testing.T) {
	r := NewPeriodicReminder(testTimeout, testTickPeriod, func(ctx ReminderContext) {})
	// Before calling Start(), time elapsed should stick to 0
	assert.Zero(t, r.GetElapsedTime())
	r.Start()
	time.Sleep(testTimeout / 2)
	// While timer is running, total time elapsed is time spent since Start()
	assert.InEpsilon(t, testTimeout/2, r.GetElapsedTime(), 0.3)
	time.Sleep(testTimeout)
	// When timer is done, time elapsed should stopTicking increasing
	assert.InEpsilon(t, testTimeout, r.GetElapsedTime(), 0.3)
}

// Time remaining until timer ends

func Test_retrieving_time_remaining_until_timer_ends(t *testing.T) {
	r := NewPeriodicReminder(testTimeout, testTickPeriod, func(ctx ReminderContext) {})
	// Before calling Start(), time remaining should stick to timeout
	assert.Equal(t, testTimeout, r.GetRemainingTime())
	r.Start()
	time.Sleep(testTimeout / 2)
	// While timer is running, total time remaining is timeout - time spent since Start()
	assert.InEpsilon(t, testTimeout/2, r.GetRemainingTime(), 0.3)
	time.Sleep(testTimeout)
	// When timer is done, time remaining should be 0
	assert.Zero(t, r.GetRemainingTime())
}
