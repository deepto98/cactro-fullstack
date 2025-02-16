Polling App
=================

Overview
--------
I coded the web service in Go and used HTML templates for the UI. Its deployed at https://cactro-fullstack.onrender.com/. These are the primary components of the app:

  - main.go: Contains the main entry point for the server.
  - internal/api: Contains API handlers for poll operations.
  - internal/db: Handles database initialization and connections.
  - internal/middleware: Contains middleware (e.g., logging).
  - internal/models: Contains data models for polls and options.
  - templates: Contains Go HTML templates for the UI.
  - static: Contains static assets like CSS.

API Endpoints
-------------
1. Create Poll
   - Endpoint: POST /api/polls
   - Description: Create a new poll with a question and multiple options.
   - Request Body:
     ```
      {  
       "question": "What is your favorite programming language?",
       "options": ["Go", "Python", "JavaScript"]
      }
     ```
 
   - Response (Status: 201 Created):
     ``` 
     {
       "poll_id": 1
     }
     ```
     

2. Get Poll
   - Endpoint: GET /api/polls/{id}
   - Description: Retrieve details of a poll, including the question, options, and vote counts.
   - Response (Status: 200 OK):
     ```
     {
       "id": 1,
       "question": "What is your favorite programming language?",
       "options": [
         {
           "id": 1,
           "poll_id": 1,
           "option_text": "Go",
           "vote_count": 10
         },
         {
           "id": 2,
           "poll_id": 1,
           "option_text": "Python",
           "vote_count": 15
         },
         {
           "id": 3,
           "poll_id": 1,
           "option_text": "JavaScript",
           "vote_count": 5
         }
       ]
     }
     ``` 

3. Vote on Poll
   - Endpoint: POST /api/polls/{id}/vote
   - Description: Cast a vote for a specific option in a poll.
   - Request Body:
     ```
     {
       "option_id": 1
     }
     ```

   - Response (Status: 200 OK):

     ```
     {
       "message": "Vote registered"
     }
     ```

4. List All Polls
   - Endpoint: GET /api/polls
   - Description: Retrieve a list of all polls with basic details.
   - Response (Status: 200 OK):
     ```
     [
       {
         "id": 1,
         "question": "What is your favorite programming language?"
       },
       {
         "id": 2,
         "question": "Do you prefer coffee or tea?"
       }
     ]
     ```

Database Schema
---------------
Database:
  - PostgreSQL is used for storing poll data.
  - The schema consists of two main tables: polls and poll_options.
  - Database initialization and connection handling are managed in the internal/db package. 

 Table: polls
  - id: SERIAL PRIMARY KEY
  - question: TEXT NOT NULL
  - created_at: TIMESTAMP DEFAULT CURRENT_TIMESTAMP

Table: poll_options
  - id: SERIAL PRIMARY KEY
  - poll_id: INTEGER REFERENCES polls(id) ON DELETE CASCADE
  - option_text: TEXT NOT NULL
  - vote_count: INTEGER NOT NULL DEFAULT 0

Architecture Overview
---------------------
Backend:
  - Written in Golang following RESTful API design principles.
  - API endpoints for poll creation, retrieval, voting, and listing are implemented in the internal/api package. 

Frontend:
  - The UI is built using Go templates located in the templates folder.
  - Static assets (including CSS) are served from the static directory.

Middleware:
  - Custom middleware, such as logging, is implemented in the internal/middleware package.

Routing:
  - The server uses Gorilla Mux for routing, clearly separating API endpoints from UI routes.

Setup & Deployment
--------------------
Environment Variables:
  - DATABASE_URL: PostgreSQL connection string (e.g., postgres://username:password@localhost:5432/polling_app?sslmode=disable)
  - URL - https://cactro-fullstack.onrender.com/
 
 
