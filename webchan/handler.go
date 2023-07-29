package webchan

// type handlers struct {
// 	OnConnection    func(c *Connection, args interface{})
// 	OnAuth          func(args interface{}) error
// 	OnDisconnection func(c *Connection, message string)
// 	OnMessage       func(c *Connection, message []byte)
// }

type Handlers struct {
	OnConnection    func(c *Member, args interface{})
	OnAuth          func(args interface{}) error
	OnDisconnection func(c *Member, message string)
	OnMessage       func(c *Member, message []byte)
}
