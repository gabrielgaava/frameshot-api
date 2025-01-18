package utils

import (
	"example/web-service-gin/src/core/entity"
	"strings"
	"time"
)

func GetRequest() entity.Request {
	return entity.Request{
		ID:           1,
		UserId:       "user-123",
		UserEmail:    "user@example.com",
		VideoKey:     "video_input/teste.mp4",
		VideoSize:    200345,
		ZipOutputKey: "",
		Status:       entity.InProgress,
		CreatedAt:    time.Now(),
		FinishedAt:   time.Now(),
	}
}

func GetMockS3EventBody() string {
	return `{  
   "Records":[  
      {  
         "eventVersion":"2.1",
         "eventSource":"aws:s3",
         "awsRegion":"us-west-2",
         "eventTime":"1970-01-01T00:00:00.000Z",
         "eventName":"ObjectCreated:Put",
         "userIdentity":{  
            "principalId":"AIDAJDPLRKLG7UEXAMPLE"
         },
         "requestParameters":{  
            "sourceIPAddress":"127.0.0.1"
         },
         "responseElements":{  
            "x-amz-request-id":"C3D13FE58DE4C810",
            "x-amz-id-2":"FMyUVURIY8/IgAtTv8xRjskZQpcIZ9KG4V5Wp6S7S/JRWeUWerMUE5JgHvANOjpD"
         },
         "s3":{  
            "s3SchemaVersion":"1.0",
            "configurationId":"VideoUploaded",
            "bucket":{  
               "name":"amzn-s3-demo-bucket",
               "ownerIdentity":{  
                  "principalId":"A3NL1KOZZKExample"
               },
               "arn":"arn:aws:s3:::amzn-s3-demo-bucket"
            },
            "object":{  
               "key":"video_input/test.mp4",
               "size":1024,
               "eTag":"d41d8cd98f00b204e9800998ecf8427e",
               "versionId":"096fKKXTRTtl3on89fVO.nfljtsv6qko",
               "sequencer":"0055AED6DCD90281E5"
            }
         }
      }
   ]
}`
}

func GetMockOutputVideoEventBody(status string) string {
	body := `
	{
		"id": 1,
		"id_user" : "abc-123",
		"status": "${status}",
		"zip_s3_output": "zip_output/file.zip",
		"creation_date": "1970-01-01T00:00:00.000Z",
		"finished_date": "1970-01-01T00:00:00.000Z"
	}
	`

	return strings.ReplaceAll(body, "${status}", status)
}
