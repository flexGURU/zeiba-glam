export interface User {
  email: string;
  username: string;
}

export interface LoginResponse {
  data: {
    access_token: string;
    refresh_token: string;
  };
}

