package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"newsfeed/pkg/service/model"
	"newsfeed/pkg/util/validator"
	gen "newsfeed/proto/gen/proto"
)

type UserService interface {
	CreateUser(user *model.User) (*model.User, error)
}

type httpHandler struct {
	userGprcClient gen.UserServiceClient

	userGrpcConn *grpc.ClientConn
}

type HttpHandlerConfig struct {
	UserGrpcAddr string
}

func New(conf HttpHandlerConfig) (*httpHandler, error) {
	// Set up a connection to the server.
	conn, err := grpc.NewClient(conf.UserGrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	userGprcClient := gen.NewUserServiceClient(conn)

	handler := &httpHandler{
		userGprcClient: userGprcClient,
		userGrpcConn:   conn,
	}
	return handler, nil
}

func (h *httpHandler) Start() {
	r := gin.Default()

	r.POST("/signup", h.Signup)

	r.Run(":33333") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func (h *httpHandler) Stop() {
	h.userGrpcConn.Close()
}

// sign up
type (
	SignupRequest struct {
		Username  string `json:"user_name"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		DOB       string `json:"dob"`
		Email     string `json:"email"`
	}

	SignupUser struct {
		UserId    string
		Username  string
		FirstName string
		LastName  string
		DOB       string
		Email     string
	}
	SignupResponse struct {
		ErrCode ErrCode
		Message string
		User    SignupUser
	}
)

func (h *httpHandler) Signup(c *gin.Context) {
	var signupReq = new(SignupRequest) // &SignupRequest{}
	if err := c.ShouldBindJSON(signupReq); err != nil {
		log.Println("[DEBUG] invalid request", signupReq)
		c.JSON(http.StatusBadRequest, SignupResponse{ErrCode: ErrCodeInvalidInput, Message: "Invalid request"})
		return
	}

	dobTimestamp, err := parseStringToTimestamp(signupReq.DOB)
	if err != nil {
		c.JSON(http.StatusBadRequest, SignupResponse{ErrCode: ErrCodeInvalidInput, Message: "Invalid DOB"})
		return
	}

	// call grpc user
	req := &gen.CreateUserRequest{
		Username:  signupReq.Username,
		Password:  signupReq.Password,
		FirstName: signupReq.FirstName,
		LastName:  signupReq.LastName,
		Dob:       int32(dobTimestamp),
		Email:     signupReq.Email,
	}
	resp, err := h.userGprcClient.CreateUser(context.Background(), req)
	if err != nil {
		log.Println("failed to call user grpc CreateUser", err)
		c.JSON(http.StatusInternalServerError, SignupResponse{ErrCode: ErrCodeInternal, Message: "Something wrong in server"})
		return
	}
	createdUser := resp.GetUser()

	c.JSON(http.StatusOK, SignupResponse{ErrCode: 0, Message: "Sign up successfully",
		User: SignupUser{
			UserId:    strconv.Itoa(int(createdUser.Id)),
			Username:  createdUser.Username,
			FirstName: createdUser.FirstName,
			LastName:  createdUser.LastName,
			DOB:       parseTimestampToString(int(createdUser.Dob)),
			Email:     createdUser.Email,
		}})
}

func validateSignupReq(req *SignupRequest) error {
	// username
	if err := validator.ValidateNotEmpty(req.Username, "user_name"); err != nil {
		return err
	}
	if err := validator.ValidateLength(req.Username, 6, 20, "user_name"); err != nil {
		return err
	}
	// TODO: check ValidateAsciiCharacters
	if err := validator.ValidateAsciiCharacters(req.Username, "user_name"); err != nil {
		return err
	}

	// TODO: validate other fields

	return nil
}

// util
func parseStringToTimestamp(dateStr string) (int, error) {
	// DD-MM-YYYY
	var date time.Time
	date, err := time.Parse("02-01-2006", dateStr)
	if err != nil {
		return 0, fmt.Errorf("error parse time from str=%s", dateStr)
	}
	return int(date.Unix()), nil
}

func parseTimestampToString(timestamp int) string {
	date := time.Unix(int64(timestamp), 0)
	return date.Format("02-01-2006")
}
