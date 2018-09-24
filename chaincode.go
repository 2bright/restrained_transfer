package main

import (
    "fmt"
    "strings"
    "encoding/json"

    "github.com/hyperledger/fabric/core/chaincode/shim"
    pb "github.com/hyperledger/fabric/protos/peer"

    "github.com/shopspring/decimal"
)

type RestraintType byte

const (
    NONWAY RestraintType = '0' + iota
    ONEWAY
    ONEWAY_R
    TWOWAY
)

func (t RestraintType) isValid() bool {
    return t == NONWAY || t == ONEWAY || t == ONEWAY_R || t == TWOWAY
}

func (t RestraintType) reverse() RestraintType {
    t_int := t - '0'
    return RestraintType((((t_int & 1) << 1) | ((t_int & 2) >> 1)) + '0')
}

func (t RestraintType) allow() bool {
    return t == ONEWAY || t == TWOWAY
}

type MAP map[string]interface{}

type RestrainedTransferCC struct {
}

func (cc *RestrainedTransferCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
    // do nothing
    return shim.Success(nil)
}

func (cc *RestrainedTransferCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
    function, args := stub.GetFunctionAndParameters()
    switch function {
    case "register":
        return  cc.register(stub, args)
    case "getUserInfo":
        return  cc.getUserInfo(stub, args)
    case "getBalance":
        return  cc.getBalance(stub, args)
    case "recharge":
        return  cc.recharge(stub, args)
    case "withdraw":
        return  cc.withdraw(stub, args)
    case "transfer":
        return  cc.transfer(stub, args)
    case "setRestraint":
        return  cc.setRestraint(stub, args)
    case "getRestraintsOfUser":
        return  cc.getRestraintsOfUser(stub, args)
    case "getRestraintBetweenUsers":
        return  cc.getRestraintBetweenUsers(stub, args)
    default:
        return shim.Error("Error: unkown chaincode function " + function)
    }
}

// 注册用户
// 返回值：nil
func (cc *RestrainedTransferCC) register(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) != 2 {
        return shim.Error(`parameter error. usage: "{fcn: 'register', args: ['username', 'extras']}"`)
    }

    username := strings.TrimSpace(args[0])
    extras := args[1]

    if len(username) == 0 {
        return shim.Error("username should not be empty.")
    }

    var err error

    user_info_key, err := stub.CreateCompositeKey("u_i:", []string{username})
    if err != nil {
        return shim.Error("username is not valid. " + err.Error())
    }

    user_balance_key, err := stub.CreateCompositeKey("u_b:", []string{username})
    if err != nil {
        return shim.Error("username is not valid. " + err.Error())
    }

    userAsBytes, err := stub.GetState(user_info_key)
    if err != nil {
        return shim.Error("Failed to get state. " + err.Error())
    }
    if userAsBytes != nil {
        return shim.Error("username " + username + " already registered.")
    }

    balanceAsBytes, err := stub.GetState(user_balance_key)
    if err != nil {
        return shim.Error("Failed to get state. " + err.Error())
    }
    if balanceAsBytes != nil {
        return shim.Error("username " + username + " already registered.")
    }

	userAsBytes, err = json.Marshal(MAP{
        "name": username,
        "extras": extras,
    })
    if err != nil {
        return shim.Error("Failed to format user info. " + err.Error())
    }

    err = stub.PutState(user_info_key, userAsBytes)
    if err != nil {
        return shim.Error("Failed to put state. " + err.Error())
    }

    err = stub.PutState(user_balance_key, []byte(decimal.Zero.String()))
    if err != nil {
        return shim.Error("Failed to put state. " + err.Error())
    }

    return shim.Success(nil)
}

// 查询用户信息
// 返回值：json字符串
// {
//     "name":"user_a",
//     "extras":"balabala"
// }
func (cc *RestrainedTransferCC) getUserInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) != 1 {
        return shim.Error(`parameter error. usage: "{fcn: 'getUserInfo', args: ['username']}"`)
    }

    username := strings.TrimSpace(args[0])

    var err error

    user_info_key, err := stub.CreateCompositeKey("u_i:", []string{username})
    if err != nil {
        return shim.Error("username is not valid. " + err.Error())
    }

    userAsBytes, err := stub.GetState(user_info_key)
    if err != nil {
        return shim.Error("Failed to get state. " + err.Error())
    }
    if userAsBytes == nil {
        return shim.Error("username " + username + " is not registered.")
    }

    return shim.Success(userAsBytes)
}

// 查询余额
// 返回值: 十进制数 字符串，如 123.456
func (cc *RestrainedTransferCC) getBalance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) != 1 {
        return shim.Error(`parameter error. usage: "{fcn: 'getBalance', args: ['username']}"`)
    }

    username := strings.TrimSpace(args[0])

    var err error

    user_balance_key, err := stub.CreateCompositeKey("u_b:", []string{username})
    if err != nil {
        return shim.Error("username is not valid. " + err.Error())
    }

    balanceAsBytes, err := stub.GetState(user_balance_key)
    if err != nil {
        return shim.Error("Failed to get state. " + err.Error())
    }
    if balanceAsBytes == nil {
        return shim.Error("username " + username + " is not registered.")
    }

    return shim.Success(balanceAsBytes)
}

// 充值
// 返回值: nil
func (cc *RestrainedTransferCC) recharge(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) != 2 {
        return shim.Error(`parameter error. usage: "{fcn: 'recharge', args: ['username', 'amount']}"`)
    }

    username := strings.TrimSpace(args[0])
    amount_str := strings.TrimSpace(args[1])

    var err error

    user_balance_key, err := stub.CreateCompositeKey("u_b:", []string{username})
    if err != nil {
        return shim.Error("username is not valid. " + err.Error())
    }

    balanceAsBytes, err := stub.GetState(user_balance_key)
    if err != nil {
        return shim.Error("Failed to get state. " + err.Error())
    }
    if balanceAsBytes == nil {
        return shim.Error("username " + username + " is not registered.")
    }
    balance, err := decimal.NewFromString(string(balanceAsBytes))
    if err != nil {
        return shim.Error("Failed to parse balance stored. " + err.Error())
    }

    amount, err := decimal.NewFromString(amount_str)
    if err != nil {
        return shim.Error("Invalid recharge amount, expecting a number.")
    }
    if amount.LessThanOrEqual(decimal.Zero) {
        return shim.Error("Invalid recharge amount, expecting a number greater than 0.")
    }

    var new_balance = balance.Add(amount)

    err = stub.PutState(user_balance_key, []byte(new_balance.String()))
    if err != nil {
        return shim.Error("Failed to put state. " + err.Error())
    }

    return shim.Success(nil)
}

// 提款
// 提款金额需小于等于余额
// 返回值：nil
func (cc *RestrainedTransferCC) withdraw(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) != 2 {
        return shim.Error(`parameter error. usage: "{fcn: 'withdraw', args: ['username', 'amount']}"`)
    }

    username := strings.TrimSpace(args[0])
    amount_str := strings.TrimSpace(args[1])

    var err error

    user_balance_key, err := stub.CreateCompositeKey("u_b:", []string{username})
    if err != nil {
        return shim.Error("username is not valid. " + err.Error())
    }

    balanceAsBytes, err := stub.GetState(user_balance_key)
    if err != nil {
        return shim.Error("Failed to get state. " + err.Error())
    }
    if balanceAsBytes == nil {
        return shim.Error("username " + username + " is not registered.")
    }
    balance, err := decimal.NewFromString(string(balanceAsBytes))
    if err != nil {
        return shim.Error("Failed to parse balance stored. " + err.Error())
    }

    amount, err := decimal.NewFromString(amount_str)
    if err != nil {
        return shim.Error("Invalid recharge amount, expecting a number.")
    }
    if amount.LessThanOrEqual(decimal.Zero) {
        return shim.Error("Invalid recharge amount, expecting a number greater than 0.")
    }

    if balance.LessThan(amount) {
        return shim.Error("Failed recharge, not enough balance.")
    }

    var new_balance = balance.Sub(amount)

    err = stub.PutState(user_balance_key, []byte(new_balance.String()))
    if err != nil {
        return shim.Error("Failed to put state. " + err.Error())
    }

    return shim.Success(nil)
}

// 设置两个用户之间的转账约束
// 必需先通过 setRestraint 设置转账约束，才能调用 transfer 进行转账。即两个用户之间默认的转账约束是 0 即不允许转账。
// 转账约束有4个枚举值：
// 0 : a b 之间相互都不可转账
// 1 : a 可以转给 b, 但 b 不可转给 a
// 2 : b 可以转给 a, 但 a 不可转给 b
// 3 : a b 之间可以互转
// 返回值：nil
func (cc *RestrainedTransferCC) setRestraint(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) != 3 {
        return shim.Error(`parameter error. usage: "{fcn: 'setRestraint', args: ['username_a', 'username_b', 'restraint_type']}"`)
    }

    username_a := strings.TrimSpace(args[0])
    username_b := strings.TrimSpace(args[1])
    restraint_str := strings.TrimSpace(args[2])

    if username_a == username_b {
        return shim.Error("username_a and username_b must not be equal")
    }

    var err error

    if len(restraint_str) != 1 || !RestraintType(restraint_str[0]).isValid() {
        return shim.Error("transfer restraint type got " + restraint_str + ", expected one of [0, 1, 2, 3], respectively [forbid transfer, allow a to b, allow b to a, allow two-way]")
    }

    restraint := RestraintType(restraint_str[0])

    a_info_key, err := stub.CreateCompositeKey("u_i:", []string{username_a})
    if err != nil {
        return shim.Error("username_a is not valid. " + err.Error())
    }

    aAsBytes, err := stub.GetState(a_info_key)
    if err != nil {
        return shim.Error("Failed to get state. " + err.Error())
    }
    if aAsBytes == nil {
        return shim.Error("username_a " + username_a + " is not registered.")
    }

    b_info_key, err := stub.CreateCompositeKey("u_i:", []string{username_b})
    if err != nil {
        return shim.Error("username_b is not valid. " + err.Error())
    }

    bAsBytes, err := stub.GetState(b_info_key)
    if err != nil {
        return shim.Error("Failed to get state. " + err.Error())
    }
    if bAsBytes == nil {
        return shim.Error("username_b " + username_b + " is not registered.")
    }

    ab_restraint_key, err := stub.CreateCompositeKey("u_r:", []string{username_a, username_b})
    if err != nil {
        return shim.Error("username_a or username_b is not valid. " + err.Error())
    }

    ba_restraint_key, err := stub.CreateCompositeKey("u_r:", []string{username_b, username_a})
    if err != nil {
        return shim.Error("username_a or username_b is not valid. " + err.Error())
    }

    if restraint == NONWAY {
        err = stub.DelState(ab_restraint_key)
        if err != nil {
            return shim.Error("Failed to del state. " + err.Error())
        }

        err = stub.DelState(ba_restraint_key)
        if err != nil {
            return shim.Error("Failed to del state. " + err.Error())
        }
    } else {
        err = stub.PutState(ab_restraint_key, []byte{byte(restraint)})
        if err != nil {
            return shim.Error("Failed to put state. " + err.Error())
        }

        err = stub.PutState(ba_restraint_key, []byte{byte(restraint.reverse())})
        if err != nil {
            return shim.Error("Failed to put state. " + err.Error())
        }
    }

    return shim.Success(nil)
}

// 查询用户的所有转账约束
// 两个用户之间如果没有这种转账约束，则默认是 0 即相互都不可转账
// 枚举值的含义见 setRestraint
// 返回值：json字符串
// {
//     "user_x":"2"
//     "user_y":"1"
// }
func (cc *RestrainedTransferCC) getRestraintsOfUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) != 1 {
        return shim.Error(`parameter error. usage: "{fcn: 'getRestraintsOfUser', args: ['username']}"`)
    }

    username := strings.TrimSpace(args[0])

    var err error

    user_info_key, err := stub.CreateCompositeKey("u_i:", []string{username})
    if err != nil {
        return shim.Error("username is not valid. " + err.Error())
    }

    userAsBytes, err := stub.GetState(user_info_key)
    if err != nil {
        return shim.Error("Failed to get state. " + err.Error())
    }
    if userAsBytes == nil {
        return shim.Error("username " + username + " is not registered.")
    }

    itr, err := stub.GetStateByPartialCompositeKey("u_r:", []string{username})
    if err != nil {
        return shim.Error("Failed to get state. " + err.Error())
    }
    defer itr.Close()

    userRestraints := MAP{}

    for itr.HasNext() {
        kv, err := itr.Next()
        if err != nil {
            return shim.Error("Failed to get restraint stored. " + err.Error())
        }
        _, compositeKeyParts, err := stub.SplitCompositeKey(kv.Key)
        if err != nil {
            return shim.Error("Failed to parse restraint stored. " + err.Error())
        }
        userRestraints[compositeKeyParts[1]] = string(kv.Value)
    }

    userRestraintsAsBytes, err := json.Marshal(userRestraints)
    if err != nil {
        return shim.Error("Failed to format user restraints. " + err.Error())
    }

    return shim.Success(userRestraintsAsBytes)
}

// 查询从 username_a 到 username_b 的转账限制
// 两个用户之间如果没有这种转账约束，则默认是 "0" 即相互都不可转账
// 枚举值的含义见 setRestraint
// 返回值：字符串 "0" 或 "1" 或 "2" 或 "3"
func (cc *RestrainedTransferCC) getRestraintBetweenUsers(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) != 2 {
        return shim.Error(`parameter error. usage: "{fcn: 'getRestraintBetweenUsers', args: ['username_a', 'username_b']}"`)
    }

    username_a := strings.TrimSpace(args[0])
    username_b := strings.TrimSpace(args[1])

    var err error

    a_info_key, err := stub.CreateCompositeKey("u_i:", []string{username_a})
    if err != nil {
        return shim.Error("username_a is not valid. " + err.Error())
    }

    aAsBytes, err := stub.GetState(a_info_key)
    if err != nil {
        return shim.Error("Failed to get state. " + err.Error())
    }
    if aAsBytes == nil {
        return shim.Error("username_a " + username_a + " is not registered.")
    }

    b_info_key, err := stub.CreateCompositeKey("u_i:", []string{username_b})
    if err != nil {
        return shim.Error("username_b is not valid. " + err.Error())
    }

    bAsBytes, err := stub.GetState(b_info_key)
    if err != nil {
        return shim.Error("Failed to get state. " + err.Error())
    }
    if bAsBytes == nil {
        return shim.Error("username_b " + username_b + " is not registered.")
    }

    ab_restraint_key, err := stub.CreateCompositeKey("u_r:", []string{username_a, username_b})
    if err != nil {
        return shim.Error("username_a or username_b is not valid. " + err.Error())
    }

    abRestraintAsBytes, err := stub.GetState(ab_restraint_key)
    if err != nil {
        return shim.Error("Failed to get state. " + err.Error())
    }

    if abRestraintAsBytes != nil {
        return shim.Success(abRestraintAsBytes)
    } else {
        return shim.Success([]byte{byte(NONWAY)})
    }

}

// 从 username_a 转账给 username_b
// 如果 getRestraintBetweenUsers(stub, []string{username_a, username_b}) 返回值 不是 1 或 3，则链码返回 ERROR
// 返回值：nil
func (cc *RestrainedTransferCC) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) != 3 {
        return shim.Error(`parameter error. usage: "{fcn: 'transfer', args: ['username_a', 'username_b', 'amount']}"`)
    }

    username_a := strings.TrimSpace(args[0])
    username_b := strings.TrimSpace(args[1])
    amount_str := strings.TrimSpace(args[2])

    if username_a == username_b {
        return shim.Error("username_a and username_b must not be equal")
    }

    var err error

    a_balance_key, err := stub.CreateCompositeKey("u_b:", []string{username_a})
    if err != nil {
        return shim.Error("username_a is not valid. " + err.Error())
    }

    aAsBytes, err := stub.GetState(a_balance_key)
    if err != nil {
        return shim.Error("Failed to get state. " + err.Error())
    }
    if aAsBytes == nil {
        return shim.Error("username_a " + username_a + " is not registered.")
    }

    b_balance_key, err := stub.CreateCompositeKey("u_b:", []string{username_b})
    if err != nil {
        return shim.Error("username_b is not valid. " + err.Error())
    }

    bAsBytes, err := stub.GetState(b_balance_key)
    if err != nil {
        return shim.Error("Failed to get state. " + err.Error())
    }
    if bAsBytes == nil {
        return shim.Error("username_b " + username_b + " is not registered.")
    }

    ab_restraint_key, err := stub.CreateCompositeKey("u_r:", []string{username_a, username_b})
    if err != nil {
        return shim.Error("username_a or username_b is not valid. " + err.Error())
    }

    abRestraintAsBytes, err := stub.GetState(ab_restraint_key)
    if err != nil {
        return shim.Error("Failed to get state. " + err.Error())
    }

    restraint := NONWAY

    if abRestraintAsBytes != nil {
        restraint = RestraintType(abRestraintAsBytes[0])
    }

    if !restraint.allow() {
        return shim.Error(fmt.Sprintf("transfer from %s to %s is forbidden.", username_a, username_b))
    }

    balance_a, err := decimal.NewFromString(string(aAsBytes))
    if err != nil {
        return shim.Error("Failed to parse balance stored. " + err.Error())
    }

    balance_b, err := decimal.NewFromString(string(bAsBytes))
    if err != nil {
        return shim.Error("Failed to parse balance stored. " + err.Error())
    }

    amount, err := decimal.NewFromString(amount_str)
    if err != nil {
        return shim.Error("Invalid recharge amount, expecting a number.")
    }

    if balance_a.LessThan(amount) {
        return shim.Error("Failed transfer, not enough balance.")
    }

    new_balance_a := balance_a.Sub(amount)
    new_balance_b := balance_b.Add(amount)

    err = stub.PutState(a_balance_key, []byte(new_balance_a.String()))
    if err != nil {
        return shim.Error("Failed to put state. " + err.Error())
    }

    err = stub.PutState(b_balance_key, []byte(new_balance_b.String()))
    if err != nil {
        return shim.Error("Failed to put state. " + err.Error())
    }

    return shim.Success(nil)
}

func main() {
    if err := shim.Start(new(RestrainedTransferCC)); err != nil {
        fmt.Printf("Error starting RestrainedTransferCC chaincode: %s\n", err)
    }
}
