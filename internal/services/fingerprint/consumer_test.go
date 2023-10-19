package fingerprint

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	mock_services "github.com/DIMO-Network/device-data-api/internal/services/mocks"
	"github.com/DIMO-Network/device-data-api/models"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/DIMO-Network/device-data-api/internal/test"

	pb "github.com/DIMO-Network/devices-api/pkg/grpc"
	"github.com/DIMO-Network/shared/db"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
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
		logger:           &logger,
		dbs:              s.pdb,
		deviceAPIService: s.deviceSvc,
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

func (s *ConsumerTestSuite) TestConsumer_HandleDeviceFingerprint_insert() {
	// validate that we can extract the VIN from the message, properly validate the ecrecover and set the vin in udd.signals
	ctx := context.Background()
	//ownerAddress := null.BytesFrom(common.Hex2Bytes("448cF8Fd88AD914e3585401241BC434FbEA94bbb"))
	vin := "W1N2539531F907299"

	userDevice := pb.UserDevice{
		Id:                 ksuid.New().String(),
		UserId:             "user123",
		DeviceDefinitionId: ksuid.New().String(),
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
	s.deviceSvc.EXPECT().GetUserDeviceByEthAddr(gomock.Any(), gomock.Any()).Return(&userDevice, nil)

	var event Event
	_ = json.Unmarshal([]byte(msg), &event)
	err := s.cons.HandleDeviceFingerprint(s.ctx, &event)
	require.NoError(s.T(), err)

	datum, err := models.UserDeviceData(models.UserDeviceDatumWhere.UserDeviceID.EQ(userDevice.Id)).One(ctx, s.pdb.DBS().Writer)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), vin, gjson.GetBytes(datum.Signals.JSON, "vin.value").String())
	assert.Equal(s.T(), "2023-06-30T14:51:42Z", gjson.GetBytes(datum.Signals.JSON, "vin.timestamp").String())
}

func (s *ConsumerTestSuite) TestConsumer_HandleDeviceFingerprint_update() {
	// validate that we can extract the VIN from the message, properly validate the ecrecover and set the vin in udd.signals
	ctx := context.Background()
	//ownerAddress := null.BytesFrom(common.Hex2Bytes("448cF8Fd88AD914e3585401241BC434FbEA94bbb"))
	vin := "W1N2539531F907299"

	userDevice := pb.UserDevice{
		Id:                 ksuid.New().String(),
		UserId:             "user123",
		DeviceDefinitionId: ksuid.New().String(),
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
	s.deviceSvc.EXPECT().GetUserDeviceByEthAddr(gomock.Any(), gomock.Any()).Return(&userDevice, nil)
	// insert existing record, validate doesn't modifiy any existing data
	udd := models.UserDeviceDatum{
		UserDeviceID:  userDevice.Id,
		IntegrationID: autoPiIntegrationID,
		Signals:       null.JSONFrom([]byte(`{"odometer": {"value": 1234.5}}`)),
	}
	_ = udd.Insert(ctx, s.pdb.DBS().Writer, boil.Infer())

	var event Event
	_ = json.Unmarshal([]byte(msg), &event)
	err := s.cons.HandleDeviceFingerprint(s.ctx, &event)
	require.NoError(s.T(), err)

	datum, err := models.UserDeviceData(models.UserDeviceDatumWhere.UserDeviceID.EQ(userDevice.Id)).One(ctx, s.pdb.DBS().Writer)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), vin, gjson.GetBytes(datum.Signals.JSON, "vin.value").String())
	assert.Equal(s.T(), "2023-06-30T14:51:42Z", gjson.GetBytes(datum.Signals.JSON, "vin.timestamp").String())
	assert.Equal(s.T(), 1234.5, gjson.GetBytes(datum.Signals.JSON, "odometer.value").Float()) // check existing value wasn't modified
}
