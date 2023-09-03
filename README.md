# VoD_microservice

1.Launch service center  
                 
      docker-compose up 
      
2.Run services

      cd VoD microservice/server
      export CHASSIS_HOME=$PWD
      go run main.go
  

      cd VoD microservice1/server
      export CHASSIS_HOME=$PWD
      go run main.go
  
# API

--> Import "PostmanCollection.json" file to Postman.

1-> Create_End_User : POST
    
    url: http://localhost:3000/user/signUp

    Headers :     {
                        language:"en"
                  }

    sampleInput : {
                        "User_type":"ADMIN",
                        "Username":"Shashank",
                        "Email" :"Shashank@gmail.com",
                        "Password": "Shashank@123"
                  }

2-> Login : POST
    
    url: http://localhost:3000/user/login

    Headers :     {
                        language:"en"
                  }

    sampleInput : {
                        "Email" :"Shashank@gmail.com",
                        "Password": "Shashank@123"
                  }

3-> Get ALl USers : GET

    url: http://localhost:3000/users  

    Headers :     {
                        token:"",
                        language:"en"
                  }

4-> Upload Videos: POST

    url: http://localhost:3001/video/upload

    Headers :     {
                        token:"",
                        language:"en"
                  }
    

5-> Get All Videos : GET

    url: http://localhost:3001/videos?page=1&size=10&filter={"title":{"$regex":"Be"}}

    Headers :     {
                        token:"",
                        language:"en"
                  }

5-> Get By ID : GET

    url: http://localhost:3001/video/{id}

         (http://localhost:3001/video/64f49efe8e439cc351d59e36)

    Headers :     {
                        token:"",
                        language:"en"
                  }

  
