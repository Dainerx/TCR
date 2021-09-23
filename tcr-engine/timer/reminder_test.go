package timer

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const testTimeout = 2500 * time.Millisecond
const testTickPeriod = 1000 * time.Millisecond

// Timeout

func Test_default_timeout_is_5_min(t *testing.T) {
	r := New(0, testTickPeriod, func(tickIndex int, timestamp time.Time) {})
	assert.Equal(t, 5*time.Minute, r.timeout)
}

func Test_init_with_non_default_timeout(t *testing.T) {
	r := New(testTimeout, testTickPeriod, func(tickIndex int, timestamp time.Time) {})
	assert.Equal(t, testTimeout, r.timeout)
}

func Test_ticking_stops_after_timeout(t *testing.T) {
	r := New(testTimeout, testTickPeriod, func(tickIndex int, timestamp time.Time) {})
	r.Start()
	time.Sleep(testTimeout * 2)
	assert.Equal(t, 2, r.tickCounter)
	assert.Equal(t, StoppedAfterTimeOut, r.state)
}

// Tick Period

func Test_default_tick_period_is_1_min(t *testing.T) {
	r := New(testTimeout, 0, func(tickIndex int, timestamp time.Time) {})
	assert.Equal(t, 1*time.Minute, r.tickPeriod)
}

func Test_init_with_non_default_tick_period(t *testing.T) {
	r := New(testTimeout, testTickPeriod, func(tickIndex int, timestamp time.Time) {})
	assert.Equal(t, testTickPeriod, r.tickPeriod)
}

// Starting Reminder

func Test_start_reminder(t *testing.T) {
	r := New(testTimeout, testTickPeriod, func(tickIndex int, timestamp time.Time) {})
	time.Sleep(testTimeout)
	assert.Equal(t, 0, r.tickCounter)
	r.Start()
	time.Sleep(testTimeout)
	assert.Equal(t, 2, r.tickCounter)
}

// Stopping Reminder

func Test_stop_reminder_before_1st_tick(t *testing.T) {
	r := New(testTimeout, testTickPeriod, func(tickIndex int, timestamp time.Time) {})
	r.Start()
	time.Sleep(testTickPeriod / 2)
	r.Stop()
	time.Sleep(testTimeout)

	assert.Equal(t, 0, r.tickCounter)
	assert.Equal(t, StoppedAfterInterruption, r.state)
}

func Test_stop_reminder_between_1st_and_2nd_tick(t *testing.T) {
	r := New(testTimeout, testTickPeriod, func(tickIndex int, timestamp time.Time) {})
	r.Start()
	time.Sleep(testTickPeriod + testTickPeriod/2)
	r.Stop()
	time.Sleep(testTimeout)

	assert.Equal(t, 1, r.tickCounter)
	assert.Equal(t, StoppedAfterInterruption, r.state)
}

func Test_stop_reminder_after_timeout(t *testing.T) {
	r := New(testTimeout, testTickPeriod, func(tickIndex int, timestamp time.Time) {})
	r.Start()
	time.Sleep(testTimeout * 2)
	r.Stop()

	assert.Equal(t, 2, r.tickCounter)
	assert.Equal(t, StoppedAfterTimeOut, r.state)
}

// Reminder tick counter

func Test_can_track_number_of_ticks_fired(t *testing.T) {
	r := New(testTimeout, testTickPeriod, func(tickIndex int, timestamp time.Time) {})
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

// Reminder callback function

func Test_callback_function_can_retrieve_its_tick_index(t *testing.T) {
	var index int
	r := New(testTimeout, testTickPeriod, func(tickIndex int, timestamp time.Time) {
		index = tickIndex
	})
	r.Start()
	time.Sleep(testTickPeriod + testTickPeriod/2)
	assert.Equal(t, 0, index)
	time.Sleep(testTickPeriod)
	assert.Equal(t, 1, index)
}

func Test_callback_function_can_retrieve_its_timestamp(t *testing.T) {
	var ts [2]time.Time
	r := New(testTimeout, testTickPeriod, func(tickIndex int, timestamp time.Time) {
		ts[tickIndex] = timestamp
	})
	tsStart := time.Now()
	r.Start()
	time.Sleep(testTimeout)
	tsEnd := time.Now()

	assert.True(t, tsStart.Before(ts[0]))
	assert.True(t, ts[0].Before(ts[1]))
	assert.True(t, ts[1].Before(tsEnd))
}
