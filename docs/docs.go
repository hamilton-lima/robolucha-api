// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-01-01 16:57:36.600517097 -0500 EST m=+0.117706180

package docs

import (
	"bytes"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "swagger": "2.0",
    "info": {
        "description": "Robolucha API",
        "title": "Robolucha API",
        "contact": {},
        "license": {},
        "version": "1.0"
    },
    "host": "localhost:8080",
    "basePath": "/",
    "paths": {
        "/internal/match": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "create Match",
                "parameters": [
                    {
                        "description": "Match",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/main.Match"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/main.Match"
                        }
                    }
                }
            }
        },
        "/private/luchador": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "find or create Luchador for the current user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/main.Luchador"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Updates Luchador",
                "parameters": [
                    {
                        "description": "Luchador",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/main.Luchador"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/main.Luchador"
                        }
                    }
                }
            }
        },
        "/private/user/setting": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "find current user userSetting",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/main.UserSetting"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Updates user userSetting",
                "parameters": [
                    {
                        "description": "UserSetting",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/main.UserSetting"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/main.UserSetting"
                        }
                    }
                }
            }
        },
        "/public/login": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Logs the user",
                "parameters": [
                    {
                        "description": "LoginRequest",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/main.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/main.LoginResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.Code": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "deletedAt": {
                    "type": "string"
                },
                "event": {
                    "type": "string"
                },
                "exception": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "luchadorID": {
                    "type": "integer"
                },
                "script": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "main.LoginRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "main.LoginResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "boolean"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "main.Luchador": {
            "type": "object",
            "properties": {
                "codes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.Code"
                    }
                },
                "createdAt": {
                    "type": "string"
                },
                "deletedAt": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                },
                "userID": {
                    "type": "integer"
                }
            }
        },
        "main.Match": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "deletedAt": {
                    "type": "string"
                },
                "duration": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "lastTimeAlive": {
                    "type": "string"
                },
                "timeEnd": {
                    "type": "string"
                },
                "timeStart": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "main.UserSetting": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "deletedAt": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "lastOption": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                },
                "userID": {
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo swaggerInfo

type s struct{}

func (s *s) ReadDoc() string {
	t, err := template.New("swagger_info").Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, SwaggerInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
