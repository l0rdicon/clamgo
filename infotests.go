
package main

import (

    "github.com/l0rdicon/btcjson"
    "github.com/l0rdicon/btcutil"
    "crypto/sha256"
    "fmt"
    "errors"
    "database/sql"
    "github.com/coopernurse/gorp"
     _ "github.com/go-sql-driver/mysql"
     "log"
     "math"
     "strconv"
    //"code.google.com/p/go.crypto/ripemd160"
    //"encoding/hex"
)

type Digs struct {
    Id           int    `db:"id"`
    Txid         string `db:"txid"`
    BlockHeight  int64    `db:"height"`
    Vout         uint32    `db:"vout"`
    Address      string `db:"address"`
}



type CheckDigsInput struct {
    Txid string `json:"txid"`
    Vout uint32 `json:"vout"`
}

const (
    rpcuser = "l0rdicon"
    rpcpassword = "addsomepasswordhere"
    rpchost = "localhost:33355"
    initialDistAmount = 4.60545574
    distsets = 3208032
)

// Identifiers 
var clamID  = []byte{137}
var dogeID  = []byte{30}
var btcID  = []byte{0}
var ltcID  = []byte{48}

var dbmap *gorp.DbMap


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

func newDigs(txid string, blockheight int64, vout uint32, address string) Digs {
    return Digs{
        Txid: txid,
        BlockHeight:   blockheight,
        Vout:    vout,
        Address: address, 
    }
}

func newEmptyDigs() *Digs {
    return &Digs {}
}

func toClamAddress(string address) string {

    //Decode Base58 to []byte
    dbin := btcutil.Base58Decode(address)
    //Append clams Identifier
    addressclam := append(clamID[:], dbin[1:len(dbin)-4]...)
    //Append Checksum
    addresscomplete := append(addressclam, dblSha256(addressclam)[:4]...)
    //Encode back to Base58
    addressB58 := btcutil.Base58Encode(addresscomplete)

    return addressB58
}

func checkAddress(address string) (error) {

    dbin := btcutil.Base58Decode(address)
    //ident := []byte{dbin[0]}
    switch dbin[0] {
    case clamID[0]:
        fmt.Println("Clam")
        return nil
    case dogeID[0]:
        fmt.Println("doge")
        return nil
    case ltcID[0]:
        fmt.Println("ltc")
        return nil
    case btcID[0]:
        fmt.Println("btc")
        return nil
    }

   return errors.New("Not a valid Clam/BTC/Doge or LTC address")
}

func (t *Digs) CheckDug(hash string) (error) {
    fmt.Println(hash)
    err := dbmap.SelectOne(t, "SELECT * FROM digs where address = ?", hash)
    switch {
        case err == sql.ErrNoRows:
                return errors.New(fmt.Sprintf("Address not found in initial distrubution"))
        case err != nil:
                return  err
        default:
                return nil
        }
}


func main() {

    dbmap = initDb()
    defer dbmap.Db.Close()




  
    getinfo, _ := getinfo()
    ms := RenderFloat("#,###.##", getinfo.MoneySupply)
    ds := RenderFloat("#,###.##", getinfo.DigSupply)
    as := RenderFloat("#,###.##", getinfo.ActiveSupply)
    ss := RenderFloat("#,###.##", getinfo.StakeSupply)
    sets := RenderFloat("#,###.", getinfo.DigSupply/initialDistAmount)

    fmt.Println(ms)
    fmt.Println(ds)
    fmt.Println(as)
    fmt.Println(ss)
    fmt.Println(sets)
    //Run this to import a list of addresses included in init dist into mysql
    //includes txids and vout for each output. 
    //importDigs()


    /* example for checking is addresses was included in dist
    indist := newEmptyDigs()
    err := indist.CheckDug("xVhtAuCiHBoqPoXnQSEKUyKCL6WHF62Y2E")
    if err != nil {
        fmt.Println(err)
    }

    valid := validateoutputs([]btcjson.TransactionInput{btcjson.TransactionInput{Txid: indist.Txid, Vout: indist.Vout}})

    for _, l := range valid {
         fmt.Println(l.Status)
     }
     */

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

func validateoutputs(inputs []btcjson.TransactionInput) ([]btcjson.ValidateOutputs) {

        id := 1
        cmd, err := btcjson.NewValidateOutputsCmd(id, inputs)
        if err != nil {
            fmt.Println(err)    // Log and handle error.
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
                if info, ok := reply.Result.([]btcjson.ValidateOutputs); ok {

                       return info
                }
        }
        return nil
}




func importDigs() {

    for i := 300; i < 10000; i++ {
        blockhash := getblockhash(int64(i))
        block := getblock(blockhash)
        for l := range block.Tx {
            getrawtx := getrawtx(block.Tx[l])
            getdecodetx := decoderawtx(getrawtx)
             for g := range getdecodetx.Vout {
                dug := newDigs(block.Tx[l], block.Height, getdecodetx.Vout[g].N, getdecodetx.Vout[g].ScriptPubKey.Addresses[0])
                err := dbmap.Insert(&dug)
                checkErr(err, "Digs() DB Insert Error:")
                
            }

        }
        if i % 100 == 0  {
            fmt.Println(i) 
        }
    }   
}



func getblockhash(blknum int64) (string) {

        id := 1
        cmd, err := btcjson.NewGetBlockHashCmd(id, blknum)
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
        return reply.Result.(string)

}

func getblock(hash string) (*btcjson.BlockResult) {

        id := 1
        cmd, err := btcjson.NewGetBlockCmd(id, hash)
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
                if info, ok := reply.Result.(*btcjson.BlockResult); ok {
                       return info
                }
        }
        return nil
}



func decoderawtx(rawtx string) (*btcjson.TxRawDecodeResult) {



        id := 1
        cmd, err := btcjson.NewDecodeRawTransactionCmd(id, rawtx)
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
                if info, ok := reply.Result.(*btcjson.TxRawDecodeResult); ok {
                       return info
                }
        }

       return nil
}


func getrawtx(txid string) (string) {

        id := 1
        cmd, err := btcjson.NewGetRawTransactionCmd(id, txid)
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

       return reply.Result.(string)
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

func initDb() *gorp.DbMap {
    // connect to db using standard Go database/sql API
    // use whatever database/sql driver you wish
    
    //db, err := sql.Open("sqlite3", "/tmp/post_db.bin")
    db, err := sql.Open("mysql", "root:somepasswordhere@unix(/var/run/mysqld/mysqld.sock)/dug")
    checkErr(err, "sql.Open failed")

    // construct a gorp DbMap
    // dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
    dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

    // add a table, setting the table name to 'posts' and
    // specifying that the Id property is an auto incrementing PK
    dbmap.AddTableWithName(Digs{}, "digs").SetKeys(true, "Id")

    // create the table. in a production system you'd generally
    // use a migration tool, or create the tables via scripts
   // err = dbmap.CreateTablesIfNotExists()
    //checkErr(err, "Create tables failed")

    return dbmap
}

func checkErr(err error, msg string) {
    if err != nil {
        log.Fatalln(msg, err)
    }
}











var renderFloatPrecisionMultipliers = [10]float64{
    1,
    10,
    100,
    1000,
    10000,
    100000,
    1000000,
    10000000,
    100000000,
    1000000000,
}
 
var renderFloatPrecisionRounders = [10]float64{
    0.5,
    0.05,
    0.005,
    0.0005,
    0.00005,
    0.000005,
    0.0000005,
    0.00000005,
    0.000000005,
    0.0000000005,
}
 
func RenderFloat(format string, n float64) string {
    // Special cases:
    //   NaN = "NaN"
    //   +Inf = "+Infinity"
    //   -Inf = "-Infinity"
    if math.IsNaN(n) {
        return "NaN"
    }
    if n > math.MaxFloat64 {
        return "Infinity"
    }
    if n < -math.MaxFloat64 {
        return "-Infinity"
    }
 
    // default format
    precision := 2
    decimalStr := "."
    thousandStr := ","
    positiveStr := ""
    negativeStr := "-"
 
    if len(format) > 0 {
        // If there is an explicit format directive,
        // then default values are these:
        precision = 9
        thousandStr = ""
 
        // collect indices of meaningful formatting directives
        formatDirectiveChars := []rune(format)
        formatDirectiveIndices := make([]int, 0)
        for i, char := range formatDirectiveChars {
            if char != '#' && char != '0' {
                formatDirectiveIndices = append(formatDirectiveIndices, i)
            }
        }
 
        if len(formatDirectiveIndices) > 0 {
            // Directive at index 0:
            //   Must be a '+'
            //   Raise an error if not the case
            // index: 0123456789
            //        +0.000,000
            //        +000,000.0
            //        +0000.00
            //        +0000
            if formatDirectiveIndices[0] == 0 {
                if formatDirectiveChars[formatDirectiveIndices[0]] != '+' {
                    panic("RenderFloat(): invalid positive sign directive")
                }
                positiveStr = "+"
                formatDirectiveIndices = formatDirectiveIndices[1:]
            }
 
            // Two directives:
            //   First is thousands separator
            //   Raise an error if not followed by 3-digit
            // 0123456789
            // 0.000,000
            // 000,000.00
            if len(formatDirectiveIndices) == 2 {
                if (formatDirectiveIndices[1] - formatDirectiveIndices[0]) != 4 {
                    panic("RenderFloat(): thousands separator directive must be followed by 3 digit-specifiers")
                }
                thousandStr = string(formatDirectiveChars[formatDirectiveIndices[0]])
                formatDirectiveIndices = formatDirectiveIndices[1:]
            }
 
            // One directive:
            //   Directive is decimal separator
            //   The number of digit-specifier following the separator indicates wanted precision
            // 0123456789
            // 0.00
            // 000,0000
            if len(formatDirectiveIndices) == 1 {
                decimalStr = string(formatDirectiveChars[formatDirectiveIndices[0]])
                precision = len(formatDirectiveChars) - formatDirectiveIndices[0] - 1
            }
        }
    }
 
    // generate sign part
    var signStr string
    if n >= 0.000000001 {
        signStr = positiveStr
    } else if n <= -0.000000001 {
        signStr = negativeStr
        n = -n
    } else {
        signStr = ""
        n = 0.0
    }
 
    // split number into integer and fractional parts
    intf, fracf := math.Modf(n + renderFloatPrecisionRounders[precision])
 
    // generate integer part string
    intStr := strconv.Itoa(int(intf))
 
    // add thousand separator if required
    if len(thousandStr) > 0 {
        for i := len(intStr); i > 3; {
            i -= 3
            intStr = intStr[:i] + thousandStr + intStr[i:]
        }
    }
 
    // no fractional part, we can leave now
    if precision == 0 {
        return signStr + intStr
    }
 
    // generate fractional part
    fracStr := strconv.Itoa(int(fracf * renderFloatPrecisionMultipliers[precision]))
    // may need padding
    if len(fracStr) < precision {
        fracStr = "000000000000000"[:precision-len(fracStr)] + fracStr
    }
 
    return signStr + intStr + decimalStr + fracStr
}
 
func RenderInteger(format string, n int) string {
    return RenderFloat(format, float64(n))
}

