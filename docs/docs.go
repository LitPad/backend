// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "consumes": [
        "application/json"
    ],
    "produces": [
        "application/json"
    ],
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/google": {
            "post": {
                "description": "` + "`" + `This endpoint generates new access and refresh tokens for authentication via google` + "`" + `\n` + "`" + `Pass in token gotten from gsi client authentication here in payload to retrieve tokens for authentication` + "`" + `",
                "tags": [
                    "Auth"
                ],
                "summary": "Login a user via google",
                "parameters": [
                    {
                        "description": "User login",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemas.SocialLoginSchema"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/schemas.ResponseSchema"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "This endpoint generates new access and refresh tokens for authentication",
                "tags": [
                    "Auth"
                ],
                "summary": "Login a user",
                "parameters": [
                    {
                        "description": "User login",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemas.LoginSchema"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/schemas.ResponseSchema"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/logout": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "This endpoint logs a user out from our application",
                "tags": [
                    "Auth"
                ],
                "summary": "Logout a user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.ResponseSchema"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/refresh": {
            "post": {
                "description": "This endpoint refresh tokens by generating new access and refresh tokens for a user",
                "tags": [
                    "Auth"
                ],
                "summary": "Refresh tokens",
                "parameters": [
                    {
                        "description": "Refresh token",
                        "name": "refresh",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemas.RefreshTokenSchema"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/schemas.ResponseSchema"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "` + "`" + `This endpoint registers new users into our application.` + "`" + `",
                "tags": [
                    "Auth"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "User data",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemas.RegisterUser"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/schemas.RegisterResponseSchema"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/resend-verification-email": {
            "post": {
                "description": "` + "`" + `This endpoint resends new otp to the user's email.` + "`" + `",
                "tags": [
                    "Auth"
                ],
                "summary": "Resend Verification Email",
                "parameters": [
                    {
                        "description": "Email data",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemas.EmailRequestSchema"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.ResponseSchema"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/send-password-reset-otp": {
            "post": {
                "description": "` + "`" + `This endpoint sends new password reset otp to the user's email.` + "`" + `",
                "tags": [
                    "Auth"
                ],
                "summary": "Send Password Reset Otp",
                "parameters": [
                    {
                        "description": "Email object",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemas.EmailRequestSchema"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.ResponseSchema"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/set-new-password": {
            "post": {
                "description": "` + "`" + `This endpoint verifies the password reset otp.` + "`" + `",
                "tags": [
                    "Auth"
                ],
                "summary": "Set New Password",
                "parameters": [
                    {
                        "description": "Password reset object",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemas.SetNewPasswordSchema"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.ResponseSchema"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/verify-email": {
            "post": {
                "description": "` + "`" + `This endpoint verifies a user's email.` + "`" + `",
                "tags": [
                    "Auth"
                ],
                "summary": "Verify a user's email",
                "parameters": [
                    {
                        "description": "Verify Email object",
                        "name": "verify_email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemas.VerifyEmailRequestSchema"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.ResponseSchema"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/general/site-detail": {
            "get": {
                "description": "This endpoint retrieves few details of the site/application.",
                "tags": [
                    "General"
                ],
                "summary": "Retrieve site details",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.SiteDetailResponseSchema"
                        }
                    }
                }
            }
        },
        "/general/subscribe": {
            "post": {
                "description": "This endpoint creates a newsletter subscriber in our application",
                "tags": [
                    "General"
                ],
                "summary": "Add a subscriber",
                "parameters": [
                    {
                        "description": "Subscriber object",
                        "name": "subscriber",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Subscriber"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/schemas.SubscriberResponseSchema"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/healthcheck": {
            "get": {
                "description": "This endpoint checks the health of our application.",
                "tags": [
                    "HealthCheck"
                ],
                "summary": "HealthCheck",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/routes.HealthCheckSchema"
                        }
                    }
                }
            }
        },
        "/profiles/profile/{username}": {
            "get": {
                "description": "This endpoint views a user profile",
                "tags": [
                    "Profiles"
                ],
                "summary": "View User Profile",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username of user",
                        "name": "username",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.ResponseSchema"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/profiles/update": {
            "patch": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "This endpoint updates a user's profile",
                "tags": [
                    "Profiles"
                ],
                "summary": "Update User Profile",
                "parameters": [
                    {
                        "description": "Profile object",
                        "name": "profile",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemas.UpdateUserProfileSchema"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.ResponseSchema"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/profiles/update-password": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "This endpoint updates a user's password",
                "tags": [
                    "Profiles"
                ],
                "summary": "Update User Password",
                "parameters": [
                    {
                        "description": "Password object",
                        "name": "profile",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemas.UpdatePasswordSchema"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.ResponseSchema"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.SiteDetail": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string",
                    "example": "234, Lagos, Nigeria"
                },
                "email": {
                    "type": "string",
                    "example": "litpad@gmail.com"
                },
                "fb": {
                    "type": "string",
                    "example": "https://facebook.com"
                },
                "ig": {
                    "type": "string",
                    "example": "https://instagram.com"
                },
                "name": {
                    "type": "string"
                },
                "phone": {
                    "type": "string",
                    "example": "+234345434343"
                },
                "tw": {
                    "type": "string",
                    "example": "https://twitter.com"
                },
                "wh": {
                    "type": "string",
                    "example": "https://wa.me/2348133831036"
                }
            }
        },
        "models.Subscriber": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "minLength": 5,
                    "example": "johndoe@email.com"
                }
            }
        },
        "routes.HealthCheckSchema": {
            "type": "object",
            "properties": {
                "success": {
                    "type": "string",
                    "example": "pong"
                }
            }
        },
        "schemas.EmailRequestSchema": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "minLength": 5,
                    "example": "johndoe@email.com"
                }
            }
        },
        "schemas.LoginSchema": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "johndoe@email.com"
                },
                "password": {
                    "type": "string",
                    "example": "password"
                }
            }
        },
        "schemas.RefreshTokenSchema": {
            "type": "object",
            "required": [
                "refresh"
            ],
            "properties": {
                "refresh": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InNpbXBsZWlkIiwiZXhwIjoxMjU3ODk0MzAwfQ.Ys_jP70xdxch32hFECfJQuvpvU5_IiTIN2pJJv68EqQ"
                }
            }
        },
        "schemas.RegisterResponseSchema": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/schemas.EmailRequestSchema"
                },
                "message": {
                    "type": "string",
                    "example": "Data fetched/created/updated/deleted"
                },
                "status": {
                    "type": "string",
                    "example": "success"
                }
            }
        },
        "schemas.RegisterUser": {
            "type": "object",
            "required": [
                "email",
                "first_name",
                "last_name",
                "password",
                "username"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "minLength": 5,
                    "example": "johndoe@email.com"
                },
                "first_name": {
                    "type": "string",
                    "maxLength": 50,
                    "example": "John"
                },
                "last_name": {
                    "type": "string",
                    "maxLength": 50,
                    "example": "Doe"
                },
                "password": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 8,
                    "example": "strongpassword"
                },
                "terms_agreement": {
                    "type": "boolean"
                },
                "username": {
                    "type": "string",
                    "maxLength": 1000,
                    "example": "john-doe"
                }
            }
        },
        "schemas.ResponseSchema": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Data fetched/created/updated/deleted"
                },
                "status": {
                    "type": "string",
                    "example": "success"
                }
            }
        },
        "schemas.SetNewPasswordSchema": {
            "type": "object",
            "required": [
                "email",
                "otp",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "minLength": 5,
                    "example": "johndoe@email.com"
                },
                "otp": {
                    "type": "integer",
                    "example": 123456
                },
                "password": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 8,
                    "example": "newstrongpassword"
                }
            }
        },
        "schemas.SiteDetailResponseSchema": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/models.SiteDetail"
                },
                "message": {
                    "type": "string",
                    "example": "Data fetched/created/updated/deleted"
                },
                "status": {
                    "type": "string",
                    "example": "success"
                }
            }
        },
        "schemas.SocialLoginSchema": {
            "type": "object",
            "required": [
                "token"
            ],
            "properties": {
                "token": {
                    "type": "string",
                    "minLength": 10,
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InNpbXBsZWlkIiwiZXhwIjoxMjU3ODk0MzAwfQ.Ys_jP70xdxch32hFECfJQuvpvU5_IiTIN2pJJv68EqQ"
                }
            }
        },
        "schemas.SubscriberResponseSchema": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/models.Subscriber"
                },
                "message": {
                    "type": "string",
                    "example": "Data fetched/created/updated/deleted"
                },
                "status": {
                    "type": "string",
                    "example": "success"
                }
            }
        },
        "schemas.UpdatePasswordSchema": {
            "type": "object",
            "required": [
                "new_password",
                "old_password"
            ],
            "properties": {
                "new_password": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 8,
                    "example": "oldpassword"
                },
                "old_password": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 8,
                    "example": "newstrongpassword"
                }
            }
        },
        "schemas.UpdateUserProfileSchema": {
            "type": "object",
            "properties": {
                "username": {
                    "description": "Bio\t\t\t\t*string ` + "`" + `json:\"bio\"` + "`" + `",
                    "type": "string",
                    "maxLength": 1000,
                    "minLength": 3,
                    "example": "john-doe"
                }
            }
        },
        "schemas.VerifyEmailRequestSchema": {
            "type": "object",
            "required": [
                "email",
                "otp"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "minLength": 5,
                    "example": "johndoe@email.com"
                },
                "otp": {
                    "type": "integer",
                    "example": 123456
                }
            }
        },
        "utils.ErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "data": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "description": "Type 'Bearer jwt_string' to correctly set the API Key",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "4.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "LITPAD API",
	Description:      "`LitPAD API built with Fiber and GORM`",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
