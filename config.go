package main

import (
	"flag"
)

var (
	dbAddr     = flag.String("db_addr", "127.0.0.1:3306", "database address")
	dbUser     = flag.String("db_user", "root", "database user")
	dbName     = flag.String("db_name", "galaxy_s3_gateway", "database name")
	dbPassword = flag.String("db_passwd", "", "database password")
	port       = flag.Int("port", 5333, "proxy listening port")
	mongoDBAddr = flag.String("mongodb_addr", "127.0.0.1", "mongodb address")
	gfsMaster   = flag.String("gfs_zk_addr", "127.0.0.1:2181", "galaxyfs zookeeper cluster address")
	listenAddr  = flag.String("listen_addr", ":5050", "gateway serve address")
)

type config struct {
	dbAddr   string
	dbUser   string
	dbName   string
	dbPasswd string
	port     int
	mongoDBAddr string
}

func newConfig() *config {
	return &config{
		dbAddr:   *dbAddr,
		dbUser:   *dbUser,
		dbName:   *dbName,
		dbPasswd: *dbPassword,
		port:     *port,
		mongoDBAddr: *mongoDBAddr,
	}
}
