// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-05-14 01:00:23.0619851 -0400 DST m=+0.064065001

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
        "/internal/add-match-scores": {
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
                "summary": "saves a match score",
                "parameters": [
                    {
                        "description": "ScoreList",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/main.ScoreList"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/main.MatchScore"
                        }
                    }
                }
            }
        },
        "/internal/end-match": {
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
                "summary": "ends existing match",
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
        "/internal/game-component": {
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
                "summary": "Create Gamecomponent as Luchador",
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
        "/internal/luchador": {
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
                "summary": "find Luchador by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "int valid",
                        "name": "luchadorID",
                        "in": "query"
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
        "/internal/match-participant": {
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
                "summary": "Adds luchador to a match",
                "parameters": [
                    {
                        "description": "MatchParticipant",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/main.MatchParticipant"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/main.MatchParticipant"
                        }
                    }
                }
            }
        },
        "/internal/ready": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "summary": "returns application health check information",
                "responses": {
                    "200": {}
                }
            }
        },
        "/private/get-user": {
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
                "summary": "find The current user information",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/main.User"
                        }
                    }
                }
            }
        },
        "/private/join-match": {
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
                "summary": "Sends message with the request to join the match",
                "parameters": [
                    {
                        "description": "JoinMatch",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/main.JoinMatch"
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
                            "$ref": "#/definitions/main.UpdateLuchadorResponse"
                        }
                    }
                }
            }
        },
        "/private/mask-config/{id}": {
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
                "summary": "find maskConfig for a luchador",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Luchador ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.Config"
                            }
                        }
                    }
                }
            }
        },
        "/private/mask-random": {
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
                "summary": "create random maskConfig",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.Config"
                            }
                        }
                    }
                }
            }
        },
        "/private/match": {
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
                "summary": "find active matches",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.Match"
                            }
                        }
                    }
                }
            }
        },
        "/private/match-config": {
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
                "summary": "return luchador configs for current match",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "int valid",
                        "name": "matchID",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.Luchador"
                            }
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
        "main.Config": {
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
                "key": {
                    "type": "string"
                },
                "luchadorID": {
                    "type": "integer"
                },
                "updatedAt": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "main.JoinMatch": {
            "type": "object",
            "properties": {
                "luchadorID": {
                    "type": "integer"
                },
                "matchID": {
                    "type": "integer"
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
                "configs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.Config"
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
                "arenaHeight": {
                    "type": "integer"
                },
                "arenaWidth": {
                    "type": "integer"
                },
                "buletSpeed": {
                    "type": "integer"
                },
                "bulletSize": {
                    "type": "integer"
                },
                "createdAt": {
                    "type": "string"
                },
                "deletedAt": {
                    "type": "string"
                },
                "duration": {
                    "type": "integer"
                },
                "fps": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "lastTimeAlive": {
                    "type": "string"
                },
                "luchadorSize": {
                    "type": "integer"
                },
                "maxParticipants": {
                    "type": "integer"
                },
                "minParticipants": {
                    "type": "integer"
                },
                "participants": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.Luchador"
                    }
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
        "main.MatchParticipant": {
            "type": "object",
            "properties": {
                "luchadorID": {
                    "type": "integer"
                },
                "matchID": {
                    "type": "integer"
                }
            }
        },
        "main.MatchScore": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "deaths": {
                    "type": "integer"
                },
                "deletedAt": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "kills": {
                    "type": "integer"
                },
                "luchadorID": {
                    "type": "integer"
                },
                "matchID": {
                    "type": "integer"
                },
                "score": {
                    "type": "integer"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "main.ScoreList": {
            "type": "object",
            "properties": {
                "scores": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.MatchScore"
                    }
                }
            }
        },
        "main.UpdateLuchadorResponse": {
            "type": "object",
            "properties": {
                "errors": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "luchador": {
                    "type": "object",
                    "$ref": "#/definitions/main.Luchador"
                }
            }
        },
        "main.User": {
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
                "updatedAt": {
                    "type": "string"
                },
                "username": {
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
