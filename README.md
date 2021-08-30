# Stocks leaderboard

main/main.go -- Contains the registration of different APIs based on different URLs and types of requests


users/register.go -- Code to register the users<br>

users/updateShares.go -- Code to update the share units of a particular user and get the user information
users/readValues.go -- Code logic to get the values from Kafka (share values of each company)
users/users.go -- Code logic to get the leaderboard of users in descending order of their total share values
