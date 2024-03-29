package owner

import (
	"testing"

	mock_services "github.com/DIMO-Network/device-data-api/internal/services/mocks"
	"github.com/DIMO-Network/device-data-api/internal/test"
	pbdvc "github.com/DIMO-Network/devices-api/pkg/grpc"
	pb "github.com/DIMO-Network/users-api/pkg/grpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestOwnerMiddleware(t *testing.T) {
	userID := "louxUser"
	userAddr := "0x1ABC7154748d1ce5144478cdeB574ae244b939B5"
	otherUserID := "stanleyUser"
	otherAddr := "0x3AC4f4Ae05b75b97bfC71Ea518913007FdCaab70"
	userDeviceID := "2OeRoU9VmbFVpgpPy3BjY2WsMMm"

	logger := test.Logger()
	cont := gomock.NewController(t)
	usersClient := &test.UsersClient{}
	devicesClient := mock_services.NewMockDeviceAPIService(cont)
	middleware := New(usersClient, devicesClient, logger)

	app := test.SetupAppFiber(*logger)
	app.Get("/:userDeviceID", test.AuthInjectorTestHandler(userID), middleware, func(c *fiber.Ctx) error {
		logger := c.Locals("logger").(*zerolog.Logger)
		logger.Info().Msg("Omega croggers.")
		return nil
	})

	request := test.BuildRequest("GET", "/"+userDeviceID, "")

	cases := []struct {
		Name                string
		UserDeviceUserID    string
		DeviceUserID        string
		UserExists          bool
		UserEthereumAddress string
		DeviceOwnerAddress  string
		ExpectedCode        int
	}{
		{
			Name:         "NoDevice",
			ExpectedCode: 404,
		},
		{
			Name:             "UserIDMatch",
			UserExists:       true,
			UserDeviceUserID: userID,
			DeviceUserID:     userID,
			ExpectedCode:     200,
		},
		{
			Name:             "UserIDMismatchNoAccount",
			UserDeviceUserID: otherUserID,
			ExpectedCode:     404,
		},
		{
			Name:             "UserIDMismatchNoEthereumAddress",
			UserDeviceUserID: otherUserID,
			UserExists:       true,
			ExpectedCode:     404,
		},
		{
			Name:                "UserIDMismatchNotMinted",
			UserDeviceUserID:    userID,
			UserExists:          true,
			UserEthereumAddress: userAddr,
			ExpectedCode:        404,
		},
		{
			Name:                "UserIDMismatchEthereumAddressMatch",
			UserDeviceUserID:    otherUserID,
			DeviceUserID:        userID,
			DeviceOwnerAddress:  userAddr,
			UserExists:          true,
			UserEthereumAddress: userAddr,
			ExpectedCode:        200,
		},
		{
			Name:                "UserIDMismatchEthereumAddressMismatch",
			UserDeviceUserID:    otherUserID,
			DeviceOwnerAddress:  otherAddr,
			UserExists:          true,
			UserEthereumAddress: userAddr,
			ExpectedCode:        404,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			usersClient.Store = map[string]*pb.User{}

			if c.UserExists {
				u := &pb.User{Id: userID}
				if c.UserEthereumAddress != "" {
					u.EthereumAddress = &c.UserEthereumAddress
				}
				usersClient.Store[userID] = u
			}

			if c.DeviceUserID == "" {
				devicesClient.EXPECT().GetUserDevice(gomock.Any(), userDeviceID).Return(nil, status.Error(codes.NotFound, "Device not found."))
			} else {
				d := &pbdvc.UserDevice{Id: userDeviceID, UserId: c.DeviceUserID}
				if c.DeviceOwnerAddress != "" {
					d.OwnerAddress = common.Hex2Bytes(c.DeviceOwnerAddress)
				}
				devicesClient.EXPECT().GetUserDevice(gomock.Any(), userDeviceID).Return(d, nil)
			}

			res, err := app.Test(request)
			require.Nil(t, err)
			assert.Equal(t, c.ExpectedCode, res.StatusCode)
		})
	}
}
