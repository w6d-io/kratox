#!/bin/bash

############################################################
# Help                                                     #
############################################################
Help()
{
   # Display Help
   echo "Description of the script functions here"
   echo
   echo "This script create configuration file for kratos server and meke call to the selfservice of Kratos"
   echo
   echo "options:"
   echo " - first   parameter \${1}  is the target action."
   echo "                - The basic action is <<config>> : use to create all the config files necessary for run a kratos server"
   echo "                - The advanced action gives the feature flow you want to apply on the kratos server."
   echo "                      The flow names you can call in the first parameter for an advanced action are : "
   echo "                      - <<register>>"
   echo "                      - <<register-login>>"
   echo "                      - <<register-login-logout>>"
   echo " - seconde parameter \${2}   is the target name of the user you want apply on the action. Not necessary for <<config>> action"
   echo " - third   parameter \${3}   is the target password of the user you want apply on the action. Not necessary for <<config>> action"
   echo
   echo "Debug is present in this script if you want see all detail of the workflow"
   echo
}

############################################################
############################################################
# Main program                                             #
############################################################
############################################################

# Display the description of the script functions here.
if [ $# -eq 0 ]
then
  Help
fi

# Create kratos config files if $1 not passed to this script
if [[ "${1}" = "config" ]]; then
cat << EOF1 >  bin/identity.example.schema.json
	{
	"\$id": "https://schemas.ory.sh/presets/kratos/quickstart/email-password/identity.schema.json",
	"\$schema": "http://json-schema.org/draft-07/schema#",
	"title": "Person",
	"type": "object",
	"properties": {
	  "traits": {
	    "type": "object",
	    "properties": {
	      "email": {
	        "type": "string",
	        "format": "email",
	        "title": "E-Mail",
	        "minLength": 3,
	        "ory.sh/kratos": {
	          "credentials": {
	            "password": {
	              "identifier": true
	            }
	          },
	           "verification": {
	            "via": "email"
	          },
	          "recovery": {
	            "via": "email"
	          }
	        }
	      },
	      "name": {
	        "type": "object",
	        "properties": {
	          "first": {
	            "title": "First Name",
	            "type": "string"
	          },
	          "last": {
	            "title": "Last Name",
	            "type": "string"
	          }
	        }
	      }
	    },
	    "required": [
	      "email"
	     ]
	  }
	}
	}
EOF1

cat << EOF2 > bin/kratos.yml
version: v0.9.0-alpha.3
dsn: memory
serve:
  public:
    base_url: http://127.0.0.1:4433/
    cors:
      enabled: true
  admin:
    base_url: http://kratos:4434/
selfservice:
  default_browser_return_url: http://127.0.0.1:4455/
  allowed_return_urls:
    - http://127.0.0.1:4455
  methods:
    password:
      enabled: true
  flows:
    error:
      ui_url: http://127.0.0.1:4455/error
    settings:
      ui_url: http://127.0.0.1:4455/settings
      privileged_session_max_age: 15m
    recovery:
      enabled: true
      ui_url: http://127.0.0.1:4455/recovery
    verification:
      enabled: true
      ui_url: http://127.0.0.1:4455/verification
      after:
        default_browser_return_url: http://127.0.0.1:4455/
    logout:
      after:
        default_browser_return_url: http://127.0.0.1:4455/login
    login:
      ui_url: http://127.0.0.1:4455/login
      lifespan: 10m
    registration:
      lifespan: 10m
      ui_url: http://127.0.0.1:4455/registration
      after:
        password:
          hooks:
            -
              hook: session
log:
  level: debug
  format: text
  leak_sensitive_values: true
secrets:
  cookie:
    - PLEASE-CHANGE-ME-I-AM-VERY-INSECURE
  cipher:
    - 32-LONG-SECRET-NOT-SECURE-AT-ALL
ciphers:
  algorithm: xchacha20-poly1305
hashers:
  algorithm: bcrypt
  bcrypt:
    cost: 8
identity:
  default_schema_id: default
  schemas:
    - id: default
      url: file://$PWD/bin/identity.example.schema.json
courier:
  smtp:
    connection_uri: smtps://test:test@mailslurper:1025/?skip_ssl_verify=true
EOF2

else






############################################################
############################################################
######
###### FUNCTIONAL FLOW
###### REGISTER
######
############################################################
############################################################


if [[ "${1}" = "register" ]]; then

cookieJar=$(mktemp)

NAME=$2
PASSWORD=$3

flow=$(curl -s  --cookie $cookieJar --cookie-jar $cookieJar -X GET \
    -H "Accept: application/json" \
    "http://127.0.0.1:4433/self-service/registration/api")


# Get the flow ID for DEBUG
#flowId=$(echo $flow | jq -r '.id')

# Get the action URL
actionUrl=$(echo $flow | jq -r '.ui.action')

# Get the CSRF Token
csrfToken=$( \
  echo $flow | \
    jq -r '.ui.nodes[] | select(.attributes.name=="csrf_token") | .attributes.value' \
)

#DEBUG
#echo $flow
#echo $flowId
#echo $actionUrl
#echo $csrfToken

# Complete the registration
session=$(curl -s --cookie $cookieJar --cookie-jar $cookieJar -X POST \
    -H "Content-Type: application/json" \
    -H  "Accept: application/json" \
    -d '{"method": "password","traits.email": "'$NAME'@ory.sh","traits.name.first":"firstname_'$NAME'","traits.name.last":"lastname_'$NAME'","password": "'$PASSWORD'", "csrf_token":""}' \
    "$actionUrl")

resp=$(echo $session  )

echo $resp
fi





############################################################
############################################################
######
###### FUNCTIONAL FLOW
###### REGISTER AND LOGIN
######
############################################################
############################################################


if [[ "${1}" = "register-login" ]]; then

# We use this cookie jar to initiate the login flow
cookieJar=$(mktemp)

# Username/email and password for an existing account
NAME=$2
username="$2@ory.sh"
password=$3

flow=$(curl -s  --cookie $cookieJar --cookie-jar $cookieJar -X GET \
    -H "Accept: application/json" \
    "http://127.0.0.1:4433/self-service/registration/api")


# Get the flow ID for DEBUG
#flowId=$(echo $flow | jq -r '.id')

# Get the action URL
actionUrl=$(echo $flow | jq -r '.ui.action')

# Get the CSRF Token
csrfToken=$( \
  echo $flow | \
    jq -r '.ui.nodes[] | select(.attributes.name=="csrf_token") | .attributes.value' \
)

#DEBUG
#echo $flow
#echo $flowId
#echo $actionUrl
#echo $csrfToken

# Complete the registration
session=$(curl -s --cookie $cookieJar --cookie-jar $cookieJar -X POST \
    -H "Content-Type: application/json" \
    -H  "Accept: application/json" \
    -d '{"method": "password","traits.email": "'$NAME'@ory.sh","traits.name.first":"firstname_'$NAME'","traits.name.last":"lastname_'$NAME'","password": "'$password'", "csrf_token":""}' \
    "$actionUrl")

resp=$(echo $session  )

#DEBUG
#echo $resp

# Initialize the flow
flow=$( \
  curl -s -H "Accept: application/json"  --cookie $cookieJar --cookie-jar $cookieJar \
    'http://127.0.0.1:4433/self-service/login/browser' \
)

# Get the action URL
actionUrl=$(echo $flow | jq -r '.ui.action')

# Get the CSRF Token
csrfToken=$( \
  echo $flow | \
    jq -r '.ui.nodes[] | select(.attributes.name=="csrf_token") | .attributes.value' \
)

#DEBUG
#cat $cookieJar

# Complete the login
session=$( \
  curl -s --cookie $cookieJar --cookie-jar $cookieJar -X POST \
    -H "Accept: application/json" -H "Content-Type: application/json" \
    --data '{ "identifier": "'$username'", "password": "'$password'", "method": "password", "csrf_token": "'$csrfToken'" }' \
    "$actionUrl" \
)

#DEBUG
#echo $session | jq

resp=$(cat $cookieJar | grep -o "ory_kratos_session.*" | awk  '{print $2}')
echo -n $resp
fi





############################################################
############################################################
######
###### FUNCTIONAL FLOW
###### REGISTER LOGIN AND LOGOUT
######
############################################################
############################################################


if [[ "${1}" = "register-login-logout" ]]; then
# We use this cookie jar to initiate the login flow
cookieJar=$(mktemp)

# Username/email and password for an existing account
NAME=$2
username="$2@ory.sh"
password=$3

flow=$(curl -s  --cookie $cookieJar --cookie-jar $cookieJar -X GET \
    -H "Accept: application/json" \
    "http://127.0.0.1:4433/self-service/registration/api")


# Get the flow ID for DEBUG
#flowId=$(echo $flow | jq -r '.id')

# Get the action URL
actionUrl=$(echo $flow | jq -r '.ui.action')

# Get the CSRF Token
csrfToken=$( \
  echo $flow | \
    jq -r '.ui.nodes[] | select(.attributes.name=="csrf_token") | .attributes.value' \
)

#DEBUG
#echo $flow
#echo $flowId
#echo $actionUrl
#echo $csrfToken

# Complete the registration
session=$(curl -s --cookie $cookieJar --cookie-jar $cookieJar -X POST \
    -H "Content-Type: application/json" \
    -H  "Accept: application/json" \
    -d '{"method": "password","traits.email": "'$NAME'@ory.sh","traits.name.first":"firstname_'$NAME'","traits.name.last":"lastname_'$NAME'","password": "'$password'", "csrf_token":""}' \
    "$actionUrl")

resp=$(echo $session  )

#DEBUG
#echo $resp

# Initialize the flow
flow=$( \
  curl -s -H "Accept: application/json"  --cookie $cookieJar --cookie-jar $cookieJar \
    'http://127.0.0.1:4433/self-service/login/browser' \
)

# Get the action URL
actionUrl=$(echo $flow | jq -r '.ui.action')

# Get the CSRF Token
csrfToken=$( \
  echo $flow | \
    jq -r '.ui.nodes[] | select(.attributes.name=="csrf_token") | .attributes.value' \
)

#DEBUG
#cat $cookieJar

# Complete the login
session=$( \
  curl -s -v --cookie $cookieJar --cookie-jar $cookieJar -X POST \
    -H "Accept: application/json" -H "Content-Type: application/json" \
    --data '{ "identifier": "'$username'", "password": "'$password'", "method": "password", "csrf_token": "'$csrfToken'" }' \
    "$actionUrl" \
)

#DEBUG
#echo $session | jq
#cat $cookieJar | grep -o "ory_kratos_session.*" | awk  '{print $2}'
#cat $cookieJar

# Check the current user id
curl -s --cookie $cookieJar --cookie-jar $cookieJar -H "Accept: application/json" \
  http://127.0.0.1:4433/sessions/whoami | \
  jq -r ".id"

  # Get the Logout URL
  logoutUrl=$( \
    curl -s --cookie $cookieJar --cookie-jar $cookieJar -H "Accept: application/json" \
      http://127.0.0.1:4433/self-service/logout/browser | \
      jq -r ".logout_url" \
  )

  # Complete the logout
  curl -s --cookie $cookieJar --cookie-jar $cookieJar "$logoutUrl"

  # Check the current user id again. It should be `null` after a successful logout
  curl -s --cookie $cookieJar --cookie-jar $cookieJar -H "Accept: application/json" \
    http://127.0.0.1:4433/sessions/whoami | \
    jq -r ".id"

#DEBUG
#    cat $cookieJar
fi

fi