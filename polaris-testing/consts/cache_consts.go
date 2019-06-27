package consts

//Mail send limit keys + type + "email", eg. allstar:mail:send:limit:1473048656@qq.com
const MAIL_SEND_LIMIT_KEYS = "allstar:mail:send:limit:"
const MAIL_AUTH_COUNT_KEYS = "allstar:mail:auth:count:"
const MAIL_CODE_KEYS = "allstar:mail:code:"

//key - value : userId - token
const USER_TOKEN_KEYS = "allstar:user:token:"

//key - value : token - user
const TOKEN_USER_KEYS = "allstar:token:user:"

//Mail send limit (unit seconds), eg. 60s
const MAIL_SEND_LIMIT = 60
const MAIL_CODE_EXPIRE = 2 * 60
const MAIL_AUTO_COUNT_LIMIT = 5
const USER_TOKEN_EXPIRE = 60 * 60 * 24 * 2
