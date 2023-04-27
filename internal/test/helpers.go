package test

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	pb_devices "github.com/DIMO-Network/devices-api/pkg/grpc"
	pb "github.com/DIMO-Network/shared/api/users"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func BuildRequest(method, url, body string) *http.Request {
	req, _ := http.NewRequest(
		method,
		url,
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	return req
}

func Logger() *zerolog.Logger {
	l := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("app", "devices-api").
		Logger()
	return &l
}

// AuthInjectorTestHandler injects fake jwt with sub
func AuthInjectorTestHandler(userID string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": userID,
			"nbf": time.Now().Unix(),
		})

		c.Locals("user", token)
		return c.Next()
	}
}

// SetupAppFiber sets up app fiber with defaults for testing, like our production error handler.
func SetupAppFiber(logger zerolog.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return ErrorHandler(c, err, &logger, false)
		},
	})
	return app
}

// ErrorHandler custom handler to log recovered errors using our logger and return json instead of string
func ErrorHandler(c *fiber.Ctx, err error, logger *zerolog.Logger, isProduction bool) error {
	logger = getLogger(c, logger)

	code := fiber.StatusInternalServerError // Default 500 statuscode

	e, fiberTypeErr := err.(*fiber.Error)
	if fiberTypeErr {
		// Override status code if fiber.Error type
		code = e.Code
	}
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	codeStr := strconv.Itoa(code)

	logger.Err(err).Str("httpStatusCode", codeStr).
		Str("httpMethod", c.Method()).
		Str("httpPath", c.Path()).
		Msg("caught an error from http request")
	// return an opaque error if we're in a higher level environment and we haven't specified an fiber type err.
	if !fiberTypeErr && isProduction {
		err = fiber.NewError(fiber.StatusInternalServerError, "Internal error")
	}

	return c.Status(code).JSON(fiber.Map{
		"code":    code,
		"message": err.Error(),
	})
}

func getLogger(c *fiber.Ctx, d *zerolog.Logger) *zerolog.Logger {
	m := c.Locals("logger")
	if m == nil {
		return d
	}

	l, ok := m.(*zerolog.Logger)
	if !ok {
		return d
	}

	return l
}

type UsersClient struct {
	Store map[string]*pb.User
}

func (c *UsersClient) GetUser(_ context.Context, in *pb.GetUserRequest, _ ...grpc.CallOption) (*pb.User, error) {
	u, ok := c.Store[in.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "No user with that id found.")
	}
	return u, nil
}

type DevicesClient struct {
	Store map[string]*pb_devices.UserDevice
}

// GetUserDevice(ctx context.Context, in *devices.GetUserDeviceRequest, opts ...grpc.CallOption) (*devices.UserDevice, error)
func (c *DevicesClient) GetUserDeviceByTokenId(ctx context.Context, in *pb_devices.GetUserDeviceByTokenIdRequest, opts ...grpc.CallOption) (*pb_devices.UserDevice, error) {
	return &pb_devices.UserDevice{}, nil
}
func (c *DevicesClient) ListUserDevicesForUser(ctx context.Context, in *pb_devices.ListUserDevicesForUserRequest, opts ...grpc.CallOption) (*pb_devices.ListUserDevicesForUserResponse, error) {
	return &pb_devices.ListUserDevicesForUserResponse{}, nil
}
func (c *DevicesClient) ApplyHardwareTemplate(ctx context.Context, in *pb_devices.ApplyHardwareTemplateRequest, opts ...grpc.CallOption) (*pb_devices.ApplyHardwareTemplateResponse, error) {
	return &pb_devices.ApplyHardwareTemplateResponse{}, nil
}
func (c *DevicesClient) GetAllUserDeviceValuation(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb_devices.ValuationResponse, error) {
	return &pb_devices.ValuationResponse{}, nil
}
func (c *DevicesClient) GetUserDeviceByAutoPIUnitId(ctx context.Context, in *pb_devices.GetUserDeviceByAutoPIUnitIdRequest, opts ...grpc.CallOption) (*pb_devices.UserDeviceAutoPIUnitResponse, error) {
	return &pb_devices.UserDeviceAutoPIUnitResponse{}, nil
}
func (c *DevicesClient) GetClaimedVehiclesGrowth(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb_devices.ClaimedVehiclesGrowth, error) {
	return &pb_devices.ClaimedVehiclesGrowth{}, nil
}
func (c *DevicesClient) CreateTemplate(ctx context.Context, in *pb_devices.CreateTemplateRequest, opts ...grpc.CallOption) (*pb_devices.CreateTemplateResponse, error) {
	return &pb_devices.CreateTemplateResponse{}, nil
}
func (c *DevicesClient) RegisterUserDeviceFromVIN(ctx context.Context, in *pb_devices.RegisterUserDeviceFromVINRequest, opts ...grpc.CallOption) (*pb_devices.RegisterUserDeviceFromVINResponse, error) {
	return &pb_devices.RegisterUserDeviceFromVINResponse{}, nil
}

func (c *DevicesClient) GetUserDevice(_ context.Context, in *pb_devices.GetUserDeviceRequest, _ ...grpc.CallOption) (*pb_devices.UserDevice, error) {
	u, ok := c.Store[in.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "No user with that id found.")
	}
	return u, nil
}
