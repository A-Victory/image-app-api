# Image Social App API

A partial social media app created for users to be able to save memories as well as share these memories with other users.

#### MOTIVATION

- Devloping and improving in using Go for backend.
- To improve with working with databases using Go's standard library.
- Working and applying user based authentication with Cookies.
- Familiarity with working with cloud based storages.

#### FEATURES
   <ul>
        <li>Sign-Up</li>
        <li>Login</li>
        <li>Uploading Images</li>
        <li>Updating Infomations</li>
        <li>Cookies</li>
   </ul>
    
## CONNECTIONS

- ### Database
     - First create the necessary database tables in Database of choice. Here I used MySQL
     - In package DB connect to the database by inputing in the right parameters *(db, err := sql.Open("mysql", "user:password@/dbname"))*

- ### Cloudinary
     - Create a Cloudinary account
     - Get Cloud information from dashboard
     - Configure environmental variables for couldinary in *.env* file **(CLOUDINARY_URL=cloudinary://key:secret@cloud)**


_For the project a basic knowledge of operating and working with databases is important_

