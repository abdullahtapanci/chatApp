# Chat App

## This is a chat application that people can create an account, manage their account, find new friends and chat with them. Just like Whatsapp.

## Technologies that are used for this app
### Server side
- Go fiber
- MySQL
- Go fiber websocket
- Golang JWT

###Frontend side
- React
- Tailwindcss

### Installation
1. Clone the repository
2. Run `npm install` in the client directory
3. Then run `npm run start` to start client side
4. Run `go run main.go` in the server directory

### Database 
- You should create some tables
- users table

| Field            | Type         | Null | Key | Default | Extra          |
|------------------|--------------|------|-----|---------|----------------|
| id               | int          | NO   | PRI | NULL    | auto_increment |
| firstName        | varchar(50)  | NO   |     | NULL    |                |
| lastName         | varchar(50)  | NO   |     | NULL    |                |
| userName         | varchar(50)  | NO   | UNI | NULL    |                |
| email            | varchar(100) | NO   | UNI | NULL    |                |
| password         | varchar(50)  | NO   |     | NULL    |                |
| profileImageName | varchar(100) | YES  |     | NULL    |                |

- friendships table

| Field    | Type | Null | Key | Default | Extra |
|----------|------|------|-----|---------|-------|
| user1_id | int  | NO   | PRI | NULL    |       |
| user2_id | int  | NO   | PRI | NULL    |       |

- Do not for get to update db.go ( write your username, passsword and database name in db.go )

