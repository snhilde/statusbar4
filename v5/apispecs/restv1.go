// This file contains the JSON implementation (v1) of the RestAPI spec.

package apispecs

// RESTV1 is version 1 of the REST API specification.
var RESTV1 = `
{
	"name": "REST API",
	"prefix": "/rest/v1",
	"description": "List of endpoints for the v1 REST API",
	"version": 1.0,
	"tables": [
		{
			"name": "general",
			"description": "General endpoints for using the API",
			"endpoints": [
				{
					"method": "GET",
					"url": "/ping",
					"description": "Ping the system.",
					"response": {
						"pong": "Simple \"pong\" response, not JSON-encoded."
					},
					"callback": "HandleGetPing"
				},
				{
					"method": "GET",
					"url": "/endpoints",
					"description": "Get a list of valid endpoints.",
					"response": {
						"endpoints": [
							{
								"method": "GET",
								"url": "/url",
								"description": "endpoint description"
							}
						]
					},
					"callback": "HandleGetEndpoints"
				}
			]
		},
		{
			"name": "routines",
			"description": "Endpoints related to accessing and manipulating rountines",
			"endpoints": [
				{
					"method": "GET",
					"url": "/routines",
					"description": "Get a list of information about all routines.",
					"response": {
						"routines": {
							"routineName": {
								"name": {
									"type": "string",
									"description": "Routine's name"
								},
								"uptime": {
									"type": "number",
									"description": "Routine's uptime, in seconds"
								},
								"interval": {
									"type": "number",
									"description": "Routine's update interval, in seconds"
								},
								"active": {
									"type": "boolean",
									"description": "Whether or not routine is currently active"
								}
							}
						}
					},
					"callback": "HandleGetRoutineAll"
				},
				{
					"method": "GET",
					"url": "/routines/:routine",
					"description": "Get information about the specified routine.",
					"response": {
						"routineName": {
							"name": {
								"type": "string",
								"description": "Routine's name"
							},
							"uptime": {
								"type": "number",
								"description": "Routine's uptime, in seconds"
							},
							"interval": {
								"type": "number",
								"description": "Routine's update interval, in seconds"
							},
							"active": {
								"type": "boolean",
								"description": "Whether or not routine is currently active"
							}
						}
					},
					"callback": "HandleGetRoutine"
				},

				{
					"method": "PUT",
					"url": "/routines",
					"description": "Restart all routines.",
					"callback": "HandlePutRoutineAll"
				},
				{
					"method": "PUT",
					"url": "/routines/:routine",
					"description": "Restart the specified routine.",
					"callback": "HandlePutRoutine"
				},

				{
					"method": "PATCH",
					"url": "/routines/:routine",
					"description": "Modify the specified routine's settings.",
					"request": {
						"interval": {
							"type": "number",
							"description": "New update interval, in seconds"
						}
					},
					"callback": "HandlePatchRoutine"
				},

				{
					"method": "DELETE",
					"url": "/routines",
					"description": "Stop all routines.",
					"callback": "HandleDeleteRoutineAll"
				},
				{
					"method": "DELETE",
					"url": "/routines/:routine",
					"description": "Stop the specified routine.",
					"callback": "HandleDeleteRoutine"
				}
			]
		}
	]
}
`
