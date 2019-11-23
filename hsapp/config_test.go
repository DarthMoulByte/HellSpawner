package hsapp

import (
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestState_LatestFolder(t *testing.T) {
	assert := testify.New(t)
	state := &State{}

	state.setLatestFolder("/var/lib")
	assert.EqualValues([]string{
		"/var/lib",
	}, state.RecentFolders)


	state.setLatestFolder("/usr/bin")
	assert.EqualValues([]string{
		"/usr/bin",
		"/var/lib",
	}, state.RecentFolders)

	state.setLatestFolder("/usr/local/bin")
	assert.EqualValues([]string{
		"/usr/local/bin",
		"/usr/bin",
		"/var/lib",
	}, state.RecentFolders)

}
