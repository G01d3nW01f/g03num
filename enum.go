package main

import(
	"io/ioutil"
	"log"
	"fmt"
	"os"
	"os/user"
	"os/exec"
	"path/filepath"
	"net"
	"strings"
)



func gatherInfo(){

	userInfo,err := user.Current()
	if err != nil{
		fmt.Printf("Error getting user information: %v\n",err)
		return
	}
	fmt.Printf("User: %s\n",userInfo.Username)
	
	passwdPath := "/etc/passwd"
	dir,file := filepath.Split(passwdPath)
	fmt.Println("[+]passwd file found")
	fmt.Printf("Directory: %s%s\n",dir,file)

	versionInfo,err := exec.Command("uname","-a").Output()
	if err != nil{
		fmt.Printf("Error getting system version information: %v\n",err)
		return
	}
	fmt.Printf("SystemVersion: %s\n",versionInfo)

}


func suid_check(){

	
	fmt.Println("[+]SUID set file")
	directory := "/"

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		// skip permission error
		if err != nil && os.IsPermission(err) {
		//	fmt.Printf("Permission error: %v\n", err)
			return nil
		} else if err != nil {
		//	fmt.Printf("Error: %v\n", err)
			return nil
		}

		// check suid file
		mode := info.Mode()
		if mode&os.ModeSetuid != 0 && mode.IsRegular() {
			fmt.Println(path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the directory: %v\n", err)
	}
}


func crontab_check(){
	crontabFile := "/etc/crontab"
	content,err := ioutil.ReadFile(crontabFile)

	if err != nil{
		log.Println("Error reading crontab files:",err)
	}

	fmt.Printf(string(content))
}


func get_env(){

	envInfo := os.Environ()	
	for _, envVar := range envInfo{
		fmt.Printf("Env variable: %s\n",envVar)
	}

	procInfo,err := os.Readlink("/proc/1/exe")
	if err != nil{
		fmt.Printf("Error getting system process info: %v\n",err)
		return
	}
	fmt.Printf("SystemProcess: %s\n",procInfo)
}


func get_networkInfo(){
	netInfo,err := net.Interfaces()
	if err != nil{
		fmt.Printf("Error getting system network connection info: %v\n",err)
		return
	}
	fmt.Printf("system network connections:\n")
	for _,netConn := range netInfo{
		addrs, err := netConn.Addrs()
		
		if err != nil{
			fmt.Printf("Error getting address for interface %s\n",netConn.Name,err)
			continue
		}
		fmt.Printf("Interface: %s\n",netConn.Name)
		fmt.Printf("MAC address: %s\n",netConn.HardwareAddr)
		for _, addr := range addrs{
			fmt.Printf("IP address:  %s\n",addr.String())
		}
	}
}

func netstat_chck(){
	cmd := exec.Command("netstat","-antulp")
	output,err := cmd.Output()
	if err != nil{
		fmt.Println("Error executing netstat command: ",err)
		return
	}
	lines := strings.Split(string(output),"\n")
	for _, line := range lines{
		if strings.Contains(line, "LISTEN"){
			fmt.Println(line)
		}

	}
	for _, line := range lines{
		if strings.Contains(line, "ESTABLISHED"){
			fmt.Println(line)
		}
	}
}


func get_kernel_sysinfo(){
	
	data,err := ioutil.ReadFile("/proc/modules")
	if err != nil{
		fmt.Println("Error reading /proc/modules: ",err)
		return
	}
	fmt.Println(string(data))
}

func get_application_info(){
	appInfo, err := exec.Command("dpkg","-l").Output()
	if err != nil{
		fmt.Println("Error getting system application information: %v\n",err)
	}
	fmt.Printf("System Applications:\n%s\n",appInfo)
}


func main(){

	gatherInfo()
	suid_check()
	crontab_check()
	get_env()
	get_networkInfo()
	netstat_chck()
	get_kernel_sysinfo()
	get_application_info()
}
