# vaultConfig


## usage

./vaultConfig

optional argument : setting file path (default : setting.json)


## setting.json
```json
{
  "vaultApi": "http://127.0.0.1:8200",
  "secrets": {
    "dexSettings": {
      "taskIdSeed": "DATA",
      "signServerAddr": "DATA",
      "signServerAppName": "DATA",
      "signServerAppKey": "DATA"
    },
    "web-jwt": {
      "secret": "DATA",
      "expiration": 1234
    },
    "redis": {
      "host": "DATA",
      "port": 1234
    },
    "mariadb": {
      "username": "DATA",
      "password": "DATA",
      "addr": "DATA",
      "port": 1234,
      "schema": "DATA"
    },
    "toast-cloud": {
      "appKey": "DATA"
    },
    "java-mail": {
      "user": "DATA",
      "password": "DATA"
    },
    "coinmarketcap": {
      "apiKey": "DATA"
    },
    "mongodb": {
      "addr": "DATA",
      "port": 1234,
      "username": "DATA",
      "password": "DATA",
      "authMechanism": "DATA",
      "authSource": "DATA"
    },
    "push-aos": {
      "key": "DATA"
    },
    "push-ios": {
      "key": "DATA"
    },
    "wdxPartnerKey" : {
      "accessKey": "DATA"
    },
    "kakao_oauth" : {
      "client_id": "DATA",
      "client_secret": "DATA"
    },
    "telegram_oauth" : {
      "client_id": "DATA",
      "client_secret": "DATA"
    }
  },
  "roles": {
    "tkn-web": {
      "secret": [ "redis", "wdxPartnerKey", "kakao_oauth", "telegram_oauth" ],
      "password": "1234"
    },
    "tkn-cmu": {
      "secret": [ "redis", "mariadb", "web-jwt", "toast-cloud", "java-mail" ],
      "password": "1234"
    },
    "tkn-gov": {
      "secret": [ "redis", "mariadb", "coinmarketcap" ],
      "password": "1234"
    },
    "tkn-dex": {
      "secret": [ "redis", "mariadb", "dexSettings", "web-jwt"],
      "password": "1234"
    },
    "tkn-anc": {
      "secret": [ "redis", "mariadb"],
      "password": "1234"
    },
    "tkn-wdx": {
      "secret": [ "mongodb", "push-ios", "push-aos" ],
      "password": "1234"
    }
  }
}
```

1. roles.["ROLENAME"].hostname is optional, it will use rolename if omitted.  
both generated and given hostname will be served as hashed form
2. roles.["ROLENAME"].password is optional, it will generate random password if omitted.  
generated password will be 64byte long, and will generate new password everytime it runs  
SO, if fixed password required, set first generated password in setting.json 
3. secrets.["SECRETNAME"] can have null data as placeholder, actual secret entry will not created.


**IT'S VERY DANGEROUS TO KEEP setting.json IN PRODUCTION ENVIRONMENT, setting.json MUST BE REMOVED AFTER USE**

