# vaultConfig


## usage

./vaultConfig

optional argument : setting file path (default : setting.json)


## setting.json
<pre><code>
{
  "vaultApi": "http://127.0.0.1:8200",
  "secrets": {
    "dexKey": {
      "seed": "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    },
    "web-jwt": {
      "secret": "SECRETSECRETSECRETSECRET",
      "expiration": 604800000
    },
    "redis": {
      "host": "localhost",
      "port": 6379
    },
    "mariadb": {
      "username": "root",
      "password": "PASSWORD",
      "addr": "localhost",
      "port": 3306,
      "schema": "talken_db"
    },
    "toast-cloud": {
      "appKey": "TOASTCLOUDAPPKEY"
    },
    "coinmarketcap": {
      "apiKey": "COINMARKETCAPAPPKEY"
    },
    "mongodb": {
      "addr": "localhost",
      "port": 2701,
      "username": "USERNAME",
      "password": "PASSWORD",
      "authMechanism": "MECHANISM",
      "authSource": "SOURCE"
    },
    "push-aos": {
      "key": "value"
    },
    "push-ios": {
      "key": "value"
    }
  },
  "roles": {
    "tkn-web": {
      "secret": [
        "redis"
      ],
      "hostname": "tkn-web",
      "password": "PASSWORD"
    },
    "tkn-cmu": {
      "secret": [
        "redis",
        "mariadb",
        "web-jwt",
        "toast-cloud"
      ]
    },
    "tkn-gov": {
      "secret": [
        "redis",
        "mariadb",
        "coinmarketcap"
      ],
      "password": "PASSWORD"
    },
    "tkn-dex": {
      "secret": [
        "redis",
        "mariadb",
        "dexKey",
        "web-jwt"
      ],
      "password": "PASSWORD"
    },
    "tkn-anc": {
      "secret": [
        "redis",
        "mariadb"
      ],
      "password": "PASSWORD"
    },
    "tkn-wdx": {
      "secret": [
        "mongodb",
        "push-ios",
        "push-aos"
      ],
      "password": "PASSWORD"
    }
  }
}
</code></pre>

1. roles.["ROLENAME"].hostname is optional, it will use rolename if omitted
2. roles.["ROLENAME"].password is optional, it will generate random password if omitted


