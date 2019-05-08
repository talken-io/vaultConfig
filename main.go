package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// setting struct
type Setting struct {
	VaultApi string                            `json:"vaultApi"`
	Secrets  map[string]map[string]interface{} `json:"secrets"`
	Roles    map[string]Role                   `json:"roles"`
}

type Role struct {
	Secret   []string `json:"secret"`
	Password string   `json:"password"`
}

var setting Setting

// result struct
type Result struct {
	rolename string
	username string
	password string
}

var result []Result

// operator access token
var operatorAccessToken string

func exitHook() {
	if r := recover(); r != nil {
		fmt.Println("Error", r)
	}
}

func main() {
	var settingFile string
	if len(os.Args) > 1 {
		settingFile = os.Args[1]
	} else {
		settingFile = "setting.json"
	}

	defer exitHook()

	readSetting(settingFile)

	if setting.VaultApi == "" {
		fmt.Println("vaultApi must be specified in setting file")
		os.Exit(1)
	}

	fmt.Print("Enter Vault(" + setting.VaultApi + ") login token : ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	operatorAccessToken = scanner.Text()

	for name, kv := range setting.Secrets {
		makeSecret(name, kv)
		makeSecretPolicy(name)
	}

	for name, role := range setting.Roles {
		makeAppRole(name, role)
		makeAppRolePolicy(name)
		makeAppRoleUser(name, role)
	}

	fmt.Println("===== Role creation result =====")
	fmt.Println()
	for _, r := range result {
		fmt.Printf("Role [%s] created\n username : %s\n password : %s\n\n", r.rolename, r.username, r.password)
	}
}

func readSetting(filePath string) {
	// read setting file
	jsonFile, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	fmt.Println("Read setting.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &setting)
	if err != nil {
		panic(err)
	}

	// check integrity
	for r, v := range setting.Roles {
		for _, s := range v.Secret {
			_, exists := setting.Secrets[s]
			if !exists {
				panic("for role " + r + " secret " + s + " is not found on secrets")
			}
		}

		if v.Password == "" {
			panic("for role " + r + " no password set")
		}
	}

	// print
	fmt.Println("List of Secrets")
	for s, d := range setting.Secrets {
		fmt.Println(" - Secret : " + s)
		for k, v := range d {
			fmt.Printf("   %s : %v\n", k, v)
		}
		fmt.Println()
	}

	fmt.Println("List of Roles")
	for r, d := range setting.Roles {
		fmt.Println(" - Role : ", r)
		fmt.Println("   Secret list : ", d.Secret)
		fmt.Println("   Password : ", d.Password)
		fmt.Println()
	}
}

func sendToVault(method string, url string, data interface{}) {
	kvBytes, _ := json.Marshal(data)
	body := bytes.NewBuffer(kvBytes)

	req, err := http.NewRequest(method, setting.VaultApi+"/v1/"+url, body)
	if err != nil {
		panic(err)
	}

	req.Header.Add("X-Vault-Token", operatorAccessToken)

	fmt.Print("Request " + url + " : ")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		fmt.Println("OK")
	} else {
		fmt.Println(resp.StatusCode)
		panic(resp.Status)
	}
}

func makeSecret(s string, kv map[string]interface{}) {
	if kv != nil {
		sendToVault("POST", "secret/credentials/"+s, kv)
	} else {
		fmt.Println("secret", s, "is null, skip")
	}
}

func makeSecretPolicy(s string) {
	kv := map[string]string{"policy": "path \"secret/credentials/" + s + "\" { capabilities = [\"read\"] }"}

	sendToVault("PUT", "sys/policy/secret-"+s+"-policy", kv)
}

func makeAppRole(r string, role Role) {
	sb := strings.Builder{}

	for index, v := range role.Secret {
		sb.WriteString("secret-" + v + "-policy")
		if index < len(role.Secret)-1 {
			sb.WriteString(",")
		}
	}

	kv := make(map[string]string)
	kv["secret_id_ttl"] = "10s"
	kv["secret_id_num_uses"] = "1"
	kv["period"] = "30m"
	kv["policies"] = sb.String()

	sendToVault("POST", "auth/approle/role/"+r, kv)
}

func makeAppRolePolicy(r string) {
	kv := map[string]string{"policy": "path \"auth/approle/role/" + r + "/role-id\" { capabilities = [\"read\"] }\npath \"auth/approle/role/" + r + "/secret-id\" { capabilities = [\"read\", \"create\", \"update\", \"delete\"] }"}
	sendToVault("PUT", "sys/policy/approle-"+r+"-policy", kv)
}

func makeAppRoleUser(rolename string, role Role) {
	//var password string
	//if role.Password == "" {
	//	pass := make([]byte, 32)
	//	rand.Read(pass)
	//	password = hex.EncodeToString(pass)
	//} else {
	//	password =
	//}

	h := sha256.New()
	h.Write([]byte(rolename))
	username := hex.EncodeToString(h.Sum(nil))

	kv := make(map[string]string)
	kv["password"] = role.Password
	kv["ttl"] = "5s"
	kv["max_ttl"] = "5s"
	kv["policies"] = "approle-" + rolename + "-policy"

	sendToVault("POST", "auth/userpass/users/"+username, kv)

	result = append(result, Result{rolename: rolename, username: username, password: role.Password})
}
