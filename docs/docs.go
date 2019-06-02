// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-06-01 20:28:25.5790145 -0400 DST m=+0.057996001

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
                            "$ref": "#/definitions/main.GameComponent"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/main.GameComponent"
                        }
                    }
                }
            }
        },
        "/internal/game-definition": {
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
                "summary": "create Game definition",
                "parameters": [
                    {
                        "description": "GameDefinition",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/main.GameDefinition"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/main.GameDefinition"
                        }
                    }
                }
            }
        },
        "/internal/game-definition/{name}": {
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
                "summary": "find a game definition",
                "parameters": [
                    {
                        "type": "string",
                        "description": "GameDefinition name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/main.GameDefinition"
                        }
                    }
                }
            }
        },
        "/internal/luchador": {
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
                "summary": "find Luchador by ID",
                "parameters": [
                    {
                        "description": "FindLuchadorWithGamedefinition",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/main.FindLuchadorWithGamedefinition"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/main.GameComponent"
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
        "/internal/start-match/{name}": {
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
                        "type": "string",
                        "description": "GameDefinition name",
                        "name": "name",
                        "in": "path",
                        "required": true
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
        "/private/game-definition-all": {
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
                "summary": "find all game definitions",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.GameDefinition"
                            }
                        }
                    }
                }
            }
        },
        "/private/game-definition-id/{id}": {
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
                "summary": "find a game definition",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "GameDefinition id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/main.GameDefinition"
                        }
                    }
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
                            "$ref": "#/definitions/main.GameComponent"
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
                            "$ref": "#/definitions/main.GameComponent"
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
                                "$ref": "#/definitions/main.GameComponent"
                            }
                        }
                    }
                }
            }
        },
        "/private/match-single": {
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
                "summary": "find one match",
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
                            "type": "object",
                            "$ref": "#/definitions/main.Match"
                        }
                    }
                }
            }
        },
        "/private/tutorial": {
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
                "summary": "find tutorial GameDefinition",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.GameDefinition"
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
                "event": {
                    "type": "string"
                },
                "exception": {
                    "type": "string"
                },
                "gameDefinition": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "script": {
                    "type": "string"
                }
            }
        },
        "main.Config": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "key": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "main.FindLuchadorWithGamedefinition": {
            "type": "object",
            "properties": {
                "gameDefinitionID": {
                    "type": "integer"
                },
                "luchadorID": {
                    "type": "integer"
                }
            }
        },
        "main.GameComponent": {
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
                "gameDefinition": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "userID": {
                    "type": "integer"
                }
            }
        },
        "main.GameDefinition": {
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
                "codes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.Code"
                    }
                },
                "description": {
                    "type": "string"
                },
                "duration": {
                    "type": "integer"
                },
                "energy": {
                    "type": "integer"
                },
                "fireEnergyCost": {
                    "type": "integer"
                },
                "fps": {
                    "type": "integer"
                },
                "gameComponents": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.GameComponent"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "increaseSpeedEnergyCost": {
                    "type": "integer"
                },
                "increaseSpeedPercentage": {
                    "type": "integer"
                },
                "label": {
                    "type": "string"
                },
                "life": {
                    "type": "integer"
                },
                "luchadorSize": {
                    "type": "integer"
                },
                "maxFireAmount": {
                    "type": "integer"
                },
                "maxFireCooldown": {
                    "type": "integer"
                },
                "maxFireDamage": {
                    "type": "integer"
                },
                "maxParticipants": {
                    "type": "integer"
                },
                "minFireAmount": {
                    "type": "integer"
                },
                "minFireDamage": {
                    "type": "integer"
                },
                "minParticipants": {
                    "type": "integer"
                },
                "moveSpeed": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "punchAngle": {
                    "type": "integer"
                },
                "punchCoolDown": {
                    "type": "integer"
                },
                "punchDamage": {
                    "type": "integer"
                },
                "radarAngle": {
                    "type": "integer"
                },
                "radarRadius": {
                    "type": "integer"
                },
                "recycledLuchadorEnergyRestore": {
                    "type": "integer"
                },
                "respawnCooldown": {
                    "type": "integer"
                },
                "restoreEnergyperSecond": {
                    "type": "integer"
                },
                "sceneComponents": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.SceneComponent"
                    }
                },
                "sortOrder": {
                    "type": "integer"
                },
                "suggestedCodes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.Code"
                    }
                },
                "turnGunSpeed": {
                    "type": "integer"
                },
                "turnSpeed": {
                    "type": "integer"
                },
                "type": {
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
        "main.Match": {
            "type": "object",
            "properties": {
                "gameDefinitionID": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "lastTimeAlive": {
                    "type": "string"
                },
                "participants": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.GameComponent"
                    }
                },
                "timeEnd": {
                    "type": "string"
                },
                "timeStart": {
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
                "deaths": {
                    "type": "integer"
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
                }
            }
        },
        "main.SceneComponent": {
            "type": "object",
            "properties": {
                "blockMovement": {
                    "type": "boolean"
                },
                "codes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.Code"
                    }
                },
                "colider": {
                    "type": "boolean"
                },
                "gameDefinition": {
                    "type": "integer"
                },
                "height": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "respawn": {
                    "type": "boolean"
                },
                "rotation": {
                    "type": "integer"
                },
                "showInRadar": {
                    "type": "boolean"
                },
                "width": {
                    "type": "integer"
                },
                "x": {
                    "type": "integer"
                },
                "y": {
                    "type": "integer"
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
                    "$ref": "#/definitions/main.GameComponent"
                }
            }
        },
        "main.User": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "main.UserSetting": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "lastOption": {
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
