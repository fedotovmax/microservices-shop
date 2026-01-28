package utils

import "time"

func TimestamppbToTime(secondsPrt *int32, nanosPrt *int32) time.Time {
	if secondsPrt == nil || nanosPrt == nil {
		return time.Time{}
	}

	seconds := *secondsPrt
	nanos := *nanosPrt

	return time.Unix(int64(seconds), int64(nanos))

}
