package migrate

import (
	"fmt"
	"go.uber.org/atomic"
	"hash/crc32"
	nurl "net/url"
	"strings"
)

func suint(n int) uint {
	if n < 0 {
		panic(fmt.Sprintf("suint(%v) expects input >= 0", n))
	}
	return uint(n)
}

// FilterCustomQuery filters all query values starting with `x-`
func FilterCustomQuery(u *nurl.URL) *nurl.URL {
	ux := *u
	vx := make(nurl.Values)
	for k, v := range ux.Query() {
		if len(k) <= 1 || (len(k) > 1 && k[0:2] != "x-") {
			vx[k] = v
		}
	}
	ux.RawQuery = vx.Encode()
	return &ux
}

const advisoryLockIDSalt uint = 1486364155

func GenerateAdvisoryLockId(databaseName string, additionalNames ...string) (string, error) {
	if len(additionalNames) > 0 {
		databaseName = strings.Join(append(additionalNames, databaseName), "\x00")
	}
	sum := crc32.ChecksumIEEE([]byte(databaseName))
	sum = sum * uint32(advisoryLockIDSalt)
	return fmt.Sprint(sum), nil
}

// CasRestoreOnErr CAS wrapper to automatically restore the lock state on error
func CasRestoreOnErr(lock *atomic.Bool, o, n bool, casErr error, f func() error) error {
	if !lock.CAS(o, n) {
		return casErr
	}
	if err := f(); err != nil {
		// Automatically unlock/lock on error
		lock.Store(o)
		return err
	}
	return nil
}
