package config

var defaultConf = []byte(`

env = "dev"

log_config = """{
  level = "debug"

  formatter.name = "json"
  formatter.options {
    force-colors      = true
    disable-colors    = false
    disable-timestamp = false
    full-timestamp    = true
    timestamp-format  = "2018-01-01 23:59:59"
    disable-sorting   = false
  }
}"""
`)
