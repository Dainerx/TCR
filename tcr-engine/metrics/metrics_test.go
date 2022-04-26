/*
Copyright (c) 2022 Murex

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

package metrics

import (
	"github.com/murex/tcr/tcr-engine/events"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_compute_score(t *testing.T) {
	var timeInGreenRatio = .5
	var savingRate float64 = 60
	var changesPerCommit float64 = 3
	var expected Score = 10
	assert.Equal(t, expected, computeScore(timeInGreenRatio, savingRate, changesPerCommit))
}

func Test_compute_score_with_0_change_per_commit(t *testing.T) {
	var timeInGreenRatio = .5
	var savingRate float64 = 60
	var changesPerCommit float64
	var expected Score
	assert.Equal(t, expected, computeScore(timeInGreenRatio, savingRate, changesPerCommit))
}

func Test_compute_duration_between_2_records(t *testing.T) {
	startEvent := events.ATcrEvent(events.WithTimestamp(events.TodayAt(4, 39, 31)))
	endEvent := events.ATcrEvent(events.WithTimestamp(events.TodayAt(4, 41, 02)))
	assert.Equal(t, 1*time.Minute+31*time.Second, computeDuration(*startEvent, *endEvent))
}

func Test_compute_duration_between_2_records_with_inverted_timestamp(t *testing.T) {
	startEvent := events.ATcrEvent(events.WithDelay(1 * time.Minute))
	endEvent := events.ATcrEvent()
	assert.Equal(t, 1*time.Minute, computeDuration(*startEvent, *endEvent))
}

func Test_compute_durations_with_no_failing_tests_between_2_records(t *testing.T) {
	startEvent := events.ATcrEvent(events.WithPassingTests())
	endEvent := events.ATcrEvent(events.WithDelay(1 * time.Second))
	assert.Equal(t, 1*time.Second, computeDurationInGreen(*startEvent, *endEvent))
	assert.Equal(t, 0*time.Second, computeDurationInRed(*startEvent, *endEvent))
}

func Test_compute_durations_with_failing_tests_between_2_records(t *testing.T) {
	startEvent := events.ATcrEvent(events.WithFailingTests())
	endEvent := events.ATcrEvent(events.WithDelay(1 * time.Second))
	assert.Equal(t, 0*time.Second, computeDurationInGreen(*startEvent, *endEvent))
	assert.Equal(t, 1*time.Second, computeDurationInRed(*startEvent, *endEvent))
}

func Test_compute_time_ratios_with_no_failing_tests_between_2_records(t *testing.T) {
	startEvent := events.ATcrEvent(events.WithPassingTests())
	endEvent := events.ATcrEvent(events.WithDelay(1 * time.Second))
	assert.Equal(t, float64(1), computeTimeInGreenRatio(*startEvent, *endEvent))
	assert.Equal(t, float64(0), computeTimeInRedRatio(*startEvent, *endEvent))
}

func Test_compute_time_ratios_with_failing_tests_between_2_records(t *testing.T) {
	startEvent := events.ATcrEvent(events.WithFailingTests())
	endEvent := events.ATcrEvent(events.WithDelay(1 * time.Second))
	assert.Equal(t, float64(0), computeTimeInGreenRatio(*startEvent, *endEvent))
	assert.Equal(t, float64(1), computeTimeInRedRatio(*startEvent, *endEvent))
}

func Test_compute_time_ratios_with_no_failing_tests_between_2_records_with_same_timestamp(t *testing.T) {
	startEvent := events.ATcrEvent(events.WithPassingTests())
	endEvent := events.ATcrEvent()
	assert.Equal(t, float64(1), computeTimeInGreenRatio(*startEvent, *endEvent))
	assert.Equal(t, float64(0), computeTimeInRedRatio(*startEvent, *endEvent))
}

func Test_compute_time_ratios_with_failing_tests_between_2_records_with_same_timestamp(t *testing.T) {
	startEvent := events.ATcrEvent(events.WithFailingTests())
	endEvent := events.ATcrEvent()
	assert.Equal(t, float64(0), computeTimeInGreenRatio(*startEvent, *endEvent))
	assert.Equal(t, float64(1), computeTimeInRedRatio(*startEvent, *endEvent))
}

// TODO saving rate per hour
// TODO average size of change per commit