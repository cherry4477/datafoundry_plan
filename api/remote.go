package api

import (
	"fmt"
	"github.com/asiainfoLDP/datafoundry_plan/common"
	"github.com/asiainfoLDP/datafoundry_plan/openshift"
	userapi "github.com/openshift/origin/pkg/user/api/v1"
	kapi "k8s.io/kubernetes/pkg/api/v1"
	"net/http"
	"os"
)

const (
	GeneralRemoteCallTimeout = 10 // seconds
)

//=================================================
//get remote endpoint
//=================================================

var (
	RechargeSercice string
	DataFoundryHost string
)

func BuildServiceUrlPrefixFromEnv(name string, isHttps bool, addrEnv string, portEnv string) string {
	addr := os.Getenv(addrEnv)
	if addr == "" {
		logger.Emergency("%s env should not be null", addrEnv)

	}
	if portEnv != "" {
		port := os.Getenv(portEnv)
		if port != "" {
			addr += ":" + port
		}
	}

	prefix := ""
	if isHttps {
		prefix = fmt.Sprintf("https://%s", addr)
	} else {
		prefix = fmt.Sprintf("http://%s", addr)
	}

	logger.Info("%s = %s", name, prefix)

	return prefix
}

func InitGateWay() {
	DataFoundryHost = BuildServiceUrlPrefixFromEnv("DataFoundryHost", true, "DATAFOUNDRY_HOST_ADDR", "")
	openshift.Init(DataFoundryHost, os.Getenv("DATAFOUNDRY_ADMIN_USER"), os.Getenv("DATAFOUNDRY_ADMIN_PASS"))

	RechargeSercice = BuildServiceUrlPrefixFromEnv("ChargeSercice", false, os.Getenv("ENV_NAME_DATAFOUNDRYRECHARGE_SERVICE_HOST"), os.Getenv("ENV_NAME_DATAFOUNDRYRECHARGE_SERVICE_PORT"))
}

//=============================================================
//get username
//=============================================================

func getDFUserame(token string) (string, error) {
	//Logger.Info("token = ", token)
	//if Debug {
	//	return "liuxu", nil
	//}

	user, err := authDF(token)
	if err != nil {
		return "", err
	}
	return dfUser(user), nil
}

func authDF(userToken string) (*userapi.User, error) {
	if Debug {
		return &userapi.User{
			ObjectMeta: kapi.ObjectMeta{
				Name: "local",
			},
		}, nil
	}

	u := &userapi.User{}
	osRest := openshift.NewOpenshiftREST(openshift.NewOpenshiftClient(userToken))
	uri := "/users/~"
	osRest.OGet(uri, u)
	if osRest.Err != nil {
		logger.Info("authDF, uri(%s) error: %s", uri, osRest.Err)
		return nil, osRest.Err
	}

	return u, nil
}

func dfUser(user *userapi.User) string {
	return user.Name
}

//====================================================
//call recharge api
//====================================================

func couponRecharge(adminToken, couponSerial, username, namespace string, amount float32) error {
	body := fmt.Sprintf(
		`{"namespace":"%s", "amount":%.3f, "reason":"%s", "user":"%s"}`,
		namespace, amount, couponSerial, username,
	)

	//RechargeSercice1 := "http://datafoundry.recharge.app.dataos.io:80"
	url := fmt.Sprintf("%s/charge/v1/couponrecharge", RechargeSercice)

	response, data, err := common.RemoteCallWithJsonBody("POST", url, adminToken, "", []byte(body))
	if err != nil {
		logger.Error("recharge err: %v", err)
		return err
	}

	if response.StatusCode != http.StatusOK {
		logger.Info("makeRecharge remote (%s) status code: %d. data=%s", url, response.StatusCode, string(data))
		return fmt.Errorf("makeRecharge remote (%s) status code: %d.", url, response.StatusCode)
	}

	return nil
}
