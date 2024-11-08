package tests_test

import "github.com/gin-gonic/gin"

var (
	intKey            = ""
	authTokens        = []string{}
	authRefreshTokens = []string{}

	usersData = []gin.H{
		{
			"first_name": "TestName",
			"last_name":  "TestLastName",
			"username":   "test",
			"email":      "test@test.com",
			"password":   "test123456",
		},
	}

	organizationsData = []gin.H{
		{"name": "test", "description": "test desc"},
	}

	verificationsData = []gin.H{
		{"name": "test", "description": "test desc", "type": "SINGLE", "attributes": []gin.H{
			{"operator": "EQUAL", "value": "test"},
			{"operator": "EQUAL", "value": "test"},
		}},
		{"name": "test 2", "description": "test 2 desc", "type": "MULTI", "attributes": []gin.H{
			{"operator": "EQUAL", "value": "test"},
			{"operator": "EQUAL", "value": "test"},
		}},
	}

	schemasData = []gin.H{
		{"name": "test", "description": "test desc", "public": false, "attributes": []gin.H{
			{"name": "test_field", "description": "desc of field", "type": "TEXT"},
			{"name": "test_field_2", "description": "desc of field2", "type": "TEXT"},
		}},
		{"name": "test2", "description": "test2 desc", "public": false, "attributes": []gin.H{
			{"name": "test_field", "description": "desc of field", "type": "TEXT"},
			{"name": "test_field_2", "description": "desc of field2", "type": "TEXT"},
		}},
	}

	shorteningURLs = []string{
		"https://google.com",
		"https://app.socious.io",
	}

	recipientsData = []gin.H{
		{"first_name": "testName", "last_name": "test", "email": "test@test.com"},
		{"first_name": "testName2", "last_name": "test2", "email": "test2@test.com"},
		{"first_name": "testName3", "last_name": "test3", "email": "test3@test.com"},
	}

	credentialsData = []gin.H{
		{"name": "test", "description": "test desc", "claims": []gin.H{
			{"name": "test", "value": "valTest"},
			{"name": "test2", "value": 2},
		}},
		{"name": "test2", "description": "test 2 desc", "claims": []gin.H{
			{"name": "test", "value": "valTest"},
			{"name": "test2", "value": 2},
		}},
	}
)
