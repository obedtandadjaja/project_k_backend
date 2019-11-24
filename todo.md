Have authentication flow be in backend:
- do middleware to pull out `current_user_id` from jwt
- when passing back response to client, append new `jwt`

- have client to do initial session token fetch
  - how to access state from helper class
  - implement rest client
