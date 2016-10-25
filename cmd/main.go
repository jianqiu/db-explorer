package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"errors"
	"database/sql"

	"code.cloudfoundry.org/cflager"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/clock"
	vps "github.com/jianqiu/vm-pool-server"
	"github.com/jianqiu/vm-pool-server/db"
	"github.com/jianqiu/vm-pool-server/db/sqldb"
	"github.com/go-sql-driver/mysql"
	"github.com/jianqiu/vm-pool-server/migration"
	"github.com/jianqiu/vm-pool-server/handlers"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/grouper"
	"github.com/tedsuo/ifrit/http_server"
	"github.com/tedsuo/ifrit/sigmon"
)

var accessLogPath = flag.String(
	"accessLogPath",
	"",
	"Location of the access log",
)

var listenAddress = flag.String(
	"listenAddress",
	"",
	"The host:port that the server is bound to.",
)

var databaseConnectionString = flag.String(
	"databaseConnectionString",
	"",
	"SQL database connection string",
)

var maxDatabaseConnections = flag.Int(
	"maxDatabaseConnections",
	200,
	"Max numbers of SQL database connections",
)

var databaseDriver = flag.String(
	"databaseDriver",
	"postgres",
	"SQL database driver name",
)

var sqlCACertFile = flag.String(
	"sqlCACertFile",
	"",
	"SQL database client cert, if supplied, require TLS to SQL",
)

func main() {
	flag.Parse()

	logger, _ := cflager.New("vps")
	logger.Info("starting")

	clock := clock.NewClock()

	serviceClient := vps.NewServiceClient("cuixuex@cn.ibm.com","7eab8fbfcdda3249e780dce0b10c7e4794e5ccd0fc9af7221b9fa9b40924ba8a")

	var activeDB db.DB
	var sqlDB *sqldb.SQLDB
	var sqlConn *sql.DB

	if *databaseDriver != "" && *databaseConnectionString != "" {
		var err error
		connectionString := appendSSLConnectionStringParam(logger, *databaseDriver, *databaseConnectionString, *sqlCACertFile)

		sqlConn, err = sql.Open(*databaseDriver, connectionString)
		if err != nil {
			logger.Fatal("failed-to-open-sql", err)
		}
		defer sqlConn.Close()
		sqlConn.SetMaxOpenConns(*maxDatabaseConnections)
		sqlConn.SetMaxIdleConns(*maxDatabaseConnections)

		err = sqlConn.Ping()
		if err != nil {
			logger.Fatal("sql-failed-to-connect", err)
		}

		sqlDB = sqldb.NewSQLDB(sqlConn, clock, *databaseDriver)
		err = sqlDB.CreateConfigurationsTable(logger)
		if err != nil {
			logger.Fatal("sql-failed-create-configurations-table", err)
		}
		activeDB = sqlDB
	}

	if activeDB == nil {
		logger.Fatal("no-database-configured", errors.New("no database configured"))
	}

	migrationsDone := make(chan struct{})

	migrationManager := migration.NewManager(
		logger,
		sqlDB,
		sqlConn,
		migrationsDone,
		clock,
		*databaseDriver,
	)

	exitChan := make(chan struct{})

	var accessLogger lager.Logger
	if *accessLogPath != "" {
		accessLogger = lager.NewLogger("vps-access")
		file, err := os.OpenFile(*accessLogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			logger.Error("invalid-access-log-path", err, lager.Data{"access-log-path": *accessLogPath})
			os.Exit(1)
		}
		accessLogger.RegisterSink(lager.NewWriterSink(file, lager.INFO))
	}

	handler := handlers.New(
		logger,
		accessLogger,
		activeDB,
		serviceClient,
		migrationsDone,
		exitChan,
	)

	var server ifrit.Runner

	server = http_server.New(*listenAddress, handler)

	members := grouper.Members{
		{"server", server},
		{"migration-manager", migrationManager},
	}

	group := grouper.NewOrdered(os.Interrupt, members)

	monitor := ifrit.Invoke(sigmon.New(group))
	go func() {
		// If a handler writes to this channel, we've hit an unrecoverable error
		// and should shut down (cleanly)
		<-exitChan
		monitor.Signal(os.Interrupt)
	}()

	logger.Info("started")

	err := <-monitor.Wait()
	if sqlConn != nil {
		sqlConn.Close()
	}
	if err != nil {
		logger.Error("exited-with-failure", err)
		os.Exit(1)
	}

	logger.Info("exited")
}

func appendSSLConnectionStringParam(logger lager.Logger, driverName, databaseConnectionString, sqlCACertFile string) string {
	switch driverName {
	case "mysql":
		if sqlCACertFile != "" {
			certBytes, err := ioutil.ReadFile(sqlCACertFile)
			if err != nil {
				logger.Fatal("failed-to-read-sql-ca-file", err)
			}

			caCertPool := x509.NewCertPool()
			if ok := caCertPool.AppendCertsFromPEM(certBytes); !ok {
				logger.Fatal("failed-to-parse-sql-ca", err)
			}

			tlsConfig := &tls.Config{
				InsecureSkipVerify: false,
				RootCAs:            caCertPool,
			}

			mysql.RegisterTLSConfig("vps-tls", tlsConfig)
			databaseConnectionString = fmt.Sprintf("%s?tls=bbs-tls", databaseConnectionString)
		}
	case "postgres":
		if sqlCACertFile == "" {
			databaseConnectionString = fmt.Sprintf("%s?sslmode=disable", databaseConnectionString)
		} else {
			databaseConnectionString = fmt.Sprintf("%s?sslmode=verify-ca&sslrootcert=%s", databaseConnectionString, sqlCACertFile)
		}
	}

	return databaseConnectionString
}
