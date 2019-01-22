package config

var defaultConf = []byte(`

env = "dev"

initial_balance = 0
keep_blocks = 10000
db_cache_size = 4096

log_config = """{
  level = "info"

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
