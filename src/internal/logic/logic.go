package logic

import (
	"bdim/src/internal/logic/conf"
)

// Logic struct
type Logic struct {
	c   *conf.Config
}

// New init
func New(c *conf.Config) (l *Logic) {
	l = &Logic{
		c:            c,
	}
	return l
}

// Close close resources.
func (l *Logic) Close() {

}