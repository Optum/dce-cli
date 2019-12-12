package unit

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	util "github.com/Optum/dce-cli/internal/util"
)

func TestDurationUtil_ParseDuration(t *testing.T) {

	t.Run("Parse valid 7 days", func(t *testing.T) {

		now := time.Now()
		durUtil := util.NewDurationUtil()
		expectedDuration := time.Duration(7*24) * time.Hour

		actualDuration, err := durUtil.ParseDuration("7d")
		assert.Nil(t, err)
		assert.Equal(t, expectedDuration.Milliseconds(), actualDuration.Milliseconds())
		assert.Equal(t, now.Add(actualDuration).Unix(), now.AddDate(0, 0, 7).Unix())
	})

	t.Run("Invalid 0 days", func(t *testing.T) {

		durUtil := util.NewDurationUtil()

		_, err := durUtil.ParseDuration("0d")
		assert.NotNil(t, err)
		assert.Equal(t, "invalid zero or negative date: 0", err.Error())
	})

	t.Run("Valid for 8 hours", func(t *testing.T) {

		now := time.Now()
		durUtil := util.NewDurationUtil()
		expectedDuration := time.Duration(8) * time.Hour

		actualDuration, err := durUtil.ParseDuration("8h")
		assert.Nil(t, err)
		assert.Equal(t, expectedDuration.Milliseconds(), actualDuration.Milliseconds())
		assert.Equal(t, now.Add(actualDuration).Unix(), now.Add(time.Hour*8).Unix())
	})

}

func TestDurationUtil_ExpandEpochTime(t *testing.T) {

	durUtil := util.NewDurationUtil()

	t.Run("Expand epoch time", func(t *testing.T) {
		sevenDaysFromNow := time.Now().AddDate(0, 0, 7)
		epochAsString := strconv.FormatInt(sevenDaysFromNow.Unix(), 10)

		actualTime, err := durUtil.ExpandEpochTime(epochAsString)
		assert.Nil(t, err)
		assert.Equal(t, sevenDaysFromNow.Unix(), actualTime)
	})

	t.Run("Expand string", func(t *testing.T) {
		sevenDaysFromNow := time.Now().AddDate(0, 0, 7)
		sevenDaysAsString := "7d"

		actualTime, err := durUtil.ExpandEpochTime(sevenDaysAsString)
		assert.Nil(t, err)
		assert.Equal(t, sevenDaysFromNow.Unix(), actualTime)
	})

	t.Run("Bad input", func(t *testing.T) {
		sevenDaysAsString := "thisisbad"

		_, err := durUtil.ExpandEpochTime(sevenDaysAsString)
		assert.NotNil(t, err)
		assert.Equal(t, "invalid duration format: thisisbad", err.Error())
	})

}
