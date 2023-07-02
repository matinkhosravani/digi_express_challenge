#!/bin/bash

source_code_path="$(pwd)"


# Step 1: Copy the .env.example file
echo "Step 2: Copying .env.example..."
cp ${source_code_path}/env.example ${source_code_path}/.env
cp ${source_code_path}/testing.env.example ${source_code_path}/testing.env

# Step 2: Prompt for database name and root password
echo "Step 3: Setting up the database..."
read -p "Enter database name: ex:digiexp " dbname
read -p "Enter root password ex:123456: " rootpwd

# Step 3: Update .env file in docker-compose
echo "Step 5: Updating .env file in docker-compose..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s/MYSQL_ROOT_PASSWORD=.*/MYSQL_ROOT_PASSWORD=${rootpwd}/" ./.env
else
    sed -i "s/MYSQL_ROOT_PASSWORD=.*/MYSQL_ROOT_PASSWORD=${rootpwd}/" .env
fi

# Step 4: Create database in MySQL container
docker-compose up -d mysql
echo "Waiting for mysql to become ready..."
sleep 15
docker-compose exec mysql mysql -uroot -p${rootpwd} -e "CREATE DATABASE IF NOT EXISTS ${dbname};"
docker-compose exec mysql mysql -uroot -p${rootpwd} -e "CREATE DATABASE IF NOT EXISTS ${dbname}_test;"

# Step 5: Update .env file in Go source code
echo "Step 6: Updating .env file in Go source code..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s/DB_NAME=.*/DB_NAME=${dbname}/" $source_code_path/.env 
    sed -i '' "s/DB_PASS=.*/DB_PASS=${rootpwd}/" $source_code_path/.env
    sed -i '' "s/MYSQL_ROOT_PASSWORD=.*/MYSQL_ROOT_PASSWORD=${rootpwd}/" $source_code_path/.env
    sed -i '' "s/DB_NAME=.*/DB_NAME=${dbname}_test/" $source_code_path/testing.env 
    sed -i '' "s/DB_PASS=.*/DB_PASS=${rootpwd}/" $source_code_path/testing.env
    sed -i '' "s/MYSQL_ROOT_PASSWORD=.*/MYSQL_ROOT_PASSWORD=${rootpwd}/" $source_code_path/testing.env
else
    sed -i "s/DB_PASS=.*/DB_PASS=${rootpwd}/" $source_code_path/.env
    sed -i "s/DB_NAME=.*/DB_NAME=${dbname}/" $source_code_path/.env
    sed -i "s/MYSQL_ROOT_PASSWORD=.*/MYSQL_ROOT_PASSWORD=${rootpwd}/" $source_code_path/.env
    sed -i "s/DB_PASS=.*/DB_PASS=${rootpwd}/" $source_code_path/testing.env
    sed -i "s/DB_NAME=.*/DB_NAME=${dbname}_test/" $source_code_path/testing.env
    sed -i "s/MYSQL_ROOT_PASSWORD=.*/MYSQL_ROOT_PASSWORD=${rootpwd}/" $source_code_path/testing.env
fi

# Step 6: dockeer-compose up
docker-compose up -d --build

echo "Server is running at port 8888"

