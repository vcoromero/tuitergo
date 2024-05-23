package main

import (
	"os"
)

func main() {
}

// func CallLambda(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
// 	var res *events.APIGatewayProxyResponse

// 	awsgo.InitializeAWS()

// 	if !ValidatedParameters() {
// 		res = &events.APIGatewayProxyResponse{
// 			StatusCode: 400,
// 			Body:       "Error en las variables de entorno, deben incluir 'SecretName', 'BucketName' y 'UrlPrefix'",
// 			Headers: map[string]string{
// 				"Content-type": "application/json",
// 			},
// 		}
// 		return res, nil
// 	}
// }

func ValidatedParameters() bool {
	_, gotParameter := os.LookupEnv("SecretName")
	if !gotParameter {
		return gotParameter
	}

	_, gotParameter = os.LookupEnv("BucketName")
	if !gotParameter {
		return gotParameter
	}

	_, gotParameter = os.LookupEnv("UrlPrefix")
	if !gotParameter {
		return gotParameter
	}
	return gotParameter
}
