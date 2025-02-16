# Quick Polling App

## Overview
A simple polling app where a user can:
- Create a poll with a question and multiple options.
- Vote on a poll.
- View poll results in real-time (auto-refresh every 5 seconds).

## API Endpoints
- **POST /api/polls**  
  Create a new poll.  
  *Request Body:*  
  ```json
  {
    "question": "Your poll question",
    "options": ["Option 1", "Option 2", "Option 3"]
  }

- **GET /api/polls/{id}**
    Retrieve a poll with its options and vote counts.
-     POST /api/polls/{id}/vote
    Cast a vote for a given option.
      *Request Body:*  
  ```json
    {
      "option_id": 2
    } 

Database Schema
Table: polls
id (serial primary key)
question (text not null)
created_at (timestamp default now())
Table: poll_options
id (serial primary key)
poll_id (integer, references polls(id) on delete cascade)
option_text (text not null)
vote_count (integer 