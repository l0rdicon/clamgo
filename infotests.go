
package main

import (

    "github.com/l0rdicon/btcjson"
    //"github.com/l0rdicon/btcutil"
    "crypto/sha256"
    "fmt"
    "errors"

    //"code.google.com/p/go.crypto/ripemd160"
    //"encoding/hex"
)

const (
    rpcuser = "username"
    rpcpassword = "password"
    rpchost = "localhost:33355"
)

func toSats(val int64) string {
    return fmt.Sprintf("%.8f", float64(val)/1e8)
}

func fromSats(val float64) int64 {
    return int64(val*1e8)
}

func dblSha256(data []byte) []byte {
    sha1 := sha256.New()
    sha2 := sha256.New()
    sha1.Write(data)
    sha2.Write(sha1.Sum(nil))
    return sha2.Sum(nil)
}



func main() {

    getinfo, _ := getinfo()
    getstakinginfo, _ := getstakinginfo()

    fmt.Println(getinfo)
    fmt.Println(getstakinginfo)

}


func getinfo() (*btcjson.InfoResult, error) {

        id := 1
        cmd, err := btcjson.NewGetInfoCmd(id)
        if err != nil {
                // Log and handle error.
        }

        // Send the message to server using the appropriate username and
        // password.
        reply, err := btcjson.RpcSend(rpcuser, rpcpassword, rpchost, cmd)
        if err != nil {
                fmt.Println(err)
                // Log and handle error.
        }

        // Ensure there is a result and type assert it to a btcjson.InfoResult.
        if reply.Result != nil {
                if info, ok := reply.Result.(*btcjson.InfoResult); ok {
                       return info, nil
                }
        }
    return nil, errors.New("Some error occured somewhere at sometime doing something to someone")
}


func getstakinginfo() (*btcjson.InfoResultStaking, error) {

        id := 1
        cmd, err := btcjson.NewGetStakingInfoCmd(id)
        if err != nil {
                // Log and handle error.
        }

        // Send the message to server using the appropriate username and
        // password.
        reply, err := btcjson.RpcSend(rpcuser, rpcpassword, rpchost, cmd)
        if err != nil {
                fmt.Println(err)
                // Log and handle error.
        }

        // Ensure there is a result and type assert it to a btcjson.InfoResult.
        if reply.Result != nil {
                if info, ok := reply.Result.(*btcjson.InfoResultStaking); ok {
                       return info, nil
                }
        }
    return nil, errors.New("Some error occured somewhere at sometime doing something to someone")
}
