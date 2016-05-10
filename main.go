package main

import (
  "database/sql"
  "gopkg.in/gorp.v1"
  "log"
  "strconv"
  "ginGorp/model"
  "github.com/gin-gonic/gin"
  _ "github.com/go-sql-driver/mysql"
)



var dbmap = initDb()

func initDb() *gorp.DbMap {
  connectString:="root:root@/test"
  db, err := sql.Open("mysql", connectString)
  checkErr(err, "sql.Open failed")
  dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
  dbmap.AddTableWithName(model.User{}, "User").SetKeys(true, "Id")
  err = dbmap.CreateTablesIfNotExists()
  checkErr(err, "Create tables failed")
  return dbmap
}

func checkErr(err error, msg string) {
  if err != nil {
    log.Fatalln(msg, err)
  }
}

func Cors() gin.HandlerFunc {
  return func(c *gin.Context) {
    c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
    c.Next()
  }
}

func main() {
  r := gin.Default()

  r.Use(Cors())

  v1 := r.Group("api/v1")
  {
  v1.GET("/users", GetUsers)
  v1.GET("/users/:id", GetUser)
  v1.POST("/users", PostUser)
  v1.PUT("/users/:id", UpdateUser)
  v1.DELETE("/users/:id", DeleteUser)
  v1.OPTIONS("/users", OptionsUser)     // POST
  v1.OPTIONS("/users/:id", OptionsUser) // PUT, DELETE
}

r.Run(":8080")
}

func GetUsers(c *gin.Context) {
  var users []model.User
  _, err := dbmap.Select(&users, "SELECT * FROM user")

  if err == nil {
    c.JSON(200, users)
    } else {
      c.JSON(200, gin.H{"error": "no user(s) into the table"})
    }

    // curl -i http://localhost:8080/api/v1/users
  }

  func GetUser(c *gin.Context) {
    id := c.Params.ByName("id")
    var user model.User
    err := dbmap.SelectOne(&user, "SELECT * FROM user WHERE id=? LIMIT 1", id)

    if err == nil {
      user_id, _ := strconv.ParseInt(id, 0, 64)

      content := &model.User{
        Id:        user_id,
        Firstname: user.Firstname,
        Lastname:  user.Lastname,
      }
      c.JSON(200, content)
      } else {
        c.JSON(404, gin.H{"error": "user not found"})
      }

      // curl -i http://localhost:8080/api/v1/users/1
    }

    func PostUser(c *gin.Context) {
      var user model.User
      c.Bind(&user)

      log.Println(user)

      if user.Firstname != "" && user.Lastname != "" {

        if insert, _ := dbmap.Exec(`INSERT INTO user (firstname, lastname) VALUES (?, ?)`, user.Firstname, user.Lastname); insert != nil {
          user_id, err := insert.LastInsertId()
          if err == nil {
            content := &model.User{
              Id:        user_id,
              Firstname: user.Firstname,
              Lastname:  user.Lastname,
            }
            c.JSON(201, content)
            } else {
              checkErr(err, "Insert failed")
            }
          }

          } else {
            c.JSON(400, gin.H{"error": "Fields are empty"})
          }

          // curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/users
        }

        func UpdateUser(c *gin.Context) {
          id := c.Params.ByName("id")
          var user model.User
          err := dbmap.SelectOne(&user, "SELECT * FROM user WHERE id=?", id)

          if err == nil {
            var json model.User
            c.Bind(&json)

            user_id, _ := strconv.ParseInt(id, 0, 64)

            user := model.User{
              Id:        user_id,
              Firstname: json.Firstname,
              Lastname:  json.Lastname,
            }

            if user.Firstname != "" && user.Lastname != "" {
              _, err = dbmap.Update(&user)

              if err == nil {
                c.JSON(200, user)
                } else {
                  checkErr(err, "Updated failed")
                }

                } else {
                  c.JSON(400, gin.H{"error": "fields are empty"})
                }

                } else {
                  c.JSON(404, gin.H{"error": "user not found"})
                }

                // curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/users/1
              }

              func DeleteUser(c *gin.Context) {
                id := c.Params.ByName("id")

                var user model.User
                err := dbmap.SelectOne(&user, "SELECT * FROM user WHERE id=?", id)

                if err == nil {
                  _, err = dbmap.Delete(&user)

                  if err == nil {
                    c.JSON(200, gin.H{"id #" + id: "deleted"})
                    } else {
                      checkErr(err, "Delete failed")
                    }

                    } else {
                      c.JSON(404, gin.H{"error": "user not found"})
                    }

                    // curl -i -X DELETE http://localhost:8080/api/v1/users/1
                  }

                  func OptionsUser(c *gin.Context) {
                    c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE,POST, PUT")
                    c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
                    c.Next()
                  }
