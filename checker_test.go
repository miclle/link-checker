package checker

import (
	"testing"

	"log"

	"github.com/stretchr/testify/assert"
)

func TestChecker(t *testing.T) {
	assert := assert.New(t)

	err := Check("http://developer.qiniu.com", 1)

	if err != nil {
		log.Fatal("err:", err.Error())
	}

	assert.Nil(err)

}
