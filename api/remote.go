package api

import (
	"fmt"
	"github.com/asiainfoLDP/datafoundry_plan/models"
	"github.com/asiainfoLDP/datafoundry_plan/openshift"
	userapi "github.com/openshift/origin/pkg/user/api/v1"
	kapi "k8s.io/kubernetes/pkg/api/v1"
	"os"
)

const RemoteAddr = "dev.dataos.io:8443"

//=================================================
//get remote endpoint
//=================================================

var (
	RechargeSercice string
	DataFoundryHost string
)

func BuildServiceUrlPrefixFromEnv(name string, isHttps bool, addrEnv string, portEnv string) string {
	var addr string
	if models.SetPlatform {
		addr = RemoteAddr
	} else {
		addr = os.Getenv(addrEnv)
	}
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

	//RechargeSercice = BuildServiceUrlPrefixFromEnv("ChargeSercice", false, os.Getenv("ENV_NAME_DATAFOUNDRYRECHARGE_SERVICE_HOST"), os.Getenv("ENV_NAME_DATAFOUNDRYRECHARGE_SERVICE_PORT"))
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
