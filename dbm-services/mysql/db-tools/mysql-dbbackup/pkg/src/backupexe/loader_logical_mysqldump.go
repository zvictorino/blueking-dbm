package backupexe

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"dbm-services/common/go-pubpkg/cmutil"
	"dbm-services/mysql/db-tools/mysql-dbbackup/pkg/config"
	"dbm-services/mysql/db-tools/mysql-dbbackup/pkg/src/dbareport"
	"dbm-services/mysql/db-tools/mysql-dbbackup/pkg/src/logger"
	"dbm-services/mysql/db-tools/mysql-dbbackup/pkg/src/mysqlconn"
)

// LogicalLoaderMysqldump this logical loader is used to load logical backup with mysql(client)
type LogicalLoaderMysqldump struct {
	cnf          *config.BackupConfig
	dbbackupHome string
	dbConn       *sql.DB
	initConnect  string
}

// initConfig initializes the configuration for the logical loader[mysql]
func (l *LogicalLoaderMysqldump) initConfig(_ *dbareport.IndexContent) error {
	if l.cnf == nil {
		return errors.New("logical loader params is nil")
	}
	if cmdPath, err := os.Executable(); err != nil {
		return err
	} else {
		l.dbbackupHome = filepath.Dir(cmdPath)
	}

	return nil
}

// preExecute preprocess before loading data
func (l *LogicalLoaderMysqldump) preExecute() error {
	// 临时清理 init_connect
	dbListDrop := l.cnf.LogicalLoad.DBListDropIfExists
	var initConnect string
	if err := l.dbConn.QueryRow("select @@init_connect").Scan(&initConnect); err != nil {
		return err
	}
	l.initConnect = initConnect
	if l.initConnect != "" && strings.TrimSpace(dbListDrop) != "" {
		logger.Log.Info("set global init_connect='' for safe")
		if _, err := l.dbConn.Exec("set global init_connect=''"); err != nil {
			return err
		}
	}

	// handle DBListDropIfExists
	// 如果有设置这个选项，会在运行前执行 drop database if exists 命令，来清理脏库
	if strings.TrimSpace(dbListDrop) != "" {
		logger.Log.Info("load logical DBListDropIfExists:", dbListDrop)
		if strings.Contains(dbListDrop, "`") {
			return errors.Errorf("DBListDropIfExists has invalid character %s", dbListDrop)
		}
		SysDBs := []string{"mysql", "sys", "information_schema", "performance_schema", "test"}
		dblist := strings.Split(dbListDrop, ",")
		dblistNew := []string{}
		for _, dbName := range dblist {
			dbName = strings.TrimSpace(dbName)
			if dbName == "" {
				continue
			} else if cmutil.StringsHas(SysDBs, dbName) {
				return errors.Errorf("DBListDropIfExists should not contain sys db: %s", dbListDrop)
			} else {
				dblistNew = append(dblistNew, dbName)
			}
		}

		for _, dbName := range dblistNew {
			dropDbSql := fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", dbName)
			logger.Log.Warn("DBListDropIfExists sql:", dropDbSql)
			if _, err := l.dbConn.Exec(dropDbSql); err != nil {
				return errors.Wrap(err, "DBListDropIfExists err")
			}
		}
		return nil
	}
	return nil
}

// Execute execute loading backup with logical load tool [mysql]
// for the mysqldump backup, we use mysql to load it
func (l *LogicalLoaderMysqldump) Execute() (err error) {
	cnfPublic := config.Public{
		MysqlHost:    l.cnf.LogicalLoad.MysqlHost,
		MysqlPort:    l.cnf.LogicalLoad.MysqlPort,
		MysqlUser:    l.cnf.LogicalLoad.MysqlUser,
		MysqlPasswd:  l.cnf.LogicalLoad.MysqlPasswd,
		MysqlCharset: l.cnf.LogicalLoad.MysqlCharset,
	}
	l.dbConn, err = mysqlconn.InitConn(&cnfPublic)
	if err != nil {
		return err
	}
	defer func() {
		_ = l.dbConn.Close()
	}()
	if err = l.preExecute(); err != nil {
		return err
	}

	defer func() {
		if l.initConnect != "" {
			logger.Log.Info("set global init_connect back:", l.initConnect)
			if _, err = l.dbConn.Exec(fmt.Sprintf(`set global init_connect="%s"`, l.initConnect)); err != nil {
				//return err
				logger.Log.Warn("fail set global init_connect back:", l.initConnect)
			}
		}
	}()

	var binPath string
	if l.cnf.LogicalLoadMysqldump.BinPath != "" {
		binPath = l.cnf.LogicalLoadMysqldump.BinPath
	} else {
		binPath = filepath.Join(l.dbbackupHome, "bin/mysql")
	}

	args := []string{
		"-h" + l.cnf.LogicalLoadMysqldump.MysqlHost,
		"-P" + strconv.Itoa(l.cnf.LogicalLoadMysqldump.MysqlPort),
		"-u" + l.cnf.LogicalLoadMysqldump.MysqlUser,
		"-p" + l.cnf.LogicalLoadMysqldump.MysqlPasswd,
		"--max_allowed_packet=1073741824 ",
		fmt.Sprintf("--default-character-set=%s", l.cnf.LogicalLoadMysqldump.MysqlCharset),
	}

	// ExtraOpt is to freely add command line arguments
	if l.cnf.LogicalLoadMysqldump.ExtraOpt != "" {
		args = append(args, []string{
			fmt.Sprintf(`%s`, l.cnf.LogicalLoadMysqldump.ExtraOpt),
		}...)
	}

	args = append(args, []string{
		"<", fmt.Sprintf(`'%s'`, l.cnf.LogicalLoadMysqldump.MysqlLoadFilePath),
	}...)

	pwd, _ := os.Getwd()
	logfile := filepath.Join(pwd, "logs", fmt.Sprintf("mysqldump_load_%d.log", int(time.Now().Weekday())))
	_ = os.MkdirAll(filepath.Dir(logfile), 0755)

	args = append(args, ">>", logfile, "2>&1")
	logger.Log.Info("load logical command:", binPath+" ", strings.Join(args, " "))

	outStr, errStr, err := cmutil.ExecCommand(true, "", binPath, args...)
	if err != nil {
		logger.Log.Error("load backup failed: ", err, errStr)
		return errors.Wrap(err, errStr)
	}
	logger.Log.Info("load backup success: ", outStr)
	return nil
}
