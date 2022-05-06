package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/gonzalo91/go-serverless/pkg/user"
)

var ErrorMethodNotAllowed = "method not allowed"

type ErrorBody struct{
	ErrorMsg *string `json:"error,omitempty"`
}

func Show(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(
	*events.APIGatewayProxyResponse, error,
){

	email := req.QueryStringParameters["email"]
	if len(email) > 0 {
		result, err := user.FetchUser(email, tableName, dynaClient)
		if err!= nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}
		return apiResponse(http.StatusOK, result)
	}

	result, err := user.FetchUsers(tableName, dynaClient)
	if err!= nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusOK, result)

}

func Store(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(
	*events.APIGatewayProxyResponse, error,
){
	result, err := user.CreateUser(req, tableName, dynaClient)
	if err!=nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusCreated, result)
}

func Update(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(
	*events.APIGatewayProxyResponse, error,
){
	result, err := user.UpdateUser(req, tableName, dynaClient)
	if err!= nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusOK, result)
}

func Destroy(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(
	*events.APIGatewayProxyResponse, error,
){
	err := user.DeleteUser(req, tableName, dynaClient)

	if err!= nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusOK, nil)
}

func NotFound()(*events.APIGatewayProxyResponse, error){
	return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}