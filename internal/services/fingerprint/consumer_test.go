package fingerprint

import (
	"context"
	"encoding/json"
	"fmt"
	mock_services "github.com/DIMO-Network/device-data-api/internal/services/mocks"

	"os"
	"testing"

	"github.com/DIMO-Network/device-data-api/internal/test"

	pb "github.com/DIMO-Network/devices-api/pkg/grpc"
	"github.com/DIMO-Network/shared/db"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/mock/gomock"
)

type ConsumerTestSuite struct {
	suite.Suite
	pdb       db.Store
	container testcontainers.Container
	ctx       context.Context
	mockCtrl  *gomock.Controller
	topic     string
	cons      *Consumer
	deviceSvc *mock_services.MockDeviceAPIService
}

const migrationsDirRelPath = "../../../migrations"

// SetupSuite starts container db
func (s *ConsumerTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.pdb, s.container = test.StartContainerDatabase(context.Background(), s.T(), migrationsDirRelPath)
	s.mockCtrl = gomock.NewController(s.T())
	s.topic = "topic.fingerprint"

	s.deviceSvc = mock_services.NewMockDeviceAPIService(s.mockCtrl)

	gitSha1 := os.Getenv("GIT_SHA1")
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("app", "devices-api").
		Str("git-sha1", gitSha1).
		Logger()

	s.cons = &Consumer{
		logger: &logger,
		dbs:    s.pdb,
	}

}

// TearDownSuite cleanup at end by terminating container
func (s *ConsumerTestSuite) TearDownSuite() {
	fmt.Printf("shutting down postgres at with session: %s \n", s.container.SessionID())
	if err := s.container.Terminate(s.ctx); err != nil {
		s.T().Fatal(err)
	}
	s.mockCtrl.Finish()
}

func TestConsumerTestSuite(t *testing.T) {
	suite.Run(t, new(ConsumerTestSuite))
}

func (s *ConsumerTestSuite) TestVinCredentialerHandler_DeviceFingerprint() {

	deviceID := ksuid.New().String()
	//ownerAddress := null.BytesFrom(common.Hex2Bytes("448cF8Fd88AD914e3585401241BC434FbEA94bbb"))
	vin := "W1N2539531F907299"
	userDeviceID := "userDeviceID1"
	deiceDefID := "deviceDefID"

	userDevice := pb.UserDevice{
		Id:                 deviceID,
		UserId:             userDeviceID,
		DeviceDefinitionId: deiceDefID,
		VinConfirmed:       true,
		Vin:                &vin,
	}

	msg :=
		`{
	"data": {"rpiUptimeSecs":39,"batteryVoltage":13.49,"timestamp":1688136702634,"vin":"W1N2539531F907299","protocol":"7"},
	"id": "2RvhwjUbtoePjmXN7q9qfjLQgwP",
	"signature": "7c31e54ddcffc2a548ccaf10ed64b7e4bdd239bbaa3e5f6dba41d3e4051d930b7fbdf184724c2fb8d3b2ac8ac82662d2ed74e881dd01c09c4b2a9b4e62ede5db1b",
	"source": "aftermarket/device/fingerprint",
	"specversion": "1.0",
	"subject": "0x448cF8Fd88AD914e3585401241BC434FbEA94bbb",
	"type": "zone.dimo.aftermarket.device.fingerprint"
}`

	cases := []struct {
		Name             string
		ReturnsError     bool
		ExpectedResponse string
		UserDeviceTable  pb.UserDevice
	}{
		{
			Name:             "No corresponding aftermarket device for address",
			ReturnsError:     true,
			ExpectedResponse: "sql: no rows in result set",
		},
		{
			Name:            "active credential",
			ReturnsError:    false,
			UserDeviceTable: userDevice,
		},
		{
			Name:            "inactive credential",
			ReturnsError:    false,
			UserDeviceTable: userDevice,
		},
		{
			Name:            "invalid token id",
			ReturnsError:    false,
			UserDeviceTable: userDevice,
		},
	}

	for _, c := range cases {
		s.T().Run(c.Name, func(t *testing.T) {

			s.deviceSvc.EXPECT().GetUserDeviceByEthAddr(gomock.Any(), gomock.Any()).Return(&c.UserDeviceTable, nil)

			var event Event
			_ = json.Unmarshal([]byte(msg), &event)
			err := s.cons.HandleDeviceFingerprint(s.ctx, &event)

			if c.ReturnsError {
				assert.ErrorContains(t, err, c.ExpectedResponse)
			} else {
				require.NoError(t, err)
			}
		})
	}

}
