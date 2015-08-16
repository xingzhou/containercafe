package limit

import(
	"conf"
)


func OpenConn(resource string, max int) bool{
	if max == 0 {
		return true
	}

	//incr and get counter (from Redis)
	count := conf.RedisIncr(resource)

	//if counter > max {decrement counter; return false} else {return true}
	if count > max {
		conf.RedisDecr(resource)
		return false
	}

	return true
}

func CloseConn(resource string, max int){
	if max == 0 {
		return
	}

	conf.RedisDecr(resource)
}
