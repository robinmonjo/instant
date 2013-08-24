package container

import(
  "fmt"
  "os/exec"
  "os"
  "log"
  "io"
)

const lxc_template = "ubuntu"
const base_cn_name = "baseCN"
const base_cn_rootfs = "/var/lib/lxc/baseCN/rootfs"

func baseCnExists () bool {
  _, err := os.Stat(base_cn_rootfs)
  return err == nil
}

func CreateBaseCn () {

  if baseCnExists() {
    return
  }

  _, err := exec.LookPath("lxc")
  if err != nil {
    log.Fatal("lxc not installed")
  }

  fmt.Println("First usage, creating baseCN ...")

  cmd := exec.Command("sudo", "lxc-create", "-n", base_cn_name, "-t", lxc_template)
  stdout, err := cmd.StdoutPipe()
  if err != nil {
    log.Fatal(err)
  }
  stderr, err := cmd.StderrPipe()
  if err != nil {
    log.Fatal(err)
  }
  err = cmd.Start()
  if err != nil {
    log.Fatal(err)
  }

  io.Copy(os.Stdout, stdout) 
  io.Copy(os.Stderr, stderr) 

  cmd.Wait()
  fmt.Println("Base container created")
}

func CreateUserCn () {
  //create /containers
  //random id
  //aufs mount
  //single ip generate (keep a file on the disk (JSON) that stores cont id, ip, hwaddr)
  //single hwaddr generate
  //write config
  //return a container object (start command + sshCommand + pipe web_sockets to ssh stdin and ssh stdout to websocket)
}
