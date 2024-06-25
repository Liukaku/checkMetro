package metro

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	structs "github.com/liukaku/checkMetro/cmd/util"
)

func ternaryString(boolean bool, string1 string, string2 string) string {
	if boolean {
		return string1
	} else {
		return string2
	}
}

//sorry future me for this type, should probably make a struct
func getRecipientInfo(emailTo string)[1]map[string]map[string]string{
	recipientsMap := map[string]map[string]string {
		"address": {
			"email": emailTo,
		},
	}

	recipientsArr := [1]map[string]map[string]string{recipientsMap}

	return recipientsArr
}

func getEmailContent(emailFrom string, subjectLine string, bodyContent string) structs.EmailContent{

	emailFromMap := map[string]map[string]string {
		"from": {
			"email": emailFrom,
			"name": "Metrolink Check",
		},
	}

	htmlBody := fmt.Sprintf(`
	<html>
            <body>
            <div>
            <h2>%s</h2>
            </div>
            </body>
            </html>
	`, bodyContent)

	returnVal := structs.EmailContent{From: emailFromMap, Subject: subjectLine, Html: htmlBody}

	return returnVal
}

func SendEmail(success bool, emailAuth string, emailTo string, emailFrom string, postmanUrl string) {
	subjectLine := ternaryString(success, "Metrolink is running", "Metrolink stops need updating")
	bodyContent := ternaryString(success, "hooray it works this week", "bugger")

	recipients := getRecipientInfo(emailTo)
	emailContent := getEmailContent(emailFrom, subjectLine, bodyContent)

	postBody := structs.EmailBody{Recipients: recipients, Content: emailContent}

	postBodyJson, _ := json.Marshal(postBody)
	
	postBodyBuffer := bytes.NewBuffer(postBodyJson)


	httpClient := &http.Client{}
	req, _ := http.NewRequest("POST", postmanUrl, postBodyBuffer)
	req.Header.Set("Authorization", emailAuth)
	req.Header.Set("Content-Type", "application/json")
	res, err := httpClient.Do(req)

	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
