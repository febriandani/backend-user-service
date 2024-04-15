// ./internal/userservices.go

package api

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/febriandani/backend-user-service/internal/db"
	"github.com/febriandani/backend-user-service/internal/infra"
	"github.com/febriandani/backend-user-service/internal/utils"
	userValidate "github.com/febriandani/backend-user-service/internal/validate"
	"github.com/febriandani/backend-user-service/protogen/golang/users"
	"github.com/sirupsen/logrus"
)

// UserService should implement the UsersServer interface generated from grpc.
//
// UnimplementedUsersServer must be embedded to have forwarded compatible implementations.
type UserService struct {
	db     *db.DB
	dbConn *infra.DatabaseList
	log    *logrus.Logger
	conf   *infra.AppService
	users.UnimplementedUsersServer
}

// NewUserService creates a new UserService
func NewUserService(db *db.DB, logger *logrus.Logger, dbList *infra.DatabaseList, conf *infra.AppService) UserService {
	return UserService{
		db:     db,
		log:    logger,
		dbConn: dbList,
		conf:   conf,
	}
}

// RegistrationUser implements the RegistrationUser method of the grpc UsersServer interface to add a new user
func (us *UserService) RegistrationUser(ctx context.Context, req *users.PayloadWithSingleUser) (*users.RegistrationUserResponse, error) {
	log.Printf("Received an add user request ")

	//validate input
	message := userValidate.ValidateUserRegistration(req.User)
	if message != nil {
		return &users.RegistrationUserResponse{
			ResponseMap: message,
		}, errors.New("data not valid")
	}

	//start transaction db
	txUser, err := us.dbConn.Backend.Write.Begin()
	if err != nil {
		us.log.WithField("request: ", "transactionUserDBBegin").WithError(err).Errorf("AddUser | Failed to txUserBegin")
		return &users.RegistrationUserResponse{
			ResponseMap: map[string]string{
				"en": "There is an error in the system, please wait for a while our team will fix it immediately.",
				"id": "Terdapat kesalahan pada sistem, mohon tunggu beberapa saat tim kami akan segera memperbaikinya.",
			},
		}, err
	}

	//check Username and email isexist
	isExist, err := us.db.CheckIsExistUser(ctx, req.User)
	if err != nil {
		us.log.WithField("request: ", req.User).WithError(err).Errorf("AddUser | Failed to check is exist user")
		return &users.RegistrationUserResponse{
			ResponseMap: map[string]string{
				"en": "There is an error in the system, please wait for a while our team will fix it immediately.",
				"id": "Terdapat kesalahan pada sistem, mohon tunggu beberapa saat tim kami akan segera memperbaikinya.",
			},
		}, err
	}

	if isExist {
		us.log.WithField("request: ", req.User).WithError(err).Errorf("AddUser | Failed to create user, username or email already exists")
		return &users.RegistrationUserResponse{
			ResponseMap: map[string]string{
				"en": "Failed to create user, username or email already exists.",
				"id": "Gagal membuat pengguna, nama pengguna atau email sudah ada.",
			},
		}, err
	}

	//compare password and re-password
	if req.User.GetPassword() != req.User.GetRepassword() {
		us.log.WithField("request: ", req.User).WithError(err).Errorf("AddUser | Failed to create user, password not same")
		return &users.RegistrationUserResponse{
			ResponseMap: map[string]string{
				"en": "Password and re-password are not the same.",
				"id": "Kata sandi dan kata sandi ulang tidak sama.",
			},
		}, err
	}

	//generate password
	password, err := utils.GeneratePassword(req.User.Password)
	if err != nil {
		us.log.WithField("request", utils.StructToString(nil)).WithError(err).Errorf("ForgotPassword | fail to generate password")
		return &users.RegistrationUserResponse{
			ResponseMap: map[string]string{
				"en": "There was an error changing the password",
				"id": "Ada kesalahan dalam mengubah kata sandi",
			},
		}, err
	}

	//save into db
	_, err = us.db.SaveUser(ctx, txUser, &users.User{
		Username:  req.User.Username,
		Email:     req.User.Email,
		Password:  password,
		IsActive:  true,
		CreatedBy: "system",
	})
	if err != nil {
		us.log.WithField("request: ", req.User).WithError(err).Errorf("AddUser | Failed to save user")
		return &users.RegistrationUserResponse{
			ResponseMap: map[string]string{
				"en": "There is an error in the system, please wait for a while our team will fix it immediately.",
				"id": "Terdapat kesalahan pada sistem, mohon tunggu beberapa saat tim kami akan segera memperbaikinya.",
			},
		}, err
	}

	//commit transaction db
	err = txUser.Commit()
	if err != nil {
		us.log.WithField("request: ", "transactionUserDBCommit").WithError(err).Errorf("AddUser | Failed to txUserCommit")
		return &users.RegistrationUserResponse{
			ResponseMap: map[string]string{
				"en": "There is an error in the system, please wait for a while our team will fix it immediately.",
				"id": "Terdapat kesalahan pada sistem, mohon tunggu beberapa saat tim kami akan segera memperbaikinya.",
			},
		}, err
	}

	return &users.RegistrationUserResponse{
		ResponseMap: map[string]string{
			"en": "Account successfully created",
			"id": "Akun berhasil dibuat",
		},
	}, nil
}

// LoginV1 implements the LoginV1 method of the grpc UsersServer interface to login
func (us *UserService) LoginV1(ctx context.Context, req *users.PayloadWithSingleUser) (*users.LoginResponse, error) {

	log.Printf("Received an add user login request")

	//validate input
	message := userValidate.ValidateUserLogin(req.User)
	if message != nil {
		return &users.LoginResponse{
			ResponseMap: message,
		}, errors.New("data not valid")
	}

	//check Username and email isexist
	isExist, err := us.db.CheckIsExistUser(ctx, req.User)
	if err != nil {
		us.log.WithField("request: ", req.User).WithError(err).Errorf("LoginUser | Failed to check is exist user")
		return &users.LoginResponse{
			ResponseMap: map[string]string{
				"en": "There is an error in the system, please wait for a while our team will fix it immediately.",
				"id": "Terdapat kesalahan pada sistem, mohon tunggu beberapa saat tim kami akan segera memperbaikinya.",
			},
		}, err
	}

	if !isExist {
		us.log.WithField("request: ", req.User).WithError(err).Errorf("LoginUser | Failed to login, username or email not exists")
		return &users.LoginResponse{
			ResponseMap: map[string]string{
				"en": "Login Failed, username or email not exists.",
				"id": "Gagal login, nama pengguna atau email tidak ada.",
			},
		}, err
	}

	userData, err := us.db.GetUserByEmailOrUsername(req.User.Email)
	if err != nil {
		us.log.WithField("request", utils.StructToString(req.User.Email)).WithError(err).Errorf("LoginUser | Failed to login, error from db")
		return &users.LoginResponse{
			ResponseMap: map[string]string{
				"en": "There is an error in the system, please wait for a while our team will fix it immediately.",
				"id": "Terdapat kesalahan pada sistem, mohon tunggu beberapa saat tim kami akan segera memperbaikinya.",
			},
		}, err
	}

	if !userData.IsActive {
		us.log.WithField("response: ", utils.StructToString(userData)).WithError(err).Errorf("LoginUser | Failed to login, user status not active")
		return &users.LoginResponse{
			ResponseMap: map[string]string{
				"en": "Login Failed, status not active.",
				"id": "Gagal login, status tidak aktif.",
			},
		}, err
	}

	isValid, err := utils.ComparePassword(userData.Password, req.User.Password)
	if err != nil {
		us.log.WithField("request: ", utils.StructToString(req)).WithError(err).Errorf("LoginUser | Failed to login, failed to compare password")
		return &users.LoginResponse{
			ResponseMap: map[string]string{
				"en": "There is an error in the system, please wait for a while our team will fix it immediately.",
				"id": "Terdapat kesalahan pada sistem, mohon tunggu beberapa saat tim kami akan segera memperbaikinya.",
			},
		}, err
	}

	if !isValid {
		us.log.WithField("request", utils.StructToString(req)).WithError(err).Errorf("LoginUser | Failed to login, password is incorrect")
		return &users.LoginResponse{
			ResponseMap: map[string]string{
				"en": "Login Failed, password is incorrect",
				"id": "Login Failed, password salah.",
			},
		}, err
	}

	session, err := utils.GetEncrypt([]byte(us.conf.KeyData.User), utils.StructToString(users.CredentialData{
		Id:       userData.GetUserId(),
		Username: userData.GetUsername(),
		Email:    userData.GetEmail(),
	}))
	if err != nil {
		us.log.WithField("request: ", utils.StructToString(req)).WithError(err).Errorf("LoginUser | Failed to login, failed to encrypt jwt")
		return &users.LoginResponse{
			ResponseMap: map[string]string{
				"en": "There is an error in the system, please wait for a while our team will fix it immediately.",
				"id": "Terdapat kesalahan pada sistem, mohon tunggu beberapa saat tim kami akan segera memperbaikinya.",
			},
		}, err
	}

	generateTime := time.Now().UTC()

	accessToken, renewToken, err := infra.GenerateJWT(session)
	if err != nil {
		us.log.WithField("request: ", utils.StructToString(req)).WithError(err).Errorf("LoginUser | Failed to login, failed to generate jwt token")
		return &users.LoginResponse{
			ResponseMap: map[string]string{
				"en": "There is an error in the system, please wait for a while our team will fix it immediately.",
				"id": "Terdapat kesalahan pada sistem, mohon tunggu beberapa saat tim kami akan segera memperbaikinya.",
			},
		}, err
	}

	return &users.LoginResponse{
		UserId:         userData.GetUserId(),
		Username:       userData.GetUsername(),
		Email:          userData.GetEmail(),
		ProfilePicture: "",
		JwtAccess: &users.JWTAccess{
			AccessToken:        accessToken,
			AccessTokenExpired: generateTime.Add(time.Duration(us.conf.Authorization.JWT.AccessTokenDuration) * time.Minute).Format(time.RFC3339),
			RenewToken:         renewToken,
			RenewTokenExpired:  generateTime.Add(time.Duration(us.conf.Authorization.JWT.RefreshTokenDuration) * time.Minute).Format(time.RFC3339),
		},
		ResponseMap: map[string]string{
			"en": "Login Successfully.",
			"id": "Berhasil Login",
		},
	}, nil
}

// GetUser implements the GetUser method of the grpc UsersServer interface to fetch an user for a given userID
func (us *UserService) GetUser(ctx context.Context, req *users.PayloadWithUserID) (*users.PayloadWithSingleUser, error) {
	log.Printf("Received get user request")

	user, err := us.db.GetUserByID(ctx, req.GetUserId())
	if err != nil {
		us.log.WithField("request: ", req).WithError(err).Errorf("GetUser | Failed to get data user")
		return &users.PayloadWithSingleUser{
			User: user,
			ResponseMap: map[string]string{
				"en": "There is an error in the system, please wait for a while our team will fix it immediately.",
				"id": "Terdapat kesalahan pada sistem, mohon tunggu beberapa saat tim kami akan segera memperbaikinya.",
			},
		}, err
	}
	if user == nil {
		return &users.PayloadWithSingleUser{
			User: user,
			ResponseMap: map[string]string{
				"en": "Data not found.",
				"id": "Data tidak ditemukan.",
			},
		}, err
	}

	return &users.PayloadWithSingleUser{
		User: user,
		ResponseMap: map[string]string{
			"en": "Successfully retrieved data",
			"id": "Berhasil menampilkan data",
		},
	}, nil
}

// UpdateUser implements the UpdateUser method of the grpc usersServer interface to update an user
func (o *UserService) UpdateUser(_ context.Context, req *users.PayloadWithSingleUser) (*users.Empty, error) {
	log.Printf("Received an update user request")

	o.db.UpdateUser(req.GetUser())

	return &users.Empty{}, nil
}

// RemoveUser implements the RemoveUser method of the grpc usersServer interface to remove an user
func (o *UserService) RemoveUser(_ context.Context, req *users.PayloadWithUserID) (*users.Empty, error) {
	log.Printf("Received a remove user request")

	o.db.RemoveUser(req.GetUserId())

	return &users.Empty{}, nil
}
