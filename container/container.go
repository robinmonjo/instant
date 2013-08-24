package container

import(
  "fmt"
  "os/exec"
  "os"
  "log"
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

  fmt.Println("First usage, creating baseCN ...")

  out, err := exec.Command("/bin/bash", "sudo lxc-create -n " + base_cn_name + " -t " + lxc_template).Output()
  if err != nil {
    fmt.Println(out)
    log.Fatal(err)
  }

  fmt.Println(out)
  fmt.Println("Base container created")
}