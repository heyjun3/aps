package migrate

import (
	"fmt"
	"time"
	"strconv"
)

func ConvKeepaTimeToTime(keepa_time string) (time.Time, error) {
	k, err := strconv.Atoi(keepa_time)
	if err != nil {
		fmt.Printf("keepa time arg is valid %v", keepa_time)
		return time.Time{}, err
	}
	unix_time := (k + 21564000) * 60
	return time.Unix(int64(unix_time), 0), nil
}
