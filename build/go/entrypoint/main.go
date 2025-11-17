package main

import (
	"os"
	"syscall"
	"strings"

	"github.com/11notes/go"
)

var(
	Eleven eleven.New = eleven.New{}
)

const ROOT_SSL string = "/minio/ssl"

func CreateClusterCertificate(){
	// check for RootCA certificate of cluster
	caCertificate, err := Eleven.Container.GetSecret("MINIO_ROOT_CA_CRT", "MINIO_ROOT_CA_CRT_FILE")
	if err != nil {
		Eleven.LogFatal("you must set MINIO_ROOT_CA_CRT or MINIO_ROOT_CA_CRT_FILE!")
	}
	Eleven.Util.WriteFile(ROOT_SSL + "/CAs/ca.crt", caCertificate)

	// check for Root CA key of cluster
	caKey, err := Eleven.Container.GetSecret("MINIO_ROOT_CA_KEY", "MINIO_ROOT_CA_KEY_FILE")
	if err != nil {
		Eleven.LogFatal("you must set MINIO_ROOT_CA_KEY or MINIO_ROOT_CA_KEY_FILE!")
	}
	Eleven.Util.WriteFile(ROOT_SSL + "/CAs/ca.key", caKey)

	// create node certificate signed by Root CA
	_, err = Eleven.Util.Run("/usr/local/bin/openssl", []string{"req", "-x509", "-newkey", "rsa:4096", "-sha256", "-days", "3650", "-nodes", "-keyout", ROOT_SSL + "/private.key", "-out", ROOT_SSL + "/public.crt", "-subj", "/CN=" + os.Getenv("HOSTNAME"), "-CA", ROOT_SSL + "/CAs/ca.crt", "-CAkey", ROOT_SSL + "/CAs/ca.key", "-addext", "subjectAltName=DNS:" + os.Getenv("HOSTNAME")})
	if err != nil {
		Eleven.LogFatal("openssl: %s", err.Error())
	}
}

func main(){
	// check for root password
	password, err := Eleven.Container.GetSecret("MINIO_ROOT_PASSWORD", "MINIO_ROOT_PASSWORD_FILE")
	if err != nil {
		Eleven.LogFatal("you must set MINIO_ROOT_PASSWORD or MINIO_ROOT_PASSWORD_FILE!")
	}

	// prepare run statement and check if any additional user command is for a cluster or not
	if(len(os.Args) > 1){
		minio := []string{"minio", "server", "--anonymous", "--json", "--certs-dir", ROOT_SSL, "--address", "0.0.0.0:9000", "--console-address", "0.0.0.0:9001"}
		cluster := false
		args := os.Args[1:]
		for _, value := range args{
			minio = append(minio, value)
			if(strings.HasPrefix(value, "http")){
				cluster = true
			}
		}
		env := append(os.Environ(),"MINIO_ROOT_PASSWORD=" + password)
		if(cluster){
			// create cluster certificate if cluster
			CreateClusterCertificate()
		}else{
			// create stand alone certificate
			_, err = Eleven.Util.Run("/usr/local/bin/openssl", []string{"req", "-x509", "-newkey", "rsa:4096", "-sha256", "-days", "3650", "-nodes", "-keyout", ROOT_SSL + "/private.key", "-out", ROOT_SSL + "/public.crt", "-subj", "/CN=" + os.Getenv("HOSTNAME")})
			if err != nil {
				Eleven.LogFatal("openssl: %s", err.Error())
			}
		}

		// start minio
		if err := syscall.Exec("/usr/local/bin/minio", minio, env); err != nil {
			os.Exit(1)
		}
	}else{
		Eleven.LogFatal("you must specify minio pool address or /mnt!")	
	}
}