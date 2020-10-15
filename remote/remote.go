package remote

import (
	"gaad/common"
	"gaad/gssh"
	"gaad/gssh/auth"
	"gaad/gssh/comp"
	"gaad/models"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/ssh"
	"log"
	"strconv"
	"sync"
)

var (
	fileRWMutex sync.RWMutex
)

func InitDockerSwarmMaster(node models.Node) {
	clientConfig, err := auth.PasswordKey(node.Username, node.Password, ssh.InsecureIgnoreHostKey())
	if err != nil {
		log.Fatal(err)
		return
	}
	client := gssh.NewClient(node.Ip+":"+strconv.Itoa(node.Port), clientConfig)
	gssh.RunShellFile(client, "./shell/hello.sh")

}
func FollowDockerSwarmMaster(node models.Node) {

}
func InitDockerSwarmSlaver(node models.Node) {

}

func InitKubernetesMaster(node models.Node) {
	clientConfig, err := auth.PasswordKey(node.Username, node.Password, ssh.InsecureIgnoreHostKey())
	if err != nil {
		log.Fatal(err)
		return
	}
	client := gssh.NewClient(node.Ip+":"+strconv.Itoa(node.Port), clientConfig)

	fileRWMutex.Lock()

	configMap := make(map[string]string)

	/* map插入key - value对,各个国家对应的首都 */
	configMap["MASTER_HOST"] = "k8s-master-1"
	configMap["MASTER_IP"] = node.Ip
	//configMap["K8S_VERSION"] = "v1.18.6"

	common.GenConfigFile(configMap, "./script/k8s/config")
	uid := uuid.NewV4()
	compFileName := uid.String() + ".tar.gz"
	//沒有/tmp文件夾自动生成
	common.CreateFile("./tmp")
	compFile := "./tmp/" + compFileName
	comp.Compress("./script/k8s", compFile)

	fileRWMutex.Unlock()

	gssh.Scp(client, compFile, "/opt/", true)
	gssh.RunShell(client, "pwd;cd /opt/k8s;chmod +x -R ./*;./install.sh;")

}
func FollowKubernetesMaster(node models.Node) {

}
func InitKubernetesSlaver(node models.Node) {

}
