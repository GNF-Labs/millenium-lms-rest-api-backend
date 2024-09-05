Here’s the continuation of your documentation with details for each endpoint, including request payloads, responses, and usage:

---

# Route Lists

## 1. Auth Routes

### `/login`

- **Method**: `POST`
- **Description**: Authenticates a user and returns a JWT token.
- **Request Payload**:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "message": "Login successful",
      "token": "string"
    }
    ```
  - **Error (400 Bad Request / 401 Unauthorized)**:
    ```json
    {
      "error": "string"
    }
    ```

### `/register`

- **Method**: `POST`
- **Description**: Registers a new user.
- **Request Payload**:
  ```json
  {
    "username": "string",
    "password": "string",
    "email": "string"
  }
  ```
- **Response**:
  - **Success (201 Created)**:
    ```json
    {
      "message": "User registered successfully"
    }
    ```
  - **Error (400 Bad Request / 409 Conflict)**:
    ```json
    {
      "error": "string"
    }
    ```

## 2. Profile Routes

### `/profile/:username`

- **Method**: `GET`
- **Description**: Retrieves the profile information for a specified username.
- **Request Headers**:
  - `Authorization`: `Bearer <token>`
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "message": "Profile data retrieved",
      "data": {
        "username": "string",
        "email": "string",
        "full_name": "string",
        "about": "string"
      }
    }
    ```
  - **Error (401 Unauthorized / 404 Not Found)**:
    ```json
    {
      "error": "string"
    }
    ```

### `/profile/:username`

- **Method**: `PUT`
- **Description**: Updates the profile information for a specified username.
- **Request Headers**:
  - `Authorization`: `Bearer <token>`
- **Request Payload**:
  ```json
  {
    "email": "string",
    "full_name": "string",
    "about": "string"
  }
  ```
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "message": "Profile updated successfully",
      "data": {
        "username": "string",
        "email": "string",
        "full_name": "string",
        "about": "string"
      }
    }
    ```
  - **Error (400 Bad Request / 401 Unauthorized / 404 Not Found)**:
    ```json
    {
      "error": "string"
    }
    ```

## 3. Token Verification Routes

### `/verify-token/:username`

- **Method**: `GET`
- **Description**: Verifies the validity of the JWT token.
- **Request Headers**:
  - `Authorization`: `Bearer <token>`
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "message": "Token is valid",
      "data": {
        "username": "string",
        "token": "string"
      }
    }
    ```
  - **Error (401 Unauthorized)**:
    ```json
    {
      "error": "string"
    }
    ```

## 4. User Interaction Routes

### `/interact`

- **Method**: `PUT`
- **Description**: Updates user interactions with courses.
- **Request Headers**:
  - `Authorization`: `Bearer <token>`
- **Request Payload**:
  ```json
  {
    "course_id": 1,
    "action": "string"  // e.g., "viewed", "registered", "completed"
  }
  ```
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "message": "Interaction recorded"
    }
    ```
  - **Error (400 Bad Request / 401 Unauthorized)**:
    ```json
    {
      "error": "string"
    }
    ```

### `/interact/:username`

- **Method**: `GET`
- **Description**: Retrieves user interactions with courses.
- **Request Headers**:
  - `Authorization`: `Bearer <token>`
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "message": "Interactions retrieved",
      "data": [
        {
          "course_id": 1,
          "action": "string"  // e.g., "viewed", "registered", "completed"
        }
      ]
    }
    ```
  - **Error (401 Unauthorized / 404 Not Found)**:
    ```json
    {
      "error": "string"
    }
    ```

## 5. Course Routes

### `/courses`

- **Method**: `GET`
- **Description**: Retrieves a list of courses with optional pagination and search.
- **Query Parameters**:
  - `page` (optional): Page number for pagination (default: 1).
  - `q` (optional): Search query to filter courses by name or other attributes.
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "message": "Courses retrieved",
      "data": [
        {
          "id": 1,
          "name": "string"
        }
      ]
    }
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "error": "string"
    }
    ```

### `/courses/:id`

- **Method**: `GET`
- **Description**: Retrieves detailed information for a specific course.
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "message": "Course data retrieved",
      "data": {
        "id": 1,
        "name": "string",
        "description": "string",
        "chapters": [
          {
            "id": 1,
            "name": "string"
          }
        ]
      }
    }
    ```
  - **Error (404 Not Found)**:
    ```json
    {
      "error": "string"
    }
    ```

### `/dashboard/:username`

- **Method**: `GET`
- **Description**: Retrieves the dashboard information for a specified user.
- **Request Headers**:
  - `Authorization`: `Bearer <token>`
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "message": "Dashboard data retrieved",
      "data": {
        "username": "string",
        "courses": [
          {
            "id": 1,
            "name": "string"
          }
        ]
      }
    }
    ```
  - **Error (401 Unauthorized / 404 Not Found)**:
    ```json
    {
      "error": "string"
    }
    ```

### `/courses/:id/:chapter_id`

- **Method**: `GET`
- **Description**: Retrieves detailed information for a specific chapter in a course.
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "message": "Chapter data retrieved",
      "data": {
        "id": 1,
        "name": "string",
        "subchapters": [
          {
            "id": 1,
            "name": "string"
          }
        ]
      }
    }
    ```
  - **Error (404 Not Found)**:
    ```json
    {
      "error": "string"
    }
    ```

### `/courses/:id/:chapter_id/:subchapter_id`

- **Method**: `GET`
- **Description**: Retrieves detailed information for a specific subchapter in a chapter.
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "message": "Subchapter data retrieved",
      "data": {
        "id": 1,
        "name": "string",
        "content": "string"
      }
    }
    ```
  - **Error (404 Not Found)**:
    ```json
    {
      "error": "string"
    }
    ```
