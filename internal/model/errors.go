package model

import "errors"

// Общие ошибки
var (
	// ErrBadCredentials - переданы неправильные параметры
	ErrBadCredentials = errors.New("bad credentials")
	// ErrUnauthorized - пользователь не авторизован
	ErrUnauthorized = errors.New("unauthorized")
)

// Ошибки модели User
var (
	// ErrUserNotFound - пользователь не найден
	ErrUserNotFound = errors.New("not found")
	// ErrLoginTaken - логин уже занят
	ErrLoginTaken = errors.New("login already taken")
	// ErrNotEnoughBalance - недостаточно средств на балансе пользователя
	ErrNotEnoughBalance = errors.New("not enough balance")
)

// Ошибки модели Order
var (
	// ErrOrderIdLuhnInvalid - номер не валиден по алгоритму Луна
	ErrOrderIdLuhnInvalid = errors.New("order id invalid")
	// ErrOrderIdAlreadyExists - номер заказа уже загружен в базу
	ErrOrderIdAlreadyExists = errors.New("order id already exists")
	// ErrOrderBelongsToAnotherUser - номер заказа принадлежит другому пользователю
	ErrOrderBelongsToAnotherUser = errors.New("order belongs to another user")
	// ErrOrderAlreadyUploaded - заказ с таким номером уже был загружен этим пользоваетем
	ErrOrderAlreadyUploaded = errors.New("order already uploaded")
	// ErrOrderListEmpty - запрашиваемых заказов не найдено
	ErrOrderListEmpty = errors.New("order list empty")
)
