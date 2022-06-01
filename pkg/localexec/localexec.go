package localexec

import (
	"autodg/service"
	"fmt"
	"github.com/pingcap/log"
	"go.uber.org/zap"
	"io"
	"os"
	"os/exec"
)

// RunLocalCommand 执行本地命令
func RunLocalCommand(cmd string) (string, error) {
	// todo: 考虑autodg是否支持windows平台？短期内不会支持
	// 此处是linux版本
	c := exec.Command("bash", "-c", cmd)

	// 此处是windows版本
	//c := exec.Command("cmd", "/C", cmd)

	stdoutStderr, err := c.CombinedOutput()
	//if err != nil {
	//	log.Fatal(zap.String(err))
	//}
	fmt.Println("[CMD output] : ", string(stdoutStderr))
	fmt.Printf("[CMD output] : %s\n", stdoutStderr)
	return string(stdoutStderr), err
}

// CopylocalFile 操作系统本地文件拷贝
func CopylocalFile(sourceFile, destinationFile string) (int64, error) {
	op, err1 := os.Open(sourceFile)
	of, err2 := os.Create(destinationFile)
	if err1 != nil || err2 != nil {
		fmt.Println("文件拷贝失败", err1, err2)
		log.Warn("standby db 本地文件拷贝失败",
			zap.String("待拷贝文件", sourceFile),
			zap.String("拷贝目标地址", destinationFile),
			//正常情况err1 为”nil“ 空
			//zap.String("待拷贝文件",fmt.Sprintf(err1.Error())),
			zap.String("拷贝目标地址", err2.Error()))
		os.Exit(1)
		//return err2
	}
	defer op.Close()
	defer of.Close()

	nBytes, err := io.Copy(of, op)
	if err != nil {
		fmt.Printf("The copy operation failed : %q\n", err)
		//os.Exit(1)
	} else {
		fmt.Printf("The copy operation Copied %d bytes! : \n", nBytes)
	}
	return nBytes, err
}

// Copytns 拷贝tnsnames 到备库$ORACLE_HOME/network/admin 目录
func Copytns(scfg service.SourceConfig) (int64, error) {
	var srctns = "tnsnames.ora"

	//todo: need to change destns
	//var destns = fmt.Sprintf(`%s/network/admin/tnsnames.ora/`,scfg.StandbyOracleHome)

	var destns = fmt.Sprintf(`%s/network/admin/tnsnames.ora`, scfg.StandbyOracleHome)
	//fmt.Println("[Copytns() scfg.StandbyOracleHome]:", scfg.StandbyOracleHome)
	//fmt.Println("[Copytns() destns]:", destns)
	log.Info("copy tnsnames to $ORACLE_HOME/network/admin/", zap.String("destns", destns))

	out, err := CopylocalFile(srctns, destns)
	if err != nil {
		return out, err
		os.Exit(1)
	}
	return out, err
}

// CopyOrapw 拷贝orapw$SID 到备库$ORACLE_HOME/dbs 目录
func CopyOrapw(scfg service.SourceConfig, oracle_sid string) (int64, error) {
	var src = "orapw" + oracle_sid
	var Dest string
	//todo: need to change dest
	//var dest = fmt.Sprintf(`%s/dbs/orawd%s`,scfg.StandbyOracleHome,oracle_sid)
	if scfg.IsRAC == "FALSE" {
		var Dest = fmt.Sprintf(`%s/orawd%s`, scfg.StandbyOracleHome, oracle_sid)
		//fmt.Println("[CopyOrapw() dest]:", scfg.StandbyOracleHome)
		//fmt.Println("[CopyOrapw() dest]:", Dest)
		log.Info("copy orapw to $ORACLE_HOME/dbs/", zap.String("destns", Dest))
	} else {
		var Dest = fmt.Sprintf(`%s/orawd%s1`, scfg.StandbyOracleHome, oracle_sid)
		//fmt.Println("[CopyOrapw() dest]:", scfg.StandbyOracleHome)
		//fmt.Println("[CopyOrapw() dest]:", Dest)
		log.Info("copy orapw to $ORACLE_HOME/dbs/", zap.String("destns", Dest))
	}
	out, err := CopylocalFile(src, Dest)
	if err != nil {
		return out, err
		os.Exit(1)
	}

	return out, err
}

// InitStdListener 初始化备库监听文件listener.ora
// single node : listener owner oracle
// RAC         : listener owner  grid
func InitStdListener(cfg service.SourceConfig) (string error) {
	if cfg.IsRAC == "FALSE" {
		cmds := fmt.Sprintf(`su - %s -c 'echo "
LISTENER =
(DESCRIPTION_LIST =
(DESCRIPTION =
(ADDRESS = (PROTOCOL = IPC)(KEY = ANYTHING))
(ADDRESS = (PROTOCOL = TCP)(HOST = %s)(PORT = 1521))
) )

SID_LIST_LISTENER =
(SID_LIST =
(SID_DESC =
	(GLOBAL_DBNAME = %s)
	(ORACLE_HOME = %s)
	(SID_NAME = %s)
) )"' >> %s/network/admin/listener.ora`, cfg.StandbyOracleHomeOwner, cfg.StandbyHostIps[0], cfg.OracleDBname, cfg.StandbyOracleHome, cfg.OracleSid, cfg.StandbyGridHome)
		fmt.Println("[init listener cmds]", cmds)
		RunLocalCommand(cmds)
		return string
	} else {
		cmds := fmt.Sprintf(`su - %s -c 'echo "
LISTENER =
(DESCRIPTION_LIST =
(DESCRIPTION =
(ADDRESS = (PROTOCOL = IPC)(KEY = ANYTHING))
(ADDRESS = (PROTOCOL = TCP)(HOST = %s)(PORT = 1521))
) )

SID_LIST_LISTENER =
(SID_LIST =
(SID_DESC =
	(GLOBAL_DBNAME = %s)
	(ORACLE_HOME = %s)
	(SID_NAME = %s)
) )"' >> %s/network/admin/listener.ora`, cfg.StandbyGridHomeOwner, cfg.StandbyHostIps[0], cfg.OracleDBname, cfg.StandbyGridHome, cfg.OracleSid, cfg.StandbyGridHome)
		fmt.Println("[init listener cmds]", cmds)
		RunLocalCommand(cmds)
		return string
	}
	return
}

// StartListener 启动备库本地监听
// single node : listener owner oracle
// RAC         : listener owner  grid
func StartListener(cfg service.SourceConfig) (res string, string error) {
	if cfg.IsRAC == "FALSE" {
		Cmds := fmt.Sprintf(`su - %s -c "%s/bin/lsnrctl start"`, cfg.StandbyOracleHomeOwner, cfg.StandbyOracleHome)
		fmt.Println("[Start local listener]", Cmds)

		res, err := RunLocalCommand(Cmds)
		if err != nil {
			return res, err
			log.Info("start standby oracle listener failed", zap.String("err", res))
			os.Exit(1)
		}
		return res, err
	} else {
		Cmds := fmt.Sprintf(`su - %s %s/bin/lsnrctl start`, cfg.StandbyGridHomeOwner, cfg.StandbyGridHome)
		fmt.Println("[Start local listener]", Cmds)

		res, err := RunLocalCommand(Cmds)
		if err != nil {
			return res, err
			log.Info("start standby grid listener failed", zap.String("err", res))
			os.Exit(1)
		}
		return res, err
	}
	return
}

// InitStdInstancePfile Create the  initialization file for the ###_STANDBY_DB_INSTANCE_### instance:
func InitStdInstancePfile(cfg service.SourceConfig) (string error) {
	//todo: need to check db_unique_name&cluster_database
	cmds := fmt.Sprintf(`cat <<EOF > %s/dbs/init%s.ora
db_name=%s
#db_unique_name=%suniq
#cluster_database=false
EOF `, cfg.StandbyOracleHome, cfg.OracleSid, cfg.OracleDBname, cfg.OracleSid+"uniq")

	fmt.Println("[Init standby instance pfile] ： ", cmds)

	RunLocalCommand(cmds)

	return string
}

// StartStdInstance start standby instance （nomount）初始化备库实例
func StartStdInstance(cfg service.SourceConfig) (string error) {
	cmds := fmt.Sprintf(`su - %s -c "%s/bin/sqlplus -S <<EOF2
  sys as sysdba
  startup nomount;
  alter system register;
  quit;
EOF2"`, cfg.StandbyOracleHomeOwner, cfg.StandbyOracleHome)

	fmt.Println("[start standby Aux instance ]", cmds)

	RunLocalCommand(cmds)

	return string
}

// StartRmanDuplicate start standby rman duplicate from active database
func StartRmanDuplicate(cfg service.SourceConfig) (string error) {
	parameter_value_convert := fmt.Sprintf(`'%s','%s'`, cfg.OracleSid, cfg.OracleUniqname+"uniq")
	db_unique_name := fmt.Sprintf(`'%s'`, cfg.OracleUniqname+"uniq")
	// todo: #5 not important,some parameters need to 定制?
	//control_files='%s/%s/control.ctl'
	//log_archive_max_processes='5'
	//fal_client := fmt.Sprintf(`'%s'`,"std1900")
	//fal_server := fmt.Sprintf(`'%s'`,"pri1900")
	//standby_file_management := fmt.Sprintf(`'%s'`,"AUTO")
	log_archive_config := fmt.Sprintf(`dg_config=(%s,%s)`, cfg.OracleUniqname, cfg.OracleUniqname+"uniq")
	log_archive_dest_1 := fmt.Sprintf(`LOCATION=%s VALID_FOR=(STANDBY_LOGFILE,STANDBY_ROLE) DB_UNIQUE_NAME=%s`, cfg.StandbyDataDg, cfg.OracleUniqname)
	audit_file_dest := fmt.Sprintf(`%s/%suniq/adump`, cfg.StandbyOracleBase, cfg.OracleSid)
	db_create_file_dest := fmt.Sprintf(`%s`, cfg.StandbyDataDg)
	db_recovery_file_dest := fmt.Sprintf(`%s`, cfg.StandbyDataDg)
	target := fmt.Sprintf(`%s@%s`, cfg.Password, "pri1900")
	auxiliary := fmt.Sprintf(`%s@%s`, cfg.Password, "std1900")

	// primary db  DG  config paramters
	pri_log_archive_config := log_archive_config
	pri_log_archive_dest_5 := fmt.Sprintf(`service=std1900 LGWR ASYNC valid_for=(online_logfiles,primary_role) db_unique_name=%suniq`, cfg.OracleUniqname)

	cmds := fmt.Sprintf(`rman <<EOF
connect target sys/%s;
connect auxiliary sys/%s;
run {
allocate channel prmy1 type disk;
allocate channel prmy2 type disk;
allocate channel prmy3 type disk;
allocate channel prmy4 type disk;
allocate auxiliary channel stby type disk;

duplicate target database for standby from active database
spfile
parameter_value_convert %s;
set db_unique_name= %s
#set db_file_name_convert=
#set log_file_name_convert=
#set control_files=
SET SGA_TARGET 4096M;
set log_archive_max_processes='5'
set fal_client='std1900'
set fal_server='pri1900'
set standby_file_management='AUTO'
set log_archive_config='%s'
set log_archive_dest_1='%s'
#set diagnostic_dest=
set audit_file_dest='%s'
set db_create_file_dest='%s'
set db_recovery_file_dest='%s'
nofilenamecheck
;

sql channel prmy1 "alter system set log_archive_config=''%s''";
sql channel prmy1 "alter system set log_archive_dest_5= ''%s''";
sql channel prmy1 "alter system set log_archive_dest_state_5=enable";
sql channel prmy1 "alter system set log_archive_max_processes=5";
sql channel prmy1 "alter system set fal_client=pri1900";
sql channel prmy1 "alter system set fal_server=std1900";
sql channel prmy1 "alter system set standby_file_management=auto";
sql channel prmy1 "alter system set parallel_execution_message_size=8192 scope=spfile sid=''*''";
sql channel prmy1 "alter system archive log current";

sql channel stby "alter database recover managed standby database using current logfile disconnect";
}
EOF`, target, auxiliary, parameter_value_convert, db_unique_name, log_archive_config, log_archive_dest_1,
		audit_file_dest, db_create_file_dest, db_recovery_file_dest, pri_log_archive_config, pri_log_archive_dest_5)

	fmt.Println("[start standby Aux RMAN ]", cmds)

	_, error := RunLocalCommand(cmds)

	return error
}

// MkdirStdOraAdump 创建备库实例adump 目录
func MkdirStdOraAdump(cfg service.SourceConfig) (string error) {
	cmds := fmt.Sprintf(`su - %s -c "mkdir -p %s/%s/adump"`, cfg.StandbyOracleBase, cfg.OracleUniqname+"uniq")

	fmt.Println("[Mkdir standby instance adump dir]", cmds)

	RunLocalCommand(cmds)

	return string
}
