#!/bin/sh

base_url=$1

countries=( "tr" "br" "de" "us" "jp" "pl" "ug" "fr" "sp" )

insert_new_data() {
    payload=$( printf '{"display_name": "yagu_%s", "country": "%s"}' $1 ${countries[$RANDOM%9]} )
    response=$( curl -s -X POST --url $base_url/user/create --data "$payload" --header 'Content-Type: application/json' )
    user_id=$( echo $response | jq .data.user_id )
    payload=$( printf '{"score_worth": %d, "user_id": %s}' $(( $RANDOM % 100 )) $user_id )
    curl -s -X POST --url $base_url/score/submit --data "$payload" --header 'Content-Type: application/json' -o /dev/null 
}

for idx in {1..1000}
do
    insert_new_data $idx &
done