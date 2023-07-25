// api.js
export async function fetchWithToken(url, options) {
    const token = document.cookie
      .split("; ")
      .find((row) => row.startsWith("sessionId="))
      ?.split("=")[1];
  
    if (!token) {
      // Handle the case when the token is not found (user not authenticated)
      throw new Error('User not authenticated');
    }
  console.log(token)
    const headers = {
      "Content-Type": "application/json",
      Authorization: `${token}`,
    };
  
    const requestOptions = {
      ...options,
      headers,
    };
  
    return fetch(url, requestOptions);
  }
  