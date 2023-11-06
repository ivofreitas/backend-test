package domain

type User struct {
	ID       int32  `json:"id"`
	Name     string `json:"name" validate:"required"`
	Age      int32  `json:"age" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Address  string `json:"address" validate:"required"`
}

// NewUserWithoutPass - creates a instance of user without password
func NewUserWithoutPass(user *User) *UserNoPass {
	response := new(UserNoPass)
	response.ID = user.ID
	response.Name = user.Name
	response.Age = user.Age
	response.Email = user.Email
	response.Address = user.Address
	return response
}

type UserNoPass struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	Age     int32  `json:"age"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

type GetByIDRequest struct {
	ID int32 `param:"id" validate:"required"`
}

type UpdateRequest struct {
	ID       int32  `param:"id" validate:"required"`
	Name     string `json:"name"`
	Age      int32  `json:"age"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

type DeleteRequest struct {
	ID int32 `param:"id" validate:"required"`
}
