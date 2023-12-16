package ui

import (
	"time"
)

const refreshInterval = 500 * time.Millisecond

func currentTimeString() string {
	t := time.Now()
	return t.Format("15:04:05")
}

func (u *ui) updateTime() {
	for {
		time.Sleep(refreshInterval)
		u.app.QueueUpdateDraw(func() {
			u.curTime.SetText(currentTimeString())
		})
	}
}
