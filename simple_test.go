// Copyright 2022 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimplePushIsNotBlockedWithoutSubscribers(t *testing.T) {
	s := NewSimpleChannel[int](0, 0)

	s.Publish(1)

	// just assert that we reach this line, for the sake of having an assert
	assert.True(t, true)
}

func TestSimpleSubscribersAreNotBlockedAfterClose(t *testing.T) {
	s := NewSimpleChannel[int](0, 0)
	ch, err := s.Subscribe()
	assert.Nil(t, err)

	s.Close()

	<-ch

	// just assert that we reach this line, for the sake of having an assert
	assert.True(t, true)
}

func TestSimpleEachSubscribersRecievesEachItem(t *testing.T) {
	s := NewSimpleChannel[int](0, 0)
	input1 := 1
	input2 := 2

	ch1, err := s.Subscribe()
	assert.Nil(t, err)
	ch2, err := s.Subscribe()
	assert.Nil(t, err)

	s.Publish(input1)

	output1Ch1 := <-ch1
	output1Ch2 := <-ch2

	s.Publish(input2)

	output2Ch1 := <-ch1
	output2Ch2 := <-ch2

	assert.Equal(t, input1, output1Ch1)
	assert.Equal(t, input1, output1Ch2)

	assert.Equal(t, input2, output2Ch1)
	assert.Equal(t, input2, output2Ch2)
}

func TestSimpleEachSubscribersRecievesEachItemGivenBufferedEventChan(t *testing.T) {
	s := NewSimpleChannel[int](0, 2)
	input1 := 1
	input2 := 2

	ch1, err := s.Subscribe()
	assert.Nil(t, err)
	ch2, err := s.Subscribe()
	assert.Nil(t, err)

	// both inputs are added first before read, using the internal chan buffer
	s.Publish(input1)
	s.Publish(input2)

	output1Ch1 := <-ch1
	output1Ch2 := <-ch2

	output2Ch1 := <-ch1
	output2Ch2 := <-ch2

	assert.Equal(t, input1, output1Ch1)
	assert.Equal(t, input1, output1Ch2)

	assert.Equal(t, input2, output2Ch1)
	assert.Equal(t, input2, output2Ch2)
}

func TestSimpleSubscribersDontRecieveItemsAfterUnsubscribing(t *testing.T) {
	s := NewSimpleChannel[int](0, 0)
	ch, err := s.Subscribe()
	assert.Nil(t, err)
	s.Unsubscribe(ch)

	s.Publish(1)

	// tiny delay to try and make sure the internal logic would have had time
	// to do its thing with the pushed item.
	time.Sleep(5 * time.Millisecond)

	// closing the channel will result in reads yielding the default value
	assert.Equal(t, 0, <-ch)
}
