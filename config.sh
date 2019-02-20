name="bench_test"
pass="12345678"
address=`echo $pass | qoscli keys add $name | awk '{if(NR==2){print $3};}'`

# add monikor, pass, address, write them into config.json
echo "{
\"name\": \"$name\",
\"address\": \"$address\",
\"password\": \"$pass\"
}" > ./config.json

# assign asset 
`qosd add-genesis-accounts $address,1000000000qos`

