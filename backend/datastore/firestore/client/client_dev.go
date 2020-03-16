// +build dev staging

package client

import (
	"context"
	"log"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"

	"github.com/mtlynch/whatgotdone/backend/gcp"
)

func New(ctx context.Context) (*firestore.Client, error) {
	log.Printf("creating dev-mode firestore client")
	// Fake service account creds for dev testing. They're not real credentials
	// anywhere, but firestore.NewClient requires something that looks like real
	// credentials, even if we're calling the firestore emulator.
	dummyCreds := []byte(strings.TrimSpace(`
{
	"type": "service_account",
	"project_id": "fake-project-id",
	"private_key_id": "8c172989cb023a8ae499aa951a0d319f9263cd7d",
	"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCx5oh84od1DS3/\nfVGKdpFbQ4hNIMAjLww4rNPB4g1IA5JS5hLHw3QjkFTYs6c/d0vrNop6D5fiQ/Qw\nlps8m/GETZNKEcJHHJm4iHwOKBSXjJ83dmRbNeBuIDOhTD7DvTlWxkfvyjVoAgtK\nlr8EQU/4jNMza9qBv2xi9SFUgWxlkifePD1gjZm0QZsJD9xiBXjaditNMRoMqUV4\nHt8RX3aTNYbLE7p9ftoDJpkuFvtmriSG/j9uUGwGvgZ0hF4mBxRmf68rO1yzcWWP\nPJVEabhco1b4zfTv0VCTi9ou3Q/Zs8kpkkY2mZXYQrOPQuwoc8+jCaGJD5ArW8N9\nwuL6tFs1AgMBAAECggEAQTG7mRCnFXVEAxoY1MZI1Io2HBXBc+Nc9jQX0jiWJ2rv\np6ObBEwTdqkA/v0vcGm2j7dIHh0yyv+eMGQw9ZAsfRC0xnMloEvR5bdWxxVXHoax\nHnErq+Vdnt38LcM0SSVCKxO07yJKWhhNrQL7c4K/3NU23ORMijntbYJpuX9Iixvd\nRT4zZWzgulgNbIY37MZFub3HCZD1lSuafBGJwy1KQA+WP+0mOjEKo/I8ULbo3ekh\n38kXnWNOuhKzhN+mQ4LjCnlCq27+FsspoOeloMWpv8/StFCAXPuU7MEmXq9rHa+G\nKemu3+ZZIwNsuQ5+CGVPDsSSd8/Dj0Qszq0nR3GcTQKBgQDYC/Y+kaPpajZB1/U9\nny+VHGxz3DEiXCzmy/GjQhexBDILae5SyUsh4TsLpW0iZoCGshILO2bAdYWMZAD3\n7b9jEUZndUwufbtrTTGqvbiGvPtd64LHylBJHdz0T2JlKeidIx5L/SgxNHv5tMIm\n/uguzRpxQMgTldy3nyV+lFnzEwKBgQDSzKllPFYfR76KkcEGhxQRz5aRiK85H9Q3\nHAML8VB5xMIzxBycFEUQParJCJd/Q3Vi+863bcXzRF/d2Vm99AyoGSgK9PyQEsyw\nCa9ZP6Ko7859SZcP4c79q8m/Ns/CthNg4Q1EubOjnuC0f12CRFFGgm9MSlg6aW8v\nUA9pyXZ5lwKBgB7nvyb+MINwZSiQGw3gmq7q7Py57/FpXCb736oqBzeUURBe6++9\nydij3o1w8aatIQ+jo38l1TIM3bjSiWzt/qXOT9L27Zns9IWJ+mPhVec4W4D48rFf\n2JJNClGMlZfBIfxwjKH0Ke64AlAbMnbfmhkvz+uJh9V9Z6CAzJ1J3YAvAoGBAKPu\nkVWvRHJzAtUUYH5JEex/+WIYX9wWypxI5n2lHqZzw2sqee1PPh5RNr28NsS7m1Bs\n7udrMOPsKnmGi+nTHvyjA6bxum/4jsHf5kOL311tkLGSRy4Mt0JDFFPltlB/9DYF\nDqKBoBgAeFMmMXwa0PH6gb9cmZxjXhn3MuVbzQzlAoGBAI5NMrLY6nbqTiw8qCkH\ne+wSWpapgAyYJjZ1YVyoFV9sQSwSIoEki0zOY+eQfq/yixZu8d8LsdjkcMEYoMtq\ns0O30IQ1yeDaKCgs8/YWlk0kP33z4O7YlLijuhS9kNZwvOoxQbeBwu1hWybKIi6A\ngdO+YKUzCj6fqJgaSu9QMcS9\n-----END PRIVATE KEY-----\n",
	"client_email": "dummy@fake-project-id.iam.gserviceaccount.com",
	"client_id": "1234",
	"auth_uri": "https://accounts.google.com/o/oauth2/auth",
	"token_uri": "https://oauth2.googleapis.com/token",
	"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
	"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/dummy%fake-project-id.iam.gserviceaccount.com"
}	`))
	return firestore.NewClient(ctx, gcp.ProjectID(), option.WithCredentialsJSON(dummyCreds))
}
