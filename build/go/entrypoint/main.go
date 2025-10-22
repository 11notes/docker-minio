package main

import (
	"os"
	"syscall"

	"github.com/11notes/go"
)

var(
	Eleven eleven.New = eleven.New{}
)

const ROOT_SSL string = "/minio/ssl"

func SSL(){
	caCertificate, err := Eleven.Container.GetSecret("MINIO_ROOT_CA_CRT", "MINIO_ROOT_CA_CRT_FILE")
	if err != nil {
		Eleven.LogFatal("ERR", "you must set MINIO_ROOT_CA_CRT or MINIO_ROOT_CA_CRT_FILE!")
	}
	Eleven.Util.WriteFile(ROOT_SSL + "/CAs/ca.crt", caCertificate)

	caKey, err := Eleven.Container.GetSecret("MINIO_ROOT_CA_KEY", "MINIO_ROOT_CA_KEY_FILE")
	if err != nil {
		Eleven.LogFatal("ERR", "you must set MINIO_ROOT_CA_KEY or MINIO_ROOT_CA_KEY_FILE!")
	}
	Eleven.Util.WriteFile(ROOT_SSL + "/CAs/ca.key", caKey)

	_, err = Eleven.Util.Run("/usr/local/bin/openssl", []string{"req", "-x509", "-newkey", "rsa:4096", "-sha256", "-days", "3650", "-nodes", "-keyout", ROOT_SSL + "/private.key", "-out", ROOT_SSL + "/public.crt", "-subj", "/CN=" + os.Getenv("HOSTNAME"), "-CA", ROOT_SSL + "/CAs/ca.crt", "-CAkey", ROOT_SSL + "/CAs/ca.key", "-addext", "subjectAltName=DNS:" + os.Getenv("HOSTNAME")})
	if err != nil {
		Eleven.LogFatal("ERR", "openssl: %s", err)
	}
}

func main() {
	SSL()

	password, err := Eleven.Container.GetSecret("MINIO_ROOT_PASSWORD", "MINIO_ROOT_PASSWORD_FILE")
	if err != nil {
		Eleven.LogFatal("ERR", "you must set MINIO_ROOT_PASSWORD or MINIO_ROOT_PASSWORD_FILE!")
	}

	if(len(os.Args) > 1){
		env := append(os.Environ(),"MINIO_ROOT_PASSWORD=" + password)
		if err := syscall.Exec("/usr/local/bin/minio", []string{"minio", "server", "--anonymous", "--json", "--certs-dir", ROOT_SSL, "--address", "0.0.0.0:9000", "--console-address", "0.0.0.0:3000", os.Args[1]}, env); err != nil {
			os.Exit(1)
		}
	}else{
		Eleven.LogFatal("ERR", "you must specify minio pool address!")	
	}
}