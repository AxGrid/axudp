package axtools

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimeInterval(t *testing.T) {
	tt := time.Date(1980, 02, 05, 16, 41, 15, 31, time.UTC)
	assert.Equal(t, GetStringTimeKey(tt, Minute), "1980-02-05 16:41:00")
	assert.Equal(t, GetStringTimeKey(tt, Hour), "1980-02-05 16:00:00")
	assert.Equal(t, GetStringTimeKey(tt, Day), "1980-02-05 00:00:00")
	assert.Equal(t, GetStringTimeKey(tt, Month), "1980-02-01 00:00:00")
	assert.Equal(t, GetStringTimeKey(tt, Week), "1980-02-03 00:00:00")
	assert.Equal(t, GetStringTimeKey(tt, Year), "1980-01-01 00:00:00")
}
