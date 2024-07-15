package scheduler

import (
	"math"
	"time"
)

func Scheduler() time.Duration {
	timeNowHour := time.Now().Hour()
	timeNowMin := time.Now().Minute()
	timeNowInMin := timeNowHour*60 + timeNowMin
	MinutesToReach := (time.Minute * 1440).Abs().Minutes()
	res := int(MinutesToReach) - timeNowInMin - 1
	// duration is nanoseconds in an int64 , therfore you have to convert the minutes to nanoseconds
	durationToSleep := int64(res) * 60 * int64(math.Pow(10, 9))
	return time.Duration(durationToSleep)

}
func ScheduleTo8Am() time.Duration {
	timeNowHour := time.Now().Hour()
	timeNowMin := time.Now().Minute()
	timeNowInMins := timeNowHour*60 + timeNowMin

	// 6 am
	MinutesToReach := (time.Minute * 360).Abs().Minutes()
	if timeNowInMins > int(MinutesToReach) {
		res := 1440 - timeNowInMins
		resFinal := res + 6*60
		return time.Duration(int64(resFinal) * 60 * int64(math.Pow(10, 9)))
	} else {

		MinutesToReach := (time.Minute * 60 * 6).Abs().Minutes()
		res := int(MinutesToReach) - timeNowInMins
		return time.Duration(int64(res) * 60 * int64(math.Pow(10, 9)))

	}

}
