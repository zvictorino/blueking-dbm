// TencentBlueKing is pleased to support the open source community by making 蓝鲸智云-DB管理系统(BlueKing-BK-DBM) available.
// Copyright (C) 2017-2023 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at https://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package config

// LogicalBackup the config of logical backup
// data or schema is controlled by Public.DataSchemaGrant
type LogicalBackup struct {
	// ChunkFilesize split tables into chunks of this output file size. This value is in MB
	ChunkFilesize uint64 `ini:"ChunkFilesize"`
	Regex         string `ini:"Regex"`
	Threads       int    `ini:"Threads"`
	// DisableCompress disable zstd compress. compress is enable by default
	DisableCompress bool   `ini:"DisableCompress"`
	FlushRetryCount int    `ini:"FlushRetryCount"`
	DefaultsFile    string `ini:"DefaultsFile"`
	ExtraOpt        string `ini:"ExtraOpt"` // other mydumper options string to be appended
}

// LogicalLoad the config of logical loading
type LogicalLoad struct {
	MysqlHost    string `ini:"MysqlHost"`
	MysqlPort    int    `ini:"MysqlPort"`
	MysqlUser    string `ini:"MysqlUser"`
	MysqlPasswd  string `ini:"MysqlPasswd"`
	MysqlCharset string `ini:"MysqlCharset"`
	MysqlLoadDir string `ini:"MysqlLoadDir"`
	Threads      int    `ini:"Threads"`
	Regex        string `ini:"Regex"`
	EnableBinlog bool   `ini:"EnableBinlog"`
	// SchemaOnly import schema,trigger,func,proc (--no-data)
	//  if you want only table schema, use ExtraOpt = -skip-triggers --skip-post
	//  mydumper doest not support data only currently, you should backup only data for your purpose
	SchemaOnly    bool   `ini:"SchemaOnly"`
	IndexFilePath string `ini:"IndexFilePath" validate:"required"`
	ExtraOpt      string `ini:"ExtraOpt"` // other myloader options string to be appended
	// DBListDropIfExists will run drop database if exists db_xxx before load data. comma separated
	DBListDropIfExists string `ini:"DBListDropIfExists"`
	// CreateTableIfNotExists true will add --append-if-not-exist for myloader
	CreateTableIfNotExists bool `ini:"CreateTableIfNotExists"`
}

// LogicalBackupMysqldump the config of logical backup with mysqldump
type LogicalBackupMysqldump struct {
	BinPath  string `ini:"BinPath"`  // the binary path of mysqldump
	ExtraOpt string `ini:"ExtraOpt"` // other mysqldump options string to be appended
}

// LogicalLoadMysqldump the config of logical loading with mysql
type LogicalLoadMysqldump struct {
	MysqlHost         string `ini:"MysqlHost"`
	MysqlPort         int    `ini:"MysqlPort"`
	MysqlUser         string `ini:"MysqlUser"`
	MysqlPasswd       string `ini:"MysqlPasswd"`
	MysqlCharset      string `ini:"MysqlCharset"`
	MysqlLoadFilePath string `ini:"MysqlLoadFilePath"`
	IndexFilePath     string `ini:"IndexFilePath" validate:"required"`
	BinPath           string `ini:"BinPath"`  // the binary path of mysql
	ExtraOpt          string `ini:"ExtraOpt"` // other mysql options string to be appended(we use mysql to load backup)
}
