package constants

const (
	Ok = `Ok`

	InsertSubscriptionIncorrectBodyMessage = `В теле ожидался JSON. Пример: {“service_name”: “Yandex Plus”, “price”: 400, “user_id”: “60601fee-2bf1-4721-ae6f-7636e79a0cba”, “start_date”: “07-2025”}`
	DeleteSubscriptionIncorrectBodyMessage = `В теле ожидался JSON. Пример: {“service_name”: “Yandex Plus”, “user_id”: “60601fee-2bf1-4721-ae6f-7636e79a0cba”}`
	UpdateSubscriptionIncorrectBodyMessage = `В теле ожидался JSON. Пример: {“service_name”: “Yandex Plus”, “price”: 400, “user_id”: “60601fee-2bf1-4721-ae6f-7636e79a0cba”, “start_date”: “07-2025”, “end_date”: “08-2025”}`

	PatchSubscriptionIncorrectBodyMessage = `В теле ожидался JSON. Пример: {“service_name”: “Yandex Plus”, “user_id”: “60601fee-2bf1-4721-ae6f-7636e79a0cba”, “price”: 400, “start_date”: “07-2025”, “end_date”: “08-2025”} — укажите хотя бы одно из полей price, start_date, end_date`

	ListSubscriptionMissingUserIDMessage = `Укажите query-параметр user_id (UUID), например: /api/v1/subscription?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba`

	IncorrectDateFormatMessage = `Некорректный формат даты для start_date или end_date. Ожидаемый формат: MM-YYYY (Например: 07-2025).`

	Internal = `Internal`

	ValidationFailedMessage = `Ошибка валидации запроса. В теле ожидался JSON. Пример: {“service_name”: “Yandex Plus”, “price”: 400, “user_id”: “60601fee-2bf1-4721-ae6f-7636e79a0cba”, “start_date”: “07-2025”}`
)
