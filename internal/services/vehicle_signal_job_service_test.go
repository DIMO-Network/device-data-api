package services

import (
	"github.com/golang/mock/gomock"
	"testing"
)

// User device data is getting a different row for all incoming integrations
func TestGetJobContext(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

}
