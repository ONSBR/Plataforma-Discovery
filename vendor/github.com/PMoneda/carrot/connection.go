package carrot

import "fmt"

type ConnectionConfig struct {
	Username string
	Password string
	Host     string
	AMQPPort string
	APIPort  string
	VHost    string
}

//GetAMQPURI returns amqp url format based on config
func (conn *ConnectionConfig) GetAMQPURI() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/%s", conn.Username, conn.Password, conn.Host, conn.GetAMQPPort(), conn.VHost)
}

//GetAMQPPort returns current config of amqp port or default port
func (conn *ConnectionConfig) GetAMQPPort() string {
	if conn.AMQPPort != "" {
		return conn.AMQPPort
	}
	return "5672"
}

//GetAPIURI return uri formatted
func (conn *ConnectionConfig) GetAPIURI() string {
	return fmt.Sprintf("http://%s:%s/", conn.Host, conn.GetAPIPort())
}

//GetAPIPort returns current api port or default
func (conn *ConnectionConfig) GetAPIPort() string {
	if conn.APIPort != "" {
		return conn.APIPort
	}
	return "15672"
}
