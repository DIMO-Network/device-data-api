package owner

import (
	"context"

	"github.com/DIMO-Network/device-data-api/internal/controllers"
	"github.com/DIMO-Network/device-data-api/internal/services"
	pb "github.com/DIMO-Network/users-api/pkg/grpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errNotFound = fiber.NewError(fiber.StatusNotFound, "Device not found.")

// New creates a new middleware handler that checks whether a user is authorized to access
// a user device. For the middleware to allow the request to proceed:
//
//   - The request must have a valid JWT, identifying a user.
//   - There must be a userDeviceID path parameter, and that device must exist.
//   - Either the user owns the device, or the user's account has an Ethereum address that
//     owns the corresponding NFT.
func New(usersClient pb.UserServiceClient, devicesClient services.DeviceAPIService, logger *zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := controllers.GetUserID(c)
		udi := c.Params("userDeviceID")
		logger := logger.With().Str("userId", userID).Str("userDeviceId", udi).Logger()

		c.Locals("userID", userID)
		c.Locals("userDeviceID", udi)
		c.Locals("logger", &logger)

		device, err := devicesClient.GetUserDevice(context.Background(), udi)
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
				return errNotFound
			}
			return err
		}

		if device.UserId == userID {
			return c.Next()
		}

		user, err := usersClient.GetUser(c.Context(), &pb.GetUserRequest{Id: userID})
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
				return errNotFound
			}
			return err
		}

		if user.EthereumAddress == nil {
			return errNotFound
		}

		if common.HexToAddress(*user.EthereumAddress) == common.BytesToAddress(device.OwnerAddress) {
			return c.Next()
		}

		return errNotFound
	}
}
