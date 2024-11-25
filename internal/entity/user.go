package entity

import "time"

type User struct {
	ID          *string    `json:"id"`
	GroupID     *string    `json:"group_id"`
	EmployeeID  *string    `json:"employee_id"`
	Name        *string    `json:"name"`
	PhoneNumber *string    `json:"phone_number"`
	Pin         *string    `json:"pin"`
	CreatedDate *time.Time `json:"created_date"`
}

type Group struct {
	ID       *string `json:"id"`
	GroupKey *string `json:"group_key"`
	Name     *string `json:"name"`
	Level    *int    `json:"level"`
	Users    []User  `json:"users"`
}
