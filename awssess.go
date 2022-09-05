package awssess

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
)

// Create a Session using profile values and get the Session Token
func awssess() *session.Session {

	//Enable CONFIG to pick the region from profile
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	var mfaCode string
	var profilename string
	profilename = os.Args[1]
	fmt.Println("--", profilename, "loaded --")
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: profilename,
	})

	_iam := iam.New(sess)
	devices, err := _iam.ListMFADevices(&iam.ListMFADevicesInput{})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Println("Error:", awsErr.Code(), awsErr.Message())
		}
	}
	sn := devices.MFADevices[0].SerialNumber

	fmt.Println(awsutil.StringValue(_iam.Config.Credentials))
	svc := sts.New(sess)
	fmt.Println("##ENTER 6 Digits MFA CODE")
	fmt.Scanln(&mfaCode)

	params := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(900),
		SerialNumber:    aws.String(*sn),
		TokenCode:       aws.String(mfaCode),
	}
	resp, err := svc.GetSessionToken(params)

	os.Setenv("AWS_SESSION_TOKEN", awsutil.StringValue(resp.Credentials.SessionToken))
	fmt.Println(awsutil.Prettify(devices.MFADevices[0].UserName), "Logged in Successfully")
	return sess
}
