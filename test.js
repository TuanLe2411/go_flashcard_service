// Example: Get flashcard by ID
fetch('http://localhost:8080', {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      'user_id': 'your-user-id'  // Required by your backend
    },
  })
    .then(data => console.log(data));