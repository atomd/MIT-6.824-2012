package lockservice

import "net/rpc"
import "crypto/rand"
import "math/big"

//
// the lockservice Clerk lives in the client
// and maintains a little state.
//
type Clerk struct {
  id int64 // client id
  servers [2]string // primary port, backup port
  req_id int64
}

// Generate numbers that have a high probability of being unique
func nrand() int64 {
  max := big.NewInt(int64(1) << 62)
  bigx, _ := rand.Int(rand.Reader, max)
  x := bigx.Int64()
  return x
}

func MakeClerk(primary string, backup string) *Clerk {
  ck := new(Clerk)
  ck.servers[0] = primary
  ck.servers[1] = backup
  ck.id = nrand()
  // ck.req_id = 0
  return ck
}

//
// call() sends an RPC to the rpcname handler on server srv
// with arguments args, waits for the reply, and leaves the
// reply in reply. the reply argument should be the address
// of a reply structure.
//
// call() returns true if the server responded, and false
// if call() was not able to contact the server. in particular,
// reply's contents are valid if and only if call() returned true.
//
// you should assume that call() will time out and return an
// error after a while if it doesn't get a reply from the server.
//
// please use call() to send all RPCs, in client.go and server.go.
// please don't change this function.
//
func call(srv string, rpcname string,
          args interface{}, reply interface{}) bool {
  c, errx := rpc.Dial("unix", srv)
  if errx != nil {
    return false
  }
  defer c.Close()

  err := c.Call(rpcname, args, reply)
  if err == nil {
    return true
  }
  return false
}


//
// ask the lock service for a lock.
// returns true if the lock service
// granted the lock, false otherwise.
//
// you will have to modify this function.
//
func (ck *Clerk) Lock(lockname string) bool {
  // prepare the arguments.
  ck.req_id++
  args := &LockArgs{ck.id, ck.req_id, lockname}
  var reply LockReply

  // if the primary does not respond, it sends an RPC to the backup.
  for i := 0; i < 2; i++ {
    // send an RPC request, wait for the reply.
    ok := call(ck.servers[i], "LockServer.Lock", args, &reply)
    if ok == true {
        return reply.OK
    }
  }
  return false
}


//
// ask the lock service to unlock a lock.
// returns true if the lock was previously held,
// false otherwise.
//

func (ck *Clerk) Unlock(lockname string) bool {

  ck.req_id++
  args := &UnlockArgs{ck.id, ck.req_id, lockname}
  var reply LockReply

  // if the primary does not respond, it sends an RPC to the backup.
  for i := 0; i < 2; i++ {
    // send an RPC request, wait for the reply.
    ok := call(ck.servers[i], "LockServer.Unlock", args, &reply)
    if ok == true {
        return reply.OK
    }
  }
  return false
}
