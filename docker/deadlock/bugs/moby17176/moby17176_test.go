// +build linux

package devmapper

import (
	"testing"
	"time"

	"github.com/docker/docker/daemon/graphdriver"
	"github.com/docker/docker/daemon/graphdriver/graphtest"
)

// Make sure devices.Lock() has been release upon return from cleanupDeletedDevices() function
func TestDevmapperLockReleasedDeviceDeletion(t *testing.T) {
	driver := graphtest.GetDriver(t, "devicemapper").(*graphtest.Driver).Driver.(*graphdriver.NaiveDiffDriver).ProtoDriver.(*Driver)
	defer graphtest.PutDriver(t)

	// Call cleanupDeletedDevices() and after the call take and release
	// DeviceSet Lock. If lock has not been released, this will hang.
	driver.DeviceSet.cleanupDeletedDevices()

	doneChan := make(chan bool)

	go func() {
		driver.DeviceSet.Lock()
		defer driver.DeviceSet.Unlock()
		doneChan <- true
	}()

	select {
	case <-time.After(time.Second * 5):
		// Timer expired. That means lock was not released upon
		// function return and we are deadlocked. Release lock
		// here so that cleanup could succeed and fail the test.
		driver.DeviceSet.Unlock()
		t.Fatalf("Could not acquire devices lock after call to cleanupDeletedDevices()")
	case <-doneChan:
	}
}
