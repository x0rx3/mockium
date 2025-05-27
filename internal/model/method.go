package model

type Method string

const (
	DEFAULTMETHOD Method = GET
	GET           Method = "GET"
	POST          Method = "POST"
	DELETE        Method = "DELETE"
	PATCH         Method = "PATCH"
	PUT           Method = "PUT"
)
