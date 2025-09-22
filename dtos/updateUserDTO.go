package dtos

type UpdateUserDTO struct {
    // Указатель на строку. Если в JSON не будет "user_name", это поле будет nil.
	UserName *string `json:"user_name"` 
}