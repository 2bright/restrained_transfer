package main

import (
    "testing"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)

func testInit(t *testing.T, stub *shim.MockStub) {
    ret := stub.MockInit("1", nil)
    if ret.Status != shim.OK {
        t.Fatal("Init failed")
    }
    if ret.Payload != nil {
        t.Fatalf("Init return %s, expected nil", string(ret.Payload))
    }
}

func testRegister(t *testing.T, stub *shim.MockStub, username, extras string) {
    ret := stub.MockInvoke("1", [][]byte{[]byte("register"), []byte(username), []byte(extras)})
    if ret.Status != shim.OK {
        t.Fatalf("Invoke register failed. %s", ret.Message)
    }
    if ret.Payload != nil {
        t.Fatalf("Invoke register return %s, expected nil", string(ret.Payload))
    }
}

func testRegisterFail(t *testing.T, stub *shim.MockStub, username, extras string) {
    ret := stub.MockInvoke("1", [][]byte{[]byte("register"), []byte(username), []byte(extras)})
    if ret.Status != shim.ERROR {
        t.Fatalf("Invoke register for %s should fail.", username)
    }
}

func testGetUserInfo(t *testing.T, stub *shim.MockStub, username string, expected string) {
    ret := stub.MockInvoke("1", [][]byte{[]byte("getUserInfo"), []byte(username)})
    if ret.Status != shim.OK {
        t.Fatal("Invoke getUserInfo failed.")
    }
    if string(ret.Payload) != expected {
        t.Fatalf("Invoke getUserInfo return %s, expected %s", string(ret.Payload), expected)
    }
}

func testGetUserInfoFail(t *testing.T, stub *shim.MockStub, username string) {
    ret := stub.MockInvoke("1", [][]byte{[]byte("getUserInfo"), []byte(username)})
    if ret.Status != shim.ERROR {
        t.Fatalf("Invoke getUserInfo for %s should fail", username)
    }
}

func testGetBalance(t *testing.T, stub *shim.MockStub, username string, expected string) {
    ret := stub.MockInvoke("1", [][]byte{[]byte("getBalance"), []byte(username)})
    if ret.Status != shim.OK {
        t.Fatal("Invoke getBalance failed.")
    }
    if string(ret.Payload) != expected {
        t.Fatalf("Invoke getBalance return %s, expected %s", string(ret.Payload), expected)
    }
}

func testGetBalanceFail(t *testing.T, stub *shim.MockStub, username string) {
    ret := stub.MockInvoke("1", [][]byte{[]byte("getBalance"), []byte(username)})
    if ret.Status != shim.ERROR {
        t.Fatalf("Invoke getBalance for %s should fail", username)
    }
}

func testRecharge(t *testing.T, stub *shim.MockStub, username string, amount string) {
    ret := stub.MockInvoke("1", [][]byte{[]byte("recharge"), []byte(username), []byte(amount)})
    if ret.Status != shim.OK {
        t.Fatal("Invoke recharge failed.")
    }
    if ret.Payload != nil {
        t.Fatalf("Invoke recharge return %s, expected nil", string(ret.Payload))
    }
}

func testRechargeFail(t *testing.T, stub *shim.MockStub, username string, amount string) {
    ret := stub.MockInvoke("1", [][]byte{[]byte("recharge"), []byte(username), []byte(amount)})
    if ret.Status != shim.ERROR {
        t.Fatalf("Invoke recharge by amount %s should failed.", amount)
    }
}

func testWithdraw(t *testing.T, stub *shim.MockStub, username string, amount string) {
    ret := stub.MockInvoke("1", [][]byte{[]byte("withdraw"), []byte(username), []byte(amount)})
    if ret.Status != shim.OK {
        t.Fatal("Invoke withdraw failed.")
    }
    if ret.Payload != nil {
        t.Fatalf("Invoke withdraw return %s, expected nil", string(ret.Payload))
    }
}

func testWithdrawFail(t *testing.T, stub *shim.MockStub, username string, amount string) {
    ret := stub.MockInvoke("1", [][]byte{[]byte("withdraw"), []byte(username), []byte(amount)})
    if ret.Status != shim.ERROR {
        t.Fatalf("Invoke withdraw by amount %s should failed.", amount)
    }
}

func testGetRestraintsOfUser(t *testing.T, stub *shim.MockStub, username string, expected string) {
    ret := stub.MockInvoke("1", [][]byte{[]byte("getRestraintsOfUser"), []byte(username)})
    if ret.Status != shim.OK {
        t.Fatal("Invoke getRestraintsOfUser failed.")
    }
    if string(ret.Payload) != expected {
        t.Fatalf("Invoke getRestraintsOfUser return %s, expected %s", string(ret.Payload), expected)
    }
}

func testGetRestraintsOfUserFail(t *testing.T, stub *shim.MockStub, username string) {
    ret := stub.MockInvoke("1", [][]byte{[]byte("getRestraintsOfUser"), []byte(username)})
    if ret.Status != shim.ERROR {
        t.Fatalf("Invoke getRestraintsOfUser for %s should failed.", username)
    }
}

func testGetRestraintBetweenUsers(t *testing.T, stub *shim.MockStub, username_a, username_b string, expected string) {
    ret := stub.MockInvoke("1", [][]byte{[]byte("getRestraintBetweenUsers"), []byte(username_a), []byte(username_b)})
    if ret.Status != shim.OK {
        t.Fatal("Invoke getRestraintBetweenUsers failed." + ret.Message)
    }
    if string(ret.Payload) != expected {
        t.Fatalf("Invoke getRestraintBetweenUsers return %s, expected %s", string(ret.Payload), expected)
    }
}

func testGetRestraintBetweenUsersFail(t *testing.T, stub *shim.MockStub, username_a, username_b string) {
    ret := stub.MockInvoke("1", [][]byte{[]byte("getRestraintBetweenUsers"), []byte(username_a), []byte(username_b)})
    if ret.Status != shim.ERROR {
        t.Fatalf("Invoke getRestraintBetweenUsers for %s and %s should failed.", username_a, username_b)
    }
}

func testSetRestraint(t *testing.T, stub *shim.MockStub, username_a, username_b, restraint string) {
    ret := stub.MockInvoke("1", [][]byte{[]byte("setRestraint"), []byte(username_a), []byte(username_b), []byte(restraint)})
    if ret.Status != shim.OK {
        t.Fatal("Invoke setRestraint failed." + ret.Message)
    }
    if ret.Payload != nil {
        t.Fatalf("Invoke setRestraint return %s, expected nil", string(ret.Payload))
    }
}

func testSetRestraintFail(t *testing.T, stub *shim.MockStub, username_a, username_b, restraint string) {
    ret := stub.MockInvoke("1", [][]byte{[]byte("setRestraint"), []byte(username_a), []byte(username_b), []byte(restraint)})
    if ret.Status != shim.ERROR {
        t.Fatalf("Invoke setRestraint for %s and %s tobe %s should failed.", username_a, username_b, restraint)
    }
}

func testTransfer(t *testing.T, stub *shim.MockStub, username_a, username_b, amount string) {
    ret := stub.MockInvoke("1", [][]byte{[]byte("transfer"), []byte(username_a), []byte(username_b), []byte(amount)})
    if ret.Status != shim.OK {
        t.Fatal("Invoke transfer failed." + ret.Message)
    }
    if ret.Payload != nil {
        t.Fatalf("Invoke transfer return %s, expected nil", string(ret.Payload))
    }
}

func testTransferFail(t *testing.T, stub *shim.MockStub, username_a, username_b, amount string) {
    ret := stub.MockInvoke("1", [][]byte{[]byte("transfer"), []byte(username_a), []byte(username_b), []byte(amount)})
    if ret.Status != shim.ERROR {
        t.Fatalf("Invoke transfer from %s to %s with %s should failed.", username_a, username_b, amount)
    }
}

func TestInvoke(t *testing.T) {
    stub := shim.NewMockStub("TestInvoke", new(RestrainedTransferCC))
    testInit(t, stub)

    testRegisterFail(t, stub, "", "test empty username")
    testRegister(t, stub, "user_a", "company A")
    testRegisterFail(t, stub, "user_a", "company A")

    testGetUserInfo(t, stub, "user_a", `{"extras":"company A","name":"user_a"}`)
    testGetUserInfoFail(t, stub, "user_b")

    testGetBalance(t, stub, "user_a", "0")
    testGetBalanceFail(t, stub, "user_b")

    testGetRestraintsOfUser(t, stub, "user_a", "{}")
    testGetRestraintsOfUserFail(t, stub, "user_b")

    testGetRestraintBetweenUsersFail(t, stub, "user_a", "user_b")

    testRegister(t, stub, "user_b", "company A")

    testGetRestraintBetweenUsers(t, stub, "user_a", "user_b", "0")
    testGetRestraintBetweenUsersFail(t, stub, "user_a", "user_c")

    testRecharge(t, stub, "user_a", "10000")
    testGetBalance(t, stub, "user_a", "10000")
    testRecharge(t, stub, "user_a", "1000")
    testGetBalance(t, stub, "user_a", "11000")
    testRechargeFail(t, stub, "user_a", "-10000")

    testWithdraw(t, stub, "user_a", "10000")
    testGetBalance(t, stub, "user_a", "1000")
    testWithdraw(t, stub, "user_a", "1000")
    testGetBalance(t, stub, "user_a", "0")
    testWithdrawFail(t, stub, "user_a", "0")
    testWithdrawFail(t, stub, "user_a", "-1")

    testSetRestraint(t, stub, "user_a", "user_b", "1")
    testSetRestraintFail(t, stub, "user_a", "user_b", "")
    testSetRestraintFail(t, stub, "user_a", "user_b", "4")
    testSetRestraintFail(t, stub, "user_a", "user_b", "-1")
    testSetRestraintFail(t, stub, "user_a", "user_c", "0")

    testGetRestraintBetweenUsers(t, stub, "user_a", "user_b", "1")
    testGetRestraintBetweenUsers(t, stub, "user_b", "user_a", "2")
    testGetRestraintBetweenUsersFail(t, stub, "user_a", "user_c")

    testGetRestraintsOfUser(t, stub, "user_a", `{"user_b":"1"}`)
    testGetRestraintsOfUser(t, stub, "user_b", `{"user_a":"2"}`)

    testRegister(t, stub, "user_c", "company C")

    testSetRestraint(t, stub, "user_a", "user_c", "2")

    testGetRestraintsOfUser(t, stub, "user_a", `{"user_b":"1","user_c":"2"}`)
    testGetRestraintsOfUser(t, stub, "user_c", `{"user_a":"1"}`)

    testGetRestraintBetweenUsers(t, stub, "user_a", "user_c", "2")
    testGetRestraintBetweenUsers(t, stub, "user_c", "user_a", "1")

    testSetRestraint(t, stub, "user_a", "user_c", "0")

    testGetRestraintsOfUser(t, stub, "user_a", `{"user_b":"1"}`)
    testGetRestraintsOfUser(t, stub, "user_c", `{}`)

    testGetRestraintBetweenUsers(t, stub, "user_a", "user_c", "0")
    testGetRestraintBetweenUsers(t, stub, "user_c", "user_a", "0")

    testTransfer(t, stub, "user_a", "user_b", "0")

    testRecharge(t, stub, "user_a", "10000")
    testGetBalance(t, stub, "user_a", "10000")
    testGetBalance(t, stub, "user_b", "0")
    testGetBalance(t, stub, "user_c", "0")

    testTransfer(t, stub, "user_a", "user_b", "1000")
    testGetBalance(t, stub, "user_a", "9000")
    testGetBalance(t, stub, "user_b", "1000")

    testTransferFail(t, stub, "user_a", "user_c", "1000")
    testGetBalance(t, stub, "user_a", "9000")
    testGetBalance(t, stub, "user_c", "0")

    testTransferFail(t, stub, "user_b", "user_a", "1000")
    testGetBalance(t, stub, "user_a", "9000")
    testGetBalance(t, stub, "user_b", "1000")

    testSetRestraint(t, stub, "user_a", "user_b", "2")

    testTransfer(t, stub, "user_b", "user_a", "1000")
    testGetBalance(t, stub, "user_a", "10000")
    testGetBalance(t, stub, "user_b", "0")

    testTransferFail(t, stub, "user_b", "user_a", "1000")

    testSetRestraint(t, stub, "user_a", "user_b", "3")

    testTransfer(t, stub, "user_a", "user_b", "1000")
    testGetBalance(t, stub, "user_a", "9000")
    testGetBalance(t, stub, "user_b", "1000")

    testTransfer(t, stub, "user_b", "user_a", "100.123")
    testGetBalance(t, stub, "user_a", "9100.123")
    testGetBalance(t, stub, "user_b", "899.877")

    testSetRestraintFail(t, stub, "user_a", "user_a", "3")
    testTransferFail(t, stub, "user_a", "user_a", "5000")
}
