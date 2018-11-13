from subprocess import call
import os
def get_immediate_subdirectories(a_dir):
    return [name for name in os.listdir(a_dir)
            if os.path.isdir(os.path.join(a_dir, name))]
def restoreInfluxDBDatabases(dir_names, hostname):
    for dirname in dir_names:
        newdbname_k8s = dirname+"_k8s"
        newdbname_TestK6 = dirname+"_TestK6"
        print ("dirName: "+ dirname)
        call(["influxd", "restore", "-portable", "-host", hostname, "-db", "k8s", "-newdb", newdbname_k8s, dirname])
        call(["influxd", "restore", "-portable", "-host", hostname, "-db", "TestK6", "-newdb", newdbname_TestK6, dirname])

dir_names = get_immediate_subdirectories("./")
restoreInfluxDBDatabases(dir_names, "141.40.254.24:8088")
